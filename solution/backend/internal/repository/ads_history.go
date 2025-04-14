package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"neJok/solution/internal/model"
)

type AdsHistoryRepo struct {
	db *mongo.Database
}

func NewAdsHistoryRepo(db *mongo.Database) *AdsHistoryRepo {
	return &AdsHistoryRepo{db}
}

func (r *AdsHistoryRepo) Add(adsHistory model.AdsHistory) error {
	collection := r.db.Collection("ads_history")
	ctx := context.TODO()

	_, err := collection.InsertOne(ctx, adsHistory)
	if err != nil {
		log.Printf("Failed to insert ads history: %v", err)
		return err
	}

	return nil
}
func (r *AdsHistoryRepo) GetOne(campaignID uuid.UUID, clientID uuid.UUID, campaignType string) (model.AdsHistory, error) {
	collection := r.db.Collection("ads_history")
	ctx := context.TODO()

	var adsHistory model.AdsHistory
	filter := bson.M{"campaign_id": campaignID, "client_id": clientID, "type": campaignType}
	err := collection.FindOne(ctx, filter).Decode(&adsHistory)
	return adsHistory, err
}

func (r *AdsHistoryRepo) GetAggregatedCampaignStats(campaignID uuid.UUID) (model.CampaignStats, error) {
	collection := r.db.Collection("ads_history")
	ctx := context.TODO()
	var stats model.CampaignStats

	filter := bson.M{"campaign_id": campaignID}
	pipeline := mongo.Pipeline{
		{{"$match", filter}},
		{{"$group", bson.M{
			"_id":               nil,
			"totalViews":        bson.M{"$sum": bson.M{"$cond": bson.M{"if": bson.M{"$eq": []interface{}{"$type", "view"}}, "then": 1, "else": 0}}},
			"totalClicks":       bson.M{"$sum": bson.M{"$cond": bson.M{"if": bson.M{"$eq": []interface{}{"$type", "click"}}, "then": 1, "else": 0}}},
			"totalProfitViews":  bson.M{"$sum": bson.M{"$cond": bson.M{"if": bson.M{"$eq": []interface{}{"$type", "view"}}, "then": "$profit", "else": bson.M{"$toDecimal": 0}}}},
			"totalProfitClicks": bson.M{"$sum": bson.M{"$cond": bson.M{"if": bson.M{"$eq": []interface{}{"$type", "click"}}, "then": "$profit", "else": bson.M{"$toDecimal": 0}}}},
		}}},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		log.Printf("Failed to aggregate ads history: %v", err)
		return stats, err
	}
	defer cursor.Close(ctx)

	if cursor.Next(ctx) {
		var result struct {
			TotalViews        int64                `bson:"totalViews"`
			TotalClicks       int64                `bson:"totalClicks"`
			TotalProfitViews  primitive.Decimal128 `bson:"totalProfitViews"`
			TotalProfitClicks primitive.Decimal128 `bson:"totalProfitClicks"`
		}

		if err := cursor.Decode(&result); err != nil {
			log.Printf("Failed to decode aggregated result: %v", err)
			return stats, err
		}

		totalViews := result.TotalViews
		totalClicks := result.TotalClicks
		totalProfitViews := result.TotalProfitViews
		totalProfitClicks := result.TotalProfitClicks

		totalProfitViewsDecimal, err := decimal.NewFromString(totalProfitViews.String())
		if err != nil {
			return stats, err
		}
		totalProfitClicksDecimal, err := decimal.NewFromString(totalProfitClicks.String())
		if err != nil {
			return stats, err
		}

		conversionRate := decimal.NewFromInt(0)
		if totalViews > 0 {
			conversionRate = decimal.NewFromInt(totalClicks).Div(decimal.NewFromInt(totalViews)).Mul(decimal.NewFromInt(100))
		}

		stats = model.CampaignStats{
			TotalViews:        totalViews,
			TotalClicks:       totalClicks,
			ConversionRate:    conversionRate.InexactFloat64(),
			TotalProfitViews:  totalProfitViewsDecimal.InexactFloat64(),
			TotalProfitClicks: totalProfitClicksDecimal.InexactFloat64(),
			TotalSpent:        totalProfitViewsDecimal.Add(totalProfitClicksDecimal).InexactFloat64(),
		}
	}

	return stats, nil
}

func (r *AdsHistoryRepo) GetAggregatedAdvertiserStats(advertiserID uuid.UUID) (model.CampaignStats, error) {
	collection := r.db.Collection("ads_history")
	ctx := context.TODO()
	var stats model.CampaignStats

	filter := bson.M{"advertiser_id": advertiserID}
	pipeline := mongo.Pipeline{
		{{"$match", filter}},
		{{"$group", bson.M{
			"_id":               nil,
			"totalViews":        bson.M{"$sum": bson.M{"$cond": bson.M{"if": bson.M{"$eq": []interface{}{"$type", "view"}}, "then": 1, "else": 0}}},
			"totalClicks":       bson.M{"$sum": bson.M{"$cond": bson.M{"if": bson.M{"$eq": []interface{}{"$type", "click"}}, "then": 1, "else": 0}}},
			"totalProfitViews":  bson.M{"$sum": bson.M{"$cond": bson.M{"if": bson.M{"$eq": []interface{}{"$type", "view"}}, "then": "$profit", "else": bson.M{"$toDecimal": 0}}}},
			"totalProfitClicks": bson.M{"$sum": bson.M{"$cond": bson.M{"if": bson.M{"$eq": []interface{}{"$type", "click"}}, "then": "$profit", "else": bson.M{"$toDecimal": 0}}}},
		}}},
	}
	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		log.Printf("Failed to aggregate ads history: %v", err)
		return stats, err
	}
	defer cursor.Close(ctx)

	if cursor.Next(ctx) {
		var result struct {
			TotalViews        int64                `bson:"totalViews"`
			TotalClicks       int64                `bson:"totalClicks"`
			TotalProfitViews  primitive.Decimal128 `bson:"totalProfitViews"`
			TotalProfitClicks primitive.Decimal128 `bson:"totalProfitClicks"`
		}

		if err := cursor.Decode(&result); err != nil {
			log.Printf("Failed to decode aggregated result: %v", err)
			return stats, err
		}

		totalProfitViewsDecimal, err := decimal.NewFromString(result.TotalProfitViews.String())
		if err != nil {
			return stats, err
		}
		totalProfitClicksDecimal, err := decimal.NewFromString(result.TotalProfitClicks.String())
		if err != nil {
			return stats, err
		}

		conversionRate := decimal.NewFromInt(0)
		if result.TotalViews > 0 {
			conversionRate = decimal.NewFromInt(result.TotalClicks).Div(decimal.NewFromInt(result.TotalViews)).Mul(decimal.NewFromInt(100))
		}

		stats = model.CampaignStats{
			TotalViews:        result.TotalViews,
			TotalClicks:       result.TotalClicks,
			ConversionRate:    conversionRate.InexactFloat64(),
			TotalProfitViews:  totalProfitViewsDecimal.InexactFloat64(),
			TotalProfitClicks: totalProfitClicksDecimal.InexactFloat64(),
			TotalSpent:        totalProfitViewsDecimal.Add(totalProfitClicksDecimal).InexactFloat64(),
		}
	}

	return stats, nil
}

func (r *AdsHistoryRepo) GetAggregatedCampaignDailyStats(campaignID uuid.UUID, startDate int32, endDate int32) (stats []model.CampaignStatsDaily, err error) {
	collection := r.db.Collection("ads_history")
	ctx := context.TODO()

	filter := bson.M{
		"campaign_id": campaignID,
		"date": bson.M{
			"$gte": startDate,
			"$lte": endDate,
		},
	}

	pipeline := mongo.Pipeline{
		{{"$match", filter}},
		{{"$group", bson.M{
			"_id":               "$date",
			"totalViews":        bson.M{"$sum": bson.M{"$cond": bson.M{"if": bson.M{"$eq": []interface{}{"$type", "view"}}, "then": 1, "else": 0}}},
			"totalClicks":       bson.M{"$sum": bson.M{"$cond": bson.M{"if": bson.M{"$eq": []interface{}{"$type", "click"}}, "then": 1, "else": 0}}},
			"totalProfitViews":  bson.M{"$sum": bson.M{"$cond": bson.M{"if": bson.M{"$eq": []interface{}{"$type", "view"}}, "then": "$profit", "else": bson.M{"$toDecimal": 0}}}},
			"totalProfitClicks": bson.M{"$sum": bson.M{"$cond": bson.M{"if": bson.M{"$eq": []interface{}{"$type", "click"}}, "then": "$profit", "else": bson.M{"$toDecimal": 0}}}},
		}}},
		{{"$sort", bson.M{"_id": 1}}},
	}
	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		log.Printf("Failed to aggregate ads history: %v", err)
		return stats, err
	}
	defer cursor.Close(ctx)

	dailyStatsMap := make(map[int32]model.CampaignStatsDaily)

	for cursor.Next(ctx) {
		var result struct {
			TotalViews        int64                `bson:"totalViews"`
			TotalClicks       int64                `bson:"totalClicks"`
			TotalProfitViews  primitive.Decimal128 `bson:"totalProfitViews"`
			TotalProfitClicks primitive.Decimal128 `bson:"totalProfitClicks"`
			Date              int32                `bson:"_id"`
		}

		if err := cursor.Decode(&result); err != nil {
			log.Printf("Failed to decode aggregated result: %v", err)
			return stats, err
		}

		totalProfitViewsDecimal, err := decimal.NewFromString(result.TotalProfitViews.String())
		if err != nil {
			return stats, err
		}
		totalProfitClicksDecimal, err := decimal.NewFromString(result.TotalProfitClicks.String())
		if err != nil {
			return stats, err
		}

		conversionRate := decimal.NewFromInt(0)
		if result.TotalViews > 0 {
			conversionRate = decimal.NewFromInt(result.TotalClicks).Div(decimal.NewFromInt(result.TotalViews)).Mul(decimal.NewFromInt(100))
		}

		dailyStatsMap[result.Date] = model.CampaignStatsDaily{
			CampaignStats: model.CampaignStats{
				TotalViews:        result.TotalViews,
				TotalClicks:       result.TotalClicks,
				ConversionRate:    conversionRate.InexactFloat64(),
				TotalProfitViews:  totalProfitViewsDecimal.InexactFloat64(),
				TotalProfitClicks: totalProfitClicksDecimal.InexactFloat64(),
				TotalSpent:        totalProfitViewsDecimal.Add(totalProfitClicksDecimal).InexactFloat64(),
			},
			Date: result.Date,
		}
	}
	for date := startDate; date <= endDate; date++ {
		if _, exists := dailyStatsMap[date]; !exists {
			dailyStatsMap[date] = model.CampaignStatsDaily{
				CampaignStats: model.CampaignStats{
					TotalViews:        0,
					TotalClicks:       0,
					ConversionRate:    0,
					TotalProfitViews:  0,
					TotalProfitClicks: 0,
					TotalSpent:        0,
				},
				Date: date,
			}
		}
	}
	for date := startDate; date <= endDate; date++ {
		stats = append(stats, dailyStatsMap[date])
	}

	return stats, nil
}

func (r *AdsHistoryRepo) GetAggregatedAdvertiserDailyStats(advertiserID uuid.UUID, startDate int32, endDate int32) (stats []model.CampaignStatsDaily, err error) {
	collection := r.db.Collection("ads_history")
	ctx := context.TODO()

	filter := bson.M{
		"advertiser_id": advertiserID,
		"date": bson.M{
			"$gte": startDate,
			"$lte": endDate,
		},
	}

	pipeline := mongo.Pipeline{
		{{"$match", filter}},
		{{"$group", bson.M{
			"_id":               "$date",
			"totalViews":        bson.M{"$sum": bson.M{"$cond": bson.M{"if": bson.M{"$eq": []interface{}{"$type", "view"}}, "then": 1, "else": 0}}},
			"totalClicks":       bson.M{"$sum": bson.M{"$cond": bson.M{"if": bson.M{"$eq": []interface{}{"$type", "click"}}, "then": 1, "else": 0}}},
			"totalProfitViews":  bson.M{"$sum": bson.M{"$cond": bson.M{"if": bson.M{"$eq": []interface{}{"$type", "view"}}, "then": "$profit", "else": bson.M{"$toDecimal": 0}}}},
			"totalProfitClicks": bson.M{"$sum": bson.M{"$cond": bson.M{"if": bson.M{"$eq": []interface{}{"$type", "click"}}, "then": "$profit", "else": bson.M{"$toDecimal": 0}}}},
		}}},
		{{"$sort", bson.M{"_id": 1}}},
	}
	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		log.Printf("Failed to aggregate ads history: %v", err)
		return stats, err
	}
	defer cursor.Close(ctx)

	dailyStatsMap := make(map[int32]model.CampaignStatsDaily)

	for cursor.Next(ctx) {
		var result struct {
			TotalViews        int64                `bson:"totalViews"`
			TotalClicks       int64                `bson:"totalClicks"`
			TotalProfitViews  primitive.Decimal128 `bson:"totalProfitViews"`
			TotalProfitClicks primitive.Decimal128 `bson:"totalProfitClicks"`
			Date              int32                `bson:"_id"`
		}

		if err := cursor.Decode(&result); err != nil {
			log.Printf("Failed to decode aggregated result: %v", err)
			return stats, err
		}

		totalProfitViewsDecimal, err := decimal.NewFromString(result.TotalProfitViews.String())
		if err != nil {
			return stats, err
		}
		totalProfitClicksDecimal, err := decimal.NewFromString(result.TotalProfitClicks.String())
		if err != nil {
			return stats, err
		}

		conversionRate := decimal.NewFromInt(0)
		if result.TotalViews > 0 {
			conversionRate = decimal.NewFromInt(result.TotalClicks).Div(decimal.NewFromInt(result.TotalViews)).Mul(decimal.NewFromInt(100))
		}

		dailyStatsMap[result.Date] = model.CampaignStatsDaily{
			CampaignStats: model.CampaignStats{
				TotalViews:        result.TotalViews,
				TotalClicks:       result.TotalClicks,
				ConversionRate:    conversionRate.InexactFloat64(),
				TotalProfitViews:  totalProfitViewsDecimal.InexactFloat64(),
				TotalProfitClicks: totalProfitClicksDecimal.InexactFloat64(),
				TotalSpent:        totalProfitViewsDecimal.Add(totalProfitClicksDecimal).InexactFloat64(),
			},
			Date: result.Date,
		}
	}
	for date := startDate; date <= endDate; date++ {
		if _, exists := dailyStatsMap[date]; !exists {
			dailyStatsMap[date] = model.CampaignStatsDaily{
				CampaignStats: model.CampaignStats{
					TotalViews:        0,
					TotalClicks:       0,
					ConversionRate:    0,
					TotalProfitViews:  0,
					TotalProfitClicks: 0,
					TotalSpent:        0,
				},
				Date: date,
			}
		}
	}
	for date := startDate; date <= endDate; date++ {
		stats = append(stats, dailyStatsMap[date])
	}

	return stats, nil
}

func (r *AdsHistoryRepo) GetViewsAndClicks(campaignID uuid.UUID) (clicks int64, views int64, err error) {
	collection := r.db.Collection("ads_history")
	ctx := context.TODO()

	filter := bson.M{"campaign_id": campaignID}
	pipeline := mongo.Pipeline{
		{{"$match", filter}},
		{{"$group", bson.M{
			"_id":         nil,
			"totalViews":  bson.M{"$sum": bson.M{"$cond": bson.M{"if": bson.M{"$eq": []interface{}{"$type", "view"}}, "then": 1, "else": 0}}},
			"totalClicks": bson.M{"$sum": bson.M{"$cond": bson.M{"if": bson.M{"$eq": []interface{}{"$type", "click"}}, "then": 1, "else": 0}}},
		}}},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		log.Printf("Failed to aggregate clicks and views: %v", err)
		return 0, 0, err
	}
	defer cursor.Close(ctx)

	if cursor.Next(ctx) {
		var result struct {
			TotalViews  int64 `bson:"totalViews"`
			TotalClicks int64 `bson:"totalClicks"`
		}
		if err := cursor.Decode(&result); err != nil {
			log.Printf("Failed to decode result: %v", err)
			return 0, 0, err
		}

		views = result.TotalViews
		clicks = result.TotalClicks
	}

	return views, clicks, nil
}
