package tests

import (
	"github.com/stretchr/testify/assert"
	"neJok/solution/config"
	"net/http"
	"net/http/httptest"
	"testing"

	"neJok/solution/app"
)

func TestPing_Success(t *testing.T) {
	cfg := config.LoadConfig()
	router := app.CreateRouter(cfg)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}
