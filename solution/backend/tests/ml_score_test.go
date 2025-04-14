package tests

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"neJok/solution/config"
	"net/http"
	"net/http/httptest"
	"testing"

	"neJok/solution/app"
)

func TestMLScoreCreateOrUpdate_Success(t *testing.T) {
	cfg := config.LoadConfig()
	router := app.CreateRouter(cfg)

	advertiser := map[string]interface{}{
		"client_id":     "3fa85f64-5717-4562-b3fc-2c963f66afa6",
		"advertiser_id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
		"score":         0,
	}

	advertiserData, err := json.Marshal(advertiser)
	if err != nil {
		t.Fatalf("Error marshalling data: %v", err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/ml-scores", bytes.NewReader(advertiserData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

func TestMLScoreCreateOrUpdate_BadRequest_EmptyBody(t *testing.T) {
	cfg := config.LoadConfig()
	router := app.CreateRouter(cfg)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/ml-scores", nil)
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
}

func TestMLScoreCreateOrUpdate_BadRequest_MissingField(t *testing.T) {
	cfg := config.LoadConfig()
	router := app.CreateRouter(cfg)

	advertiser := map[string]interface{}{
		"client_id":     "3fa85f64-5717-4562-b3fc-2c963f66afa6",
		"advertiser_id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
		// "score" field is missing
	}

	advertiserData, err := json.Marshal(advertiser)
	if err != nil {
		t.Fatalf("Error marshalling data: %v", err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/ml-scores", bytes.NewReader(advertiserData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
}

func TestMLScoreCreateOrUpdate_BadRequest_InvalidScore(t *testing.T) {
	cfg := config.LoadConfig()
	router := app.CreateRouter(cfg)

	advertiser := map[string]interface{}{
		"client_id":     "3fa85f64-5717-4562-b3fc-2c963f66afa6",
		"advertiser_id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
		"score":         "4",
	}

	advertiserData, err := json.Marshal(advertiser)
	if err != nil {
		t.Fatalf("Error marshalling data: %v", err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/ml-scores", bytes.NewReader(advertiserData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
}
