package repository

import (
	"context"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"neJok/solution/internal/model"
)

type ClientRepo struct {
	db *mongo.Database
}

func NewClientRepo(db *mongo.Database) *ClientRepo {
	return &ClientRepo{db}
}

func (r *ClientRepo) Add(clients []model.Client) error {
	collection := r.db.Collection("clients")
	ctx := context.TODO()

	var bulkOps []mongo.WriteModel

	for _, client := range clients {
		filter := bson.M{"_id": client.ClientID}
		update := bson.M{
			"$set": bson.M{
				"login":    client.Login,
				"age":      client.Age,
				"location": client.Location,
				"gender":   client.Gender,
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
			log.Printf("Failed to bulk insert/update clients: %v", err)
			return err
		}
	}

	return nil
}

func (r *ClientRepo) Get(clientID uuid.UUID) (model.Client, error) {
	collection := r.db.Collection("clients")
	ctx := context.TODO()

	var client model.Client
	filter := bson.M{"_id": clientID}
	err := collection.FindOne(ctx, filter).Decode(&client)
	return client, err
}
