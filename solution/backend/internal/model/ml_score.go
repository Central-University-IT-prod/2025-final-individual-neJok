package model

import (
	"github.com/google/uuid"
)

type MLScore struct {
	ClientID     uuid.UUID `bson:"client_id" json:"client_id" binding:"required" format:"uuid"`
	AdvertiserID uuid.UUID `bson:"advertiser_id" json:"advertiser_id" binding:"required" format:"uuid"`
	Score        *int      `bson:"score" json:"score" binding:"required,min=0"`
}
