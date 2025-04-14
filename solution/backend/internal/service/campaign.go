package service

import (
	"github.com/google/uuid"
	"neJok/solution/internal/model"
	"neJok/solution/internal/repository"
	"time"
)

type CampaignService struct {
	campaignRepo *repository.CampaignRepo
}

func NewCampaignService(campaignRepo *repository.CampaignRepo) *CampaignService {
	return &CampaignService{campaignRepo}
}

func (s *CampaignService) Add(advertiserID uuid.UUID, req model.CampaignCreate, campaignID uuid.UUID) (model.Campaign, error) {
	campaign := model.Campaign{
		CampaignCreate: req,
		AdvertiserID:   advertiserID,
		CreatedAt:      time.Now(),
		CampaignID:     campaignID,
	}

	return campaign, s.campaignRepo.Add(campaign)
}

func (s *CampaignService) GetMany(advertiserID uuid.UUID, limit int, offset int) (int64, []model.Campaign, error) {
	return s.campaignRepo.GetMany(advertiserID, limit, offset)
}

func (s *CampaignService) GetByID(campaignID uuid.UUID) (model.Campaign, error) {
	return s.campaignRepo.Get(campaignID)
}

func (s *CampaignService) DeleteByID(advertiserID uuid.UUID, campaignID uuid.UUID) (model.Campaign, error) {
	return s.campaignRepo.Delete(advertiserID, campaignID)
}

func (s *CampaignService) UpdateByID(advertiserID uuid.UUID, campaignID uuid.UUID, campaign model.CampaignUpdate) (model.Campaign, error) {
	return s.campaignRepo.Update(advertiserID, campaignID, campaign)
}

func (s *CampaignService) GetManyByTargeting(gender string, age int, location string, currentDate int, clientID uuid.UUID) ([]model.CampaignForUser, model.CampaignsDBStats, error) {
	return s.campaignRepo.GetManyByTargeting(gender, age, location, currentDate, clientID)
}
