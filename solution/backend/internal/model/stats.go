package model

type CampaignStats struct {
	TotalViews        int64   `json:"impressions_count" validate:"required"`
	TotalClicks       int64   `json:"clicks_count" validate:"required"`
	ConversionRate    float64 `json:"conversion" validate:"required"`
	TotalProfitViews  float64 `json:"spent_impressions" validate:"required"`
	TotalProfitClicks float64 `json:"spent_clicks" validate:"required"`
	TotalSpent        float64 `json:"spent_total" validate:"required"`
}

type CampaignStatsDaily struct {
	CampaignStats
	Date int32 `json:"date" bson:"date" validate:"required"`
}
