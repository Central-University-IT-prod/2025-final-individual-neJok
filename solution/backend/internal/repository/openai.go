package repository

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"neJok/solution/config"
	"net/http"
)

type OpenAIRepo struct {
	cfg *config.Config
}

func NewOpenAIRepo(cfg *config.Config) *OpenAIRepo {
	return &OpenAIRepo{cfg}
}

func (r *OpenAIRepo) ModerateText(text string) (bool, error) {
	url := r.cfg.OpenAIBaseUrl + "/v1/moderations"

	body := map[string]interface{}{
		"model": "omni-moderation-latest",
		"input": text,
	}

	requestBody, err := json.Marshal(body)
	if err != nil {
		return false, fmt.Errorf("ошибка создания тела запроса: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return false, fmt.Errorf("ошибка создания запроса: %v", err)
	}

	req.Header.Add("Authorization", "Bearer "+r.cfg.OpenAIToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, fmt.Errorf("ошибка выполнения запроса: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("ошибка - статус код %v", resp.StatusCode)
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, fmt.Errorf("ошибка чтения тела ответа: %v", err)
	}

	var response struct {
		Results []struct {
			Flagged bool `json:"flagged"`
		} `json:"results"`
	}

	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		return false, fmt.Errorf("ошибка парсинга ответа: %v", err)
	}

	for _, result := range response.Results {
		if result.Flagged {
			return false, nil
		}
	}
	return true, nil
}
