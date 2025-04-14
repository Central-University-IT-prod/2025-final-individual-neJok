package repository

import (
	"context"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"neJok/solution/internal/model"
)

type AdvertiserRepo struct {
	db *mongo.Database
}

func NewAdvertiserRepo(db *mongo.Database) *AdvertiserRepo {
	return &AdvertiserRepo{db}
}

func (r *AdvertiserRepo) Add(advertisers []model.Advertiser, currentDate int32) error {
	collection := r.db.Collection("advertisers")
	ctx := context.TODO()

	var bulkOps []mongo.WriteModel

	for _, advertiser := range advertisers {
		filter := bson.M{"_id": advertiser.AdvertiserID}
		update := bson.M{
			"$set": bson.M{
				"name": advertiser.Name,
			},
			"$setOnInsert": bson.M{
				"created_at": currentDate,
			},
		}

		bulkOps = append(bulkOps, mongo.NewUpdateOneModel().
			SetFilter(filter).
			SetUpdate(update).
			SetUpsert(true),
		)
	}

	if len(bulkOps) > 0 {
		_, err := collection.BulkWrite(ctx, bulkOps)
		if err != nil {
			log.Printf("Failed to bulk insert/update advertisers: %v", err)
			return err
		}
	}

	return nil
}

func (r *AdvertiserRepo) Get(advertiserID uuid.UUID) (model.Advertiser, error) {
	collection := r.db.Collection("advertisers")
	ctx := context.TODO()

	var advertiser model.Advertiser
	filter := bson.M{"_id": advertiserID}
	err := collection.FindOne(ctx, filter).Decode(&advertiser)
	return advertiser, err
}
