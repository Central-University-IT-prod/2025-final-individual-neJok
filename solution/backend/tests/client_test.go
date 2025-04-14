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

func TestClientCreateOrUpdate_Success(t *testing.T) {
	cfg := config.LoadConfig()
	router := app.CreateRouter(cfg)

	client := []map[string]interface{}{
		{
			"client_id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
			"login":     "string",
			"age":       25,
			"location":  "New York",
			"gender":    "MALE",
		},
	}

	clientData, err := json.Marshal(client)
	if err != nil {
		t.Fatalf("Error marshalling data: %v", err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/clients/bulk", bytes.NewReader(clientData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, 201, w.Code)
}

func TestClientCreateOrUpdate_BadRequest_EmptyBody(t *testing.T) {
	cfg := config.LoadConfig()
	router := app.CreateRouter(cfg)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/clients/bulk", nil)
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
}

func TestClientCreateOrUpdate_BadRequest_MissingField(t *testing.T) {
	cfg := config.LoadConfig()
	router := app.CreateRouter(cfg)

	client := []map[string]interface{}{
		{
			"client_id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
			// "login" field is missing
			"age":      25,
			"location": "New York",
			"gender":   "MALE",
		},
	}

	clientData, err := json.Marshal(client)
	if err != nil {
		t.Fatalf("Error marshalling data: %v", err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/clients/bulk", bytes.NewReader(clientData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
}

func TestClientCreateOrUpdate_BadRequest_InvalidAge(t *testing.T) {
	cfg := config.LoadConfig()
	router := app.CreateRouter(cfg)

	client := []map[string]interface{}{
		{
			"client_id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
			"login":     "string",
			"age":       -1,
			"location":  "New York",
			"gender":    "MALE",
		},
	}

	clientData, err := json.Marshal(client)
	if err != nil {
		t.Fatalf("Error marshalling data: %v", err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/clients/bulk", bytes.NewReader(clientData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
}

func TestClientCreateOrUpdate_BadRequest_InvalidGender(t *testing.T) {
	cfg := config.LoadConfig()
	router := app.CreateRouter(cfg)

	client := []map[string]interface{}{
		{
			"client_id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
			"login":     "string",
			"age":       25,
			"location":  "New York",
			"gender":    "INVALID_GENDER", // Invalid gender
		},
	}

	clientData, err := json.Marshal(client)
	if err != nil {
		t.Fatalf("Error marshalling data: %v", err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/clients/bulk", bytes.NewReader(clientData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
}

func TestClientCreateOrUpdate_Success_EmptyClientList(t *testing.T) {
	cfg := config.LoadConfig()
	router := app.CreateRouter(cfg)

	client := []map[string]interface{}{} // Empty list

	clientData, err := json.Marshal(client)
	if err != nil {
		t.Fatalf("Error marshalling data: %v", err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/clients/bulk", bytes.NewReader(clientData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, 201, w.Code)
}

func TestClientCreateOrUpdate_BadRequest_InvalidJSON(t *testing.T) {
	cfg := config.LoadConfig()
	router := app.CreateRouter(cfg)

	invalidJSON := `{"client_id": "3fa85f64-5717-4562-b3fc-2c963f66afa6", "login": "string", "age": 25, "location": "New York", "gender": "MALE"`

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/clients/bulk", bytes.NewReader([]byte(invalidJSON)))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
}

func TestClientGetById_Success(t *testing.T) {
	cfg := config.LoadConfig()
	router := app.CreateRouter(cfg)

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/clients/3fa85f64-5717-4562-b3fc-2c963f66afa6", nil)
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}

func TestClientGetById_NotFound(t *testing.T) {
	cfg := config.LoadConfig()
	router := app.CreateRouter(cfg)

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/clients/3fa85f64-5717-4562-b3fc-2c963f66afa9", nil)
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	assert.Equal(t, 404, w.Code)
}

func TestClientGetById_InvalidId(t *testing.T) {
	cfg := config.LoadConfig()
	router := app.CreateRouter(cfg)

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/clients/3fa85f64-5717-4562-b3fc-2c963f66a", nil)
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)
}
