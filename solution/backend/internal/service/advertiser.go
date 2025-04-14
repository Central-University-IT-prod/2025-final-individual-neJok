package service

import (
	"github.com/google/uuid"
	"neJok/solution/internal/model"
	"neJok/solution/internal/repository"
)

type AdvertiserService struct {
	advertiserRepo *repository.AdvertiserRepo
}

func NewAdvertiserService(advertiserRepo *repository.AdvertiserRepo) *AdvertiserService {
	return &AdvertiserService{advertiserRepo}
}

func (s *AdvertiserService) CreateOrUpdate(req []model.Advertiser, currentDate int32) error {
	return s.advertiserRepo.Add(req, currentDate)
}

func (s *AdvertiserService) GetByID(advertiserID uuid.UUID) (model.Advertiser, error) {
	return s.advertiserRepo.Get(advertiserID)
}
