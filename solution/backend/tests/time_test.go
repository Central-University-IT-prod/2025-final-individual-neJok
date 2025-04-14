package tests

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"neJok/solution/config"
	"neJok/solution/internal/model"
	"net/http"
	"net/http/httptest"
	"testing"

	"neJok/solution/app"
)

func TestSetCurrentDate_Success(t *testing.T) {
	cfg := config.LoadConfig()
	router := app.CreateRouter(cfg)

	body := map[string]interface{}{
		"current_date": 0,
	}

	data, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("Error marshalling data: %v", err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/time/advance", bytes.NewReader(data))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

func TestGetCurrentDate_Success(t *testing.T) {
	cfg := config.LoadConfig()
	router := app.CreateRouter(cfg)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/time/advance", nil)
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var response model.TimeSetRequest
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Error unmarshalling response: %v", err)
	}
	assert.Equal(t, 0, int(*response.CurrentDate))
}

func TestSetCurrentDate_Invalid(t *testing.T) {
	cfg := config.LoadConfig()
	router := app.CreateRouter(cfg)

	body := map[string]interface{}{
		"current_date": -1,
	}

	data, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("Error marshalling data: %v", err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/time/advance", bytes.NewReader(data))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
}
