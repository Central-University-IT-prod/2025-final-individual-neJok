package model

import (
	"github.com/google/uuid"
	"mime/multipart"
	"time"
)

type Campaign struct {
	AdvertiserID   uuid.UUID `bson:"advertiser_id" json:"advertiser_id" validate:"required" format:"uuid"`
	CampaignID     uuid.UUID `bson:"_id" json:"campaign_id" validate:"required" format:"uuid"`
	CreatedAt      time.Time `bson:"created_at" json:"-"`
	CampaignCreate `bson:",inline"`
}

type CampaignTargeting struct {
	Gender   string  `bson:"gender" json:"gender" binding:"omitempty,oneof=MALE FEMALE ALL"`
	AgeFrom  *int    `bson:"age_from" json:"age_from,omitempty" binding:"omitempty,min=0"`
	AgeTo    *int    `bson:"age_to" json:"age_to,omitempty" binding:"omitempty,min=0"`
	Location *string `bson:"location" json:"location,omitempty" binding:"omitempty"`
}

func (c *CampaignTargeting) SetDefaults() {
	if c.Gender == "" {
		c.Gender = "ALL"
	}
}

type CampaignCreate struct {
	ImpressionsLimit  *int                  `bson:"impressions_limit" json:"impressions_limit" binding:"required,min=0,gtefield=ClicksLimit" form:"impressions_limit"`
	ClicksLimit       *int                  `bson:"clicks_limit" json:"clicks_limit" binding:"required,min=0" form:"clicks_limit"`
	CostPerImpression *float64              `bson:"cost_per_impression" json:"cost_per_impression" binding:"required,min=0" form:"cost_per_impression"`
	CostPerClick      *float64              `bson:"cost_per_click" json:"cost_per_click" binding:"required,min=0" form:"cost_per_click"`
	AdTitle           string                `bson:"ad_title" json:"ad_title" binding:"required" form:"ad_title"`
	AdText            string                `bson:"ad_text" json:"ad_text" binding:"required" form:"ad_text"`
	StartDate         *int32                `bson:"start_date" json:"start_date" binding:"required,min=0" form:"start_date"`
	EndDate           *int32                `bson:"end_date" json:"end_date" binding:"required,gtefield=StartDate,min=0" form:"end_date"`
	Targeting         CampaignTargeting     `bson:"targeting" json:"targeting" form:"targeting"`
	ImageURL          *string               `bson:"image_url" json:"image_url,omitempty" form:"image_url" binding:"omitempty,url" format:"url" example:"https://domain.com/image.jpg"`
	ImageFile         *multipart.FileHeader `bson:"-" json:"-" form:"image_file"`
	FileName          *string               `bson:"file_name" json:"-"`
}

type CampaignUpdate struct {
	ImpressionsLimit  *int                  `bson:"impressions_limit" json:"impressions_limit" binding:"required,min=0" form:"impressions_limit"`
	ClicksLimit       *int                  `bson:"clicks_limit" json:"clicks_limit" binding:"required,min=0" form:"clicks_limit"`
	CostPerImpression *float64              `bson:"cost_per_impression" json:"cost_per_impression" binding:"required,min=0" form:"cost_per_impression"`
	CostPerClick      *float64              `bson:"cost_per_click" json:"cost_per_click" binding:"required,min=0" form:"cost_per_click"`
	AdTitle           string                `bson:"ad_title" json:"ad_title" binding:"required" form:"ad_title"`
	AdText            string                `bson:"ad_text" json:"ad_text" binding:"required" form:"ad_text"`
	StartDate         *int32                `bson:"start_date" json:"start_date" binding:"required,min=0" form:"start_date"`
	EndDate           *int32                `bson:"end_date" json:"end_date" binding:"required,min=0" form:"end_date"`
	Targeting         *CampaignTargeting    `bson:"targeting" json:"targeting" binding:"omitempty" form:"targeting"`
	ImageURL          *string               `bson:"image_url" json:"image_url,omitempty" form:"image_url" binding:"omitempty,url" format:"url" example:"https://domain.com/image.jpg"`
	ImageFile         *multipart.FileHeader `bson:"-" json:"-" form:"image_file"`
	FileName          *string               `bson:"file_name" json:"-"`
}

type CampaignGetManyRequest struct {
	Limit  int `form:"size" binding:"min=0"`
	Offset int `form:"page" binding:"min=0"`
}

type CampaignForUser struct {
	CampaignID        uuid.UUID `bson:"_id"`
	AdvertiserID      uuid.UUID `bson:"advertiser_id"`
	ImpressionsLimit  int       `bson:"impressions_limit"`
	ClicksLimit       int       `bson:"clicks_limit"`
	CostPerImpression float64   `bson:"cost_per_impression"`
	CostPerClick      float64   `bson:"cost_per_click"`
	AdTitle           string    `bson:"ad_title"`
	AdText            string    `bson:"ad_text"`
	EndDate           int32     `bson:"end_date"`
	Score             float64   `bson:"score"`
	IsClickedByUser   bool      `bson:"is_clicked_by_user"`
	ViewsCount        int       `bson:"views_count"`
	ImageURL          *string   `bson:"image_url"`
}

type CampaignsDBStats struct {
	MaxCostPerImpression float64 `bson:"max_cost_per_impression" json:"-"`
	MaxCostPerClick      float64 `bson:"max_cost_per_click" json:"-"`
	MinEndDate           int32   `bson:"min_end_date" json:"-"`
	MaxEndDate           int32   `bson:"max_end_date" json:"-"`
	MaxScore             int     `bson:"max_score" json:"-"`
}
