package model

import (
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AdsHistory struct {
	CampaignID   uuid.UUID            `bson:"campaign_id"`
	AdvertiserID uuid.UUID            `bson:"advertiser_id"`
	ClientID     uuid.UUID            `bson:"client_id"`
	Type         string               `bson:"type"`
	Profit       primitive.Decimal128 `bson:"profit"`
	Date         int                  `bson:"date"`
}
