package service

import (
	"github.com/google/uuid"
	"neJok/solution/internal/model"
	"neJok/solution/internal/repository"
)

type ClientService struct {
	clientRepo *repository.ClientRepo
}

func NewClientService(clientRepo *repository.ClientRepo) *ClientService {
	return &ClientService{clientRepo}
}

func (s *ClientService) CreateOrUpdate(req []model.Client) error {
	return s.clientRepo.Add(req)
}

func (s *ClientService) GetByID(clientID uuid.UUID) (model.Client, error) {
	return s.clientRepo.Get(clientID)
}
