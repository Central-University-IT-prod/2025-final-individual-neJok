package model

import "github.com/google/uuid"

type AdsRequest struct {
	ClientID string `form:"client_id" json:"client_id" binding:"required" format:"uuid"`
}

type AdsResponse struct {
	CampaignID   uuid.UUID `bson:"_id" json:"ad_id" validate:"required" format:"uuid"`
	AdvertiserID uuid.UUID `bson:"advertiser_id" json:"advertiser_id" validate:"required" format:"uuid"`
	AdTitle      string    `bson:"ad_title" json:"ad_title" validate:"required"`
	AdText       string    `bson:"ad_text" json:"ad_text" validate:"required"`
	ImageURL     *string   `bson:"image_url" json:"image_url,omitempty" format:"url" example:"https://domain.com/image.jpg"`
}
