package repository

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"neJok/solution/internal/model"
)

type CampaignRepo struct {
	db *mongo.Database
}

func NewCampaignRepo(db *mongo.Database) *CampaignRepo {
	return &CampaignRepo{db}
}

func (r *CampaignRepo) Add(campaign model.Campaign) error {
	collection := r.db.Collection("campaigns")
	ctx := context.TODO()
	_, err := collection.InsertOne(ctx, campaign)
	if err != nil {
		log.Printf("Failed to insert campaign: %v", err)
		return err
	}
	return err
}

func (r *CampaignRepo) GetMany(advertiserID uuid.UUID, limit int, offset int) (total int64, arr []model.Campaign, err error) {
	collection := r.db.Collection("campaigns")
	filter := bson.M{"advertiser_id": advertiserID}
	total, err = collection.CountDocuments(context.Background(), filter)
	if err != nil {
		return 0, nil, fmt.Errorf("error counting documents: %w", err)
	}

	findOptions := options.Find().
		SetSort(bson.M{"created_at": -1}).
		SetSkip(int64(offset)).
		SetLimit(int64(limit))

	cursor, err := collection.Find(context.Background(), filter, findOptions)
	if err != nil {
		return 0, nil, fmt.Errorf("error finding documents: %w", err)
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var campaign model.Campaign
		if err := cursor.Decode(&campaign); err != nil {
			return 0, nil, fmt.Errorf("error decoding document: %w", err)
		}
		arr = append(arr, campaign)
	}

	if err := cursor.Err(); err != nil {
		return 0, nil, fmt.Errorf("cursor error: %w", err)
	}

	return total, arr, nil
}

func (r *CampaignRepo) Get(campaignID uuid.UUID) (model.Campaign, error) {
	collection := r.db.Collection("campaigns")
	ctx := context.TODO()

	var campaign model.Campaign
	filter := bson.M{"_id": campaignID}
	err := collection.FindOne(ctx, filter).Decode(&campaign)
	return campaign, err
}

func (r *CampaignRepo) Delete(advertiserID uuid.UUID, campaignID uuid.UUID) (model.Campaign, error) {
	collection := r.db.Collection("campaigns")
	ctx := context.TODO()

	var campaign model.Campaign
	filter := bson.M{"_id": campaignID, "advertiser_id": advertiserID}
	err := collection.FindOneAndDelete(ctx, filter).Decode(&campaign)
	if err != nil {
		return model.Campaign{}, err
	}

	return campaign, nil
}

func (r *CampaignRepo) Update(advertiserID uuid.UUID, campaignID uuid.UUID, campaign model.CampaignUpdate) (model.Campaign, error) {
	collection := r.db.Collection("campaigns")
	ctx := context.TODO()

	filter := bson.M{"_id": campaignID, "advertiser_id": advertiserID}

	if campaign.Targeting == nil {
		campaign.Targeting = &model.CampaignTargeting{}
	}

	if campaign.Targeting.Gender == "" {
		campaign.Targeting.Gender = "ALL"
	}

	res := collection.FindOneAndUpdate(ctx, filter, bson.M{"$set": campaign}, options.FindOneAndUpdate().SetReturnDocument(options.After))

	if res.Err() != nil {
		return model.Campaign{}, res.Err()
	}

	var updatedCampaign model.Campaign
	if err := res.Decode(&updatedCampaign); err != nil {
		return model.Campaign{}, err
	}

	return updatedCampaign, nil
}

func (r *CampaignRepo) GetManyByTargeting(gender string, age int, location string, currentDate int, clientID uuid.UUID) (arr []model.CampaignForUser, stats model.CampaignsDBStats, err error) {
	collection := r.db.Collection("campaigns")
	filter := bson.M{
		"targeting.gender": bson.M{"$in": []string{"ALL", gender}},
		"$and": []bson.M{
			{
				"$or": []bson.M{
					{"targeting.age_from": nil},
					{"targeting.age_from": bson.M{"$lte": age}},
				},
			},
			{
				"$or": []bson.M{
					{"targeting.age_to": nil},
					{"targeting.age_to": bson.M{"$gte": age}},
				},
			},
			{
				"$or": []bson.M{
					{"targeting.location": nil},
					{"targeting.location": location},
				},
			},
		},
		"start_date": bson.M{"$lte": currentDate},
		"end_date":   bson.M{"$gte": currentDate},
	}

	pipeline := []bson.M{
		{
			"$match": filter,
		},
		{
			"$lookup": bson.M{
				"from":         "ads_history",
				"localField":   "_id",
				"foreignField": "campaign_id",
				"as":           "ad_history",
				"pipeline": []bson.M{
					{
						"$match": bson.M{
							"client_id": clientID,
						},
					},
				},
			},
		},
		{
			"$addFields": bson.M{
				"is_clicked_by_user": bson.M{
					"$in": []interface{}{"click", "$ad_history.type"},
				},
				"views_count": bson.M{
					"$size": bson.M{
						"$filter": bson.M{
							"input": "$ad_history",
							"as":    "history",
							"cond":  bson.M{"$eq": []interface{}{"$$history.type", "view"}},
						},
					},
				},
			},
		},
		{
			"$project": bson.M{
				"_id":                 1,
				"advertiser_id":       1,
				"ad_title":            1,
				"ad_text":             1,
				"cost_per_impression": 1,
				"cost_per_click":      1,
				"end_date":            1,
				"is_clicked_by_user":  1,
				"views_count":         1,
				"impressions_limit":   1,
				"clicks_limit":        1,
				"image_url":           1,
			},
		},
	}

	cursor, err := collection.Aggregate(context.Background(), pipeline)
	if err != nil {
		return nil, model.CampaignsDBStats{}, fmt.Errorf("error finding documents: %w", err)
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var campaign model.CampaignForUser
		if err := cursor.Decode(&campaign); err != nil {
			return nil, model.CampaignsDBStats{}, fmt.Errorf("error decoding document: %w", err)
		}
		if campaign.CostPerImpression > stats.MaxCostPerImpression {
			stats.MaxCostPerImpression = campaign.CostPerImpression
		}
		if campaign.CostPerClick > stats.MaxCostPerClick {
			stats.MaxCostPerClick = campaign.CostPerClick
		}
		if campaign.EndDate > stats.MaxEndDate {
			stats.MaxEndDate = campaign.EndDate
		}
		if campaign.EndDate < stats.MinEndDate {
			stats.MinEndDate = campaign.EndDate
		}
		if campaign.IsClickedByUser == true {
			continue
		}
		arr = append(arr, campaign)
	}

	if err := cursor.Err(); err != nil {
		return nil, model.CampaignsDBStats{}, fmt.Errorf("cursor error: %w", err)
	}

	return arr, stats, nil
}
