package service

import (
	"neJok/solution/internal/model"
	"neJok/solution/internal/repository"
	"time"
)

type ActCacheService struct {
	actCacheRepo *repository.ActCacheRepo
}

func NewActCacheService(actCacheRepo *repository.ActCacheRepo) *ActCacheService {
	return &ActCacheService{actCacheRepo}
}

func (s *ActCacheService) GetInt(key string) (int, error) {
	return s.actCacheRepo.GetInt(key)
}

func (s *ActCacheService) SetInt(key string, value int) error {
	return s.actCacheRepo.SetInt(key, value)
}

func (s *ActCacheService) GetList(key string) ([]model.CampaignForUser, error) {
	return s.actCacheRepo.GetList(key)
}

func (s *ActCacheService) SetList(key string, campaigns []model.CampaignForUser) error {
	return s.actCacheRepo.SetList(key, campaigns)
}

func (s *ActCacheService) DeleteKeysByPrefix(key string) error {
	return s.actCacheRepo.DeleteKeysByPrefix(key)
}

func (s *ActCacheService) GetStr(key string) (string, error) {
	return s.actCacheRepo.GetStr(key)
}

func (s *ActCacheService) SetStr(key string, value string, expire *time.Duration) error {
	return s.actCacheRepo.SetStr(key, value, expire)
}
