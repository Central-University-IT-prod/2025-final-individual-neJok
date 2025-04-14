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

func TestAdvertiserCreateOrUpdate_Success(t *testing.T) {
	cfg := config.LoadConfig()
	router := app.CreateRouter(cfg)

	advertiser := []map[string]interface{}{
		{
			"advertiser_id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
			"name":          "string",
		},
	}

	advertiserData, err := json.Marshal(advertiser)
	if err != nil {
		t.Fatalf("Error marshalling data: %v", err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/advertisers/bulk", bytes.NewReader(advertiserData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, 201, w.Code)
}

func TestAdvertiserCreateOrUpdate_BadRequest_EmptyBody(t *testing.T) {
	cfg := config.LoadConfig()
	router := app.CreateRouter(cfg)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/advertisers/bulk", nil)
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
}

func TestAdvertiserCreateOrUpdate_BadRequest_MissingField(t *testing.T) {
	cfg := config.LoadConfig()
	router := app.CreateRouter(cfg)

	advertiser := []map[string]interface{}{
		{
			"advertiser_id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
			// "name" field is missing
		},
	}

	advertiserData, err := json.Marshal(advertiser)
	if err != nil {
		t.Fatalf("Error marshalling data: %v", err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/advertisers/bulk", bytes.NewReader(advertiserData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
}

func TestAdvertiserCreateOrUpdate_BadRequest_InvalidName(t *testing.T) {
	cfg := config.LoadConfig()
	router := app.CreateRouter(cfg)

	advertiser := []map[string]interface{}{
		{
			"advertiser_id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
			"name":          "",
		},
	}

	advertiserData, err := json.Marshal(advertiser)
	if err != nil {
		t.Fatalf("Error marshalling data: %v", err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/advertisers/bulk", bytes.NewReader(advertiserData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
}

func TestAdvertiserCreateOrUpdate_Success_EmptyAdvertiserList(t *testing.T) {
	cfg := config.LoadConfig()
	router := app.CreateRouter(cfg)

	advertiser := []map[string]interface{}{} // Empty list

	advertiserData, err := json.Marshal(advertiser)
	if err != nil {
		t.Fatalf("Error marshalling data: %v", err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/advertisers/bulk", bytes.NewReader(advertiserData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, 201, w.Code)
}

func TestAdvertiserCreateOrUpdate_BadRequest_InvalidJSON(t *testing.T) {
	cfg := config.LoadConfig()
	router := app.CreateRouter(cfg)

	invalidJSON := `{"advertiser_id": "3fa85f64-5717-4562-b3fc-2c963f66afa6", "name": "string"`

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/advertisers/bulk", bytes.NewReader([]byte(invalidJSON)))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
}

func TestAdvertiserGetById_Success(t *testing.T) {
	cfg := config.LoadConfig()
	router := app.CreateRouter(cfg)

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/advertisers/3fa85f64-5717-4562-b3fc-2c963f66afa6", nil)
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}

func TestAdvertiserGetById_NotFound(t *testing.T) {
	cfg := config.LoadConfig()
	router := app.CreateRouter(cfg)

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/advertisers/3fa85f64-5717-4562-b3fc-2c963f66afa9", nil)
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	assert.Equal(t, 404, w.Code)
}

func TestAdvertiserGetById_InvalidId(t *testing.T) {
	cfg := config.LoadConfig()
	router := app.CreateRouter(cfg)

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/advertisers/3fa85f64-5717-4562-b3fc-2c963f66a", nil)
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)
}
