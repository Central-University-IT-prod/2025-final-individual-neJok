package model

import "github.com/google/uuid"

type Advertiser struct {
	AdvertiserID uuid.UUID `bson:"_id" json:"advertiser_id" binding:"required" format:"uuid"`
	Name         string    `bson:"name" json:"name" binding:"required"`
	CreatedAt    int32     `bson:"created_at" json:"-"`
}

type AdvertiserCreateOrUpdateRequest []Advertiser
