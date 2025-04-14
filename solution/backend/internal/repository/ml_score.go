package repository

import (
	"context"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"neJok/solution/internal/model"
)

type MLScoreRepo struct {
	db *mongo.Database
}

func NewMLScoreRepo(db *mongo.Database) *MLScoreRepo {
	return &MLScoreRepo{db}
}

func (r *MLScoreRepo) Add(mlScore model.MLScore) error {
	collection := r.db.Collection("ml_scores")
	ctx := context.TODO()

	filter := bson.M{"client_id": mlScore.ClientID, "advertiser_id": mlScore.AdvertiserID}
	update := bson.M{"$set": bson.M{"score": mlScore.Score}}

	upsert := true
	updateOptions := options.UpdateOptions{
		Upsert: &upsert,
	}
	_, err := collection.UpdateOne(ctx, filter, update, &updateOptions)
	if err != nil {
		log.Printf("Failed to insert mlScore: %v", err)
		return err
	}
	return err
}

func (r *MLScoreRepo) GetMax(clientID uuid.UUID) (int, error) {
	collection := r.db.Collection("ml_scores")
	ctx := context.TODO()

	filter := bson.M{"client_id": clientID}
	findOptions := options.FindOneOptions{
		Sort: bson.D{{"score", -1}},
	}
	var mlScore model.MLScore
	err := collection.FindOne(ctx, filter, &findOptions).Decode(&mlScore)
	if err != nil {
		return 0, err
	}
	return *mlScore.Score, nil
}

func (r *MLScoreRepo) GetAll() ([]model.MLScore, error) {
	collection := r.db.Collection("ml_scores")
	ctx := context.TODO()

	filter := bson.M{}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		log.Printf("Failed to retrieve mlScores: %v", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var mlScores []model.MLScore
	for cursor.Next(ctx) {
		var mlScore model.MLScore
		if err := cursor.Decode(&mlScore); err != nil {
			log.Printf("Failed to decode mlScore: %v", err)
			return nil, err
		}
		mlScores = append(mlScores, mlScore)
	}

	if err := cursor.Err(); err != nil {
		log.Printf("Cursor error: %v", err)
		return nil, err
	}

	return mlScores, nil
}
