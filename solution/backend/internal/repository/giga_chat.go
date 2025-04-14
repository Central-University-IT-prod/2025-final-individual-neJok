package repository

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io"
	"io/ioutil"
	"neJok/solution/config"
	"neJok/solution/internal/model"
	"net/http"
	"time"
)

type GigaChatRepo struct {
	cfg *config.Config
}

func NewGigaChatRepo(cfg *config.Config) *GigaChatRepo {
	return &GigaChatRepo{cfg}
}

func (r *GigaChatRepo) GetToken() (string, time.Time, error) {
	rqUID := uuid.New().String()

	url := r.cfg.GigaChatAuthBaseUrl + "/api/v2/oauth"
	payload := "scope=GIGACHAT_API_PERS"

	req, err := http.NewRequest("POST", url, bytes.NewBufferString(payload))
	if err != nil {
		return "", time.Time{}, fmt.Errorf("ошибка создания запроса: %v", err)
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("RqUID", rqUID)
	req.Header.Add("Authorization", "Basic "+r.cfg.GigaChatAuthKey)

	certPool := x509.NewCertPool()
	certData, err := ioutil.ReadFile("/certs/russian_trusted_root_ca.cer")
	if err != nil {
		return "", time.Time{}, fmt.Errorf("ошибка чтения сертификата: %v", err)
	}

	if ok := certPool.AppendCertsFromPEM(certData); !ok {
		return "", time.Time{}, fmt.Errorf("не удалось добавить сертификат в пул")
	}

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			RootCAs: certPool,
		},
	}

	client := &http.Client{
		Transport: transport,
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("ошибка выполнения запроса: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("ошибка чтения тела ответа: %v", err)
	}

	var tokenResp model.GigaChatTokenResponse
	err = json.Unmarshal(body, &tokenResp)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("ошибка парсинга ответа: %v", err)
	}

	seconds := tokenResp.ExpiresAt / 1000
	nanoseconds := (tokenResp.ExpiresAt % 1000) * 1000000
	expirationTime := time.Unix(seconds, nanoseconds)
	return tokenResp.AccessToken, expirationTime, nil
}

func (r *GigaChatRepo) GenerateText(userPrompt string, apiToken string) (model.GenerateTextResponse, error) {
	url := r.cfg.GigaChatModelsBaseUrl + "/api/v1/chat/completions"

	generatePrompt := "Ты — генератор рекламных текстов. Когда пользователь присылает название товара после текста \"Сгенерируй мне текст для:\", пожелания после текста \"Пожелания: \", а также локацию и пол, ты генерируешь рекламный текст, соответствующий запросу на языке пользователя. В ответе ты пишешь только рекламное сообщение. Если пользователь просит не генерацию рекламного текста, ты говоришь, что ты умеешь только генерировать тексты.\n\nПример:\nСгенерируй мне текст для: Смартфон\nПожелания: Выделить его уникальные функции, привлекательную цену и долгий срок службы.\nЛокация: Челябинск.\n\nОтвет:\n\"Смартфон с уникальными функциями по привлекательной цене в Челябинске! Отличная производительность и долгий срок службы — идеальный выбор для вас!\"\n"
	body := map[string]interface{}{
		"model":              "GigaChat",
		"messages":           []map[string]interface{}{{"role": "system", "content": generatePrompt}, {"role": "user", "content": userPrompt}},
		"temperature":        0.1,
		"top_p":              1.0,
		"stream":             false,
		"max_tokens":         100,
		"repetition_penalty": 1.0,
		"update_interval":    0,
	}

	requestBody, err := json.Marshal(body)
	if err != nil {
		return model.GenerateTextResponse{}, fmt.Errorf("ошибка создания тела запроса: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return model.GenerateTextResponse{}, fmt.Errorf("ошибка создания запроса: %v", err)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+apiToken)
	req.Header.Add("Accept", "application/json")

	certPool := x509.NewCertPool()
	certData, err := ioutil.ReadFile("/certs/russian_trusted_root_ca.cer")
	if err != nil {
		return model.GenerateTextResponse{}, fmt.Errorf("ошибка чтения сертификата: %v", err)
	}

	if ok := certPool.AppendCertsFromPEM(certData); !ok {
		return model.GenerateTextResponse{}, fmt.Errorf("не удалось добавить сертификат в пул")
	}

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			RootCAs: certPool,
		},
	}

	client := &http.Client{
		Transport: transport,
	}

	resp, err := client.Do(req)
	if err != nil {
		return model.GenerateTextResponse{}, fmt.Errorf("ошибка выполнения запроса: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return model.GenerateTextResponse{}, fmt.Errorf("ошибка - статус код %v", resp.StatusCode)
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return model.GenerateTextResponse{}, fmt.Errorf("ошибка чтения тела ответа: %v", err)
	}

	var response struct {
		Choices []struct {
			Message struct {
				Role    string `json:"role"`
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		return model.GenerateTextResponse{}, fmt.Errorf("ошибка парсинга ответа: %v", err)
	}

	if len(response.Choices) > 0 {
		content := response.Choices[0].Message.Content
		generateTextResponse := model.GenerateTextResponse{
			Message: content,
		}

		return generateTextResponse, nil
	}

	return model.GenerateTextResponse{}, fmt.Errorf("отсутствуют сгенерированные данные в ответе")
}
