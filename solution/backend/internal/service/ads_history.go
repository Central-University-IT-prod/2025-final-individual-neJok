package service

import (
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"neJok/solution/internal/model"
	"neJok/solution/internal/repository"
	"strconv"
)

type AdsHistoryService struct {
	adsHistoryRepo *repository.AdsHistoryRepo
}

func NewAdsHistoryService(adsHistoryRepo *repository.AdsHistoryRepo) *AdsHistoryService {
	return &AdsHistoryService{adsHistoryRepo}
}

func (s *AdsHistoryService) Add(campaignID uuid.UUID, advertiserID uuid.UUID, clientID uuid.UUID, cost float64, historyType string, date int) error {
	profit, err := primitive.ParseDecimal128(strconv.FormatFloat(cost, 'g', -1, 64))
	if err != nil {
		return err
	}
	adsHistory := model.AdsHistory{
		CampaignID:   campaignID,
		AdvertiserID: advertiserID,
		ClientID:     clientID,
		Profit:       profit,
		Type:         historyType,
		Date:         date,
	}
	return s.adsHistoryRepo.Add(adsHistory)
}

func (s *AdsHistoryService) GetOne(campaignID uuid.UUID, clientID uuid.UUID, campaignType string) (model.AdsHistory, error) {
	return s.adsHistoryRepo.GetOne(campaignID, clientID, campaignType)
}

func (s *AdsHistoryService) GetAggregatedCampaignStats(campaignID uuid.UUID) (model.CampaignStats, error) {
	return s.adsHistoryRepo.GetAggregatedCampaignStats(campaignID)
}

func (s *AdsHistoryService) GetAggregatedAdvertiserStats(advertiserID uuid.UUID) (model.CampaignStats, error) {
	return s.adsHistoryRepo.GetAggregatedAdvertiserStats(advertiserID)
}

func (s *AdsHistoryService) GetAggregatedCampaignDailyStats(campaignID uuid.UUID, startDate int32, endDate int32) ([]model.CampaignStatsDaily, error) {
	return s.adsHistoryRepo.GetAggregatedCampaignDailyStats(campaignID, startDate, endDate)
}

func (s *AdsHistoryService) GetAggregatedAdvertiserDailyStats(advertiserID uuid.UUID, startDate int32, endDate int32) ([]model.CampaignStatsDaily, error) {
	return s.adsHistoryRepo.GetAggregatedAdvertiserDailyStats(advertiserID, startDate, endDate)
}

func (s *AdsHistoryService) GetViewsAndClicks(campaignID uuid.UUID) (int64, int64, error) {
	return s.adsHistoryRepo.GetViewsAndClicks(campaignID)
}
