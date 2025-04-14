package service

import (
	"github.com/google/uuid"
	"neJok/solution/internal/model"
	"neJok/solution/internal/repository"
)

type MLScoreService struct {
	mlScoreRepo *repository.MLScoreRepo
}

func NewMLScoreService(mlScoreRepo *repository.MLScoreRepo) *MLScoreService {
	return &MLScoreService{mlScoreRepo}
}

func (s *MLScoreService) CreateOrUpdate(req model.MLScore) error {
	return s.mlScoreRepo.Add(req)
}

func (s *MLScoreService) GetMax(clientID uuid.UUID) (int, error) {
	return s.mlScoreRepo.GetMax(clientID)
}

func (s *MLScoreService) GetAll() ([]model.MLScore, error) {
	return s.mlScoreRepo.GetAll()
}
