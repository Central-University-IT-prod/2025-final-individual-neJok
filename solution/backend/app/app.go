package app

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/redis/go-redis/v9"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"neJok/solution/config"
	"neJok/solution/internal/handler"
	"neJok/solution/internal/repository"
	"neJok/solution/internal/service"
	mongoUtil "neJok/solution/pkg/mongo"
	"strings"
)

func getRedis(cfg config.Config) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr: strings.Join([]string{cfg.RedisHost, cfg.RedisPort}, ":"),
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, fmt.Errorf("connect to Redis server: %w", err)
	}

	return client, nil
}

func getDatabase(cfg config.Config) (*mongo.Database, error) {
	clientOptions := options.Client().ApplyURI(cfg.BuildDsn()).SetRegistry(mongoUtil.MongoRegistry)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, fmt.Errorf("connect to MongoDB: %s", err)
	}

	err = client.Ping(context.Background(), readpref.Primary())
	if err != nil {
		return nil, fmt.Errorf("ping MongoDB: %s", err)
	}

	db := client.Database(cfg.DBName)

	return db, nil
}

func getS3(cfg config.Config) (*minio.Client, error) {
	minioClient, err := minio.New(cfg.S3Endpoint, &minio.Options{
		Creds: credentials.NewStaticV4(cfg.S3AccessKey, cfg.S3SecretKey, ""),
	})
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	err = minioClient.MakeBucket(ctx, cfg.S3Bucket, minio.MakeBucketOptions{Region: cfg.S3Region})
	if err != nil {
		exists, errBucketExists := minioClient.BucketExists(ctx, cfg.S3Bucket)
		if errBucketExists == nil && exists {
			log.Printf("We already own %s\n", cfg.S3Bucket)
		} else {
			log.Fatalln(err)
		}
	} else {
		log.Printf("Successfully created %s\n", cfg.S3Bucket)

		policy := fmt.Sprintf(`{
			"Version": "2012-10-17",
			"Statement": [
				{
					"Action": ["s3:GetObject"],
					"Effect": "Allow",
					"Principal": {"AWS": ["*"]},
					"Resource": ["arn:aws:s3:::%s/*"],
					"Sid": ""
				}
			]
		}`, cfg.S3Bucket)
		err := minioClient.SetBucketPolicy(ctx, cfg.S3Bucket, policy)
		if err != nil {
			log.Fatalf("Error setting bucket policy: %v", err)
		} else {
			log.Printf("Bucket %s set to public read-only access", cfg.S3Bucket)
		}
	}

	return minioClient, nil
}

func InitMlScores(mlScoreSvc *service.MLScoreService, actCacheSvc *service.ActCacheService) {
	mlScores, err := mlScoreSvc.GetAll()
	if err != nil {
		log.Fatalln(err)
	}

	for _, mlScore := range mlScores {
		err = actCacheSvc.SetInt(fmt.Sprintf("%s:%s", mlScore.AdvertiserID.String(), mlScore.ClientID.String()), *mlScore.Score)
		if err != nil {
			log.Fatalln(err)
		}
	}
}

func CreateRouter(cfg config.Config) *gin.Engine {
	db, err := getDatabase(cfg)
	if err != nil {
		log.Fatal(err)
	}

	redisDb, err := getRedis(cfg)
	if err != nil {
		log.Fatal(err)
	}

	s3Session, err := getS3(cfg)
	if err != nil {
		log.Fatal(err)
	}

	clientRepo := repository.NewClientRepo(db)
	advertiserRepo := repository.NewAdvertiserRepo(db)
	mlScoreRepo := repository.NewMLScoreRepo(db)
	campaignRepo := repository.NewCampaignRepo(db)
	actCacheRepo := repository.NewActCacheRepo(redisDb)
	adsHistoryRepo := repository.NewAdsHistoryRepo(db)
	s3Repo := repository.NewS3Repo(s3Session, cfg.S3Bucket, cfg.S3Public)
	gigaChatRepo := repository.NewGigaChatRepo(&cfg)
	openAIRepo := repository.NewOpenAIRepo(&cfg)

	clientSvc := service.NewClientService(clientRepo)
	advertiserSvc := service.NewAdvertiserService(advertiserRepo)
	mlScoreSvc := service.NewMLScoreService(mlScoreRepo)
	campaignSvc := service.NewCampaignService(campaignRepo)
	actCacheSvc := service.NewActCacheService(actCacheRepo)
	adsHistorySvc := service.NewAdsHistoryService(adsHistoryRepo)
	s3Svc := service.NewS3Service(s3Repo)
	gigaChatSvc := service.NewGigaChatService(gigaChatRepo, actCacheRepo)
	openAISvc := service.NewOpenAIService(openAIRepo)

	InitMlScores(mlScoreSvc, actCacheSvc)
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	_ = router.SetTrustedProxies(nil)

	pingHandler := handler.NewPingHandler()
	router.GET("/ping", pingHandler.Ping)

	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	apiClients := router.Group("/clients")
	clientHandler := handler.NewClientHandler(clientSvc)
	apiClients.POST("/bulk", clientHandler.CreateOrUpdate)
	apiClients.GET("/:clientID", clientHandler.GetByID)

	apiAdvertisers := router.Group("/advertisers")
	advertiserHandler := handler.NewAdvertiserHandler(advertiserSvc, actCacheSvc)
	apiAdvertisers.POST("/bulk", advertiserHandler.CreateOrUpdate)
	apiAdvertisers.GET("/:advertiserID", advertiserHandler.GetByID)

	apiMLScores := router.Group("/ml-scores")
	mlScoresHandler := handler.NewMLScoreHandler(mlScoreSvc, clientSvc, advertiserSvc, actCacheSvc)
	apiMLScores.POST("", mlScoresHandler.CreateOrUpdate)

	apiCampaigns := apiAdvertisers.Group("/:advertiserID/campaigns")
	campaignsHandler := handler.NewCampaignHandler(campaignSvc, advertiserSvc, actCacheSvc, s3Svc, openAISvc)
	apiCampaigns.POST("", campaignsHandler.Create)
	apiCampaigns.GET("", campaignsHandler.GetMany)
	apiCampaigns.GET("/:campaignID", campaignsHandler.GetOne)
	apiCampaigns.DELETE("/:campaignID", campaignsHandler.DeleteOne)
	apiCampaigns.PUT("/:campaignID", campaignsHandler.UpdateOne)

	apiTime := router.Group("/time")
	timeHandler := handler.NewTimeHandler(actCacheSvc)
	apiTime.POST("/advance", timeHandler.Set)
	apiTime.GET("/advance", timeHandler.Get)

	apiAds := router.Group("/ads")
	adsHandler := handler.NewAdsHandler(clientSvc, actCacheSvc, campaignSvc, adsHistorySvc, mlScoreSvc)
	apiAds.GET("", adsHandler.GetOne)
	apiAds.POST("/:campaignID/click", adsHandler.Click)

	apiStats := router.Group("/stats")
	statsHandler := handler.NewStatsHandler(advertiserSvc, campaignSvc, adsHistorySvc, actCacheSvc)
	apiStats.GET("/campaigns/:campaignID", statsHandler.GetCampaignStats)
	apiStats.GET("/campaigns/:campaignID/daily", statsHandler.GetCampaignDailyStats)
	apiStats.GET("/advertisers/:advertiserID/campaigns", statsHandler.GetAdvertiserStats)
	apiStats.GET("/advertisers/:advertiserID/campaigns/daily", statsHandler.GetAdvertiserDailyStats)

	apiAI := router.Group("/ai")
	aiHandler := handler.NewAIHandler(gigaChatSvc, actCacheSvc)
	apiAI.POST("/text/generate", aiHandler.GenerateText)
	apiAI.POST("/text/moderation", aiHandler.SetModeration)

	return router
}
