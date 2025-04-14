package service

import (
	"neJok/solution/internal/repository"
)

type OpenAIService struct {
	openAIRepo *repository.OpenAIRepo
}

func NewOpenAIService(openAIRepo *repository.OpenAIRepo) *OpenAIService {
	return &OpenAIService{openAIRepo}
}

func (s *OpenAIService) ModerateText(text string) (bool, error) {
	return s.openAIRepo.ModerateText(text)
}
