package service

import (
	"neJok/solution/internal/model"
	"neJok/solution/internal/repository"
	"time"
)

type GigaChatService struct {
	gigaChatRepo *repository.GigaChatRepo
	actCacheRepo *repository.ActCacheRepo
}

func NewGigaChatService(gigaChatRepo *repository.GigaChatRepo, actCacheRepo *repository.ActCacheRepo) *GigaChatService {
	return &GigaChatService{gigaChatRepo, actCacheRepo}
}

func (s *GigaChatService) GetToken() (string, error) {
	var token string
	token, err := s.actCacheRepo.GetStr("giga_chat_token")
	if err == nil && token != "" {
		return token, nil
	}
	token, expire, err := s.gigaChatRepo.GetToken()
	if err != nil {
		return "", err
	}

	duration := expire.Sub(time.Now())
	s.actCacheRepo.SetStr("giga_chat_token", token, &duration)
	return token, nil
}

func (s *GigaChatService) GenerateText(req model.GenerateTextRequest) (model.GenerateTextResponse, int) {
	apiToken, err := s.GetToken()
	if err != nil {
		return model.GenerateTextResponse{
			Message: err.Error(),
		}, 500
	}
	userPrompt := "Сгенерируй мне текст для: " + req.AdTitle + "\n\nПожелания: " + *req.Wishes
	if req.Location != nil {
		userPrompt += "\nЛокация: " + *req.Location
	}
	if req.Gender != "" {
		userPrompt += "\nПол: " + req.Gender
	}
	answer, err := s.gigaChatRepo.GenerateText(userPrompt, apiToken)
	if err != nil {
		return model.GenerateTextResponse{
			Message: err.Error(),
		}, 500
	}
	return answer, 200
}
