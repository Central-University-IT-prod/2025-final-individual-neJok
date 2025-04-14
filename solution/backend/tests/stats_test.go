package tests

import (
	"bytes"
	"encoding/json"
	"neJok/solution/app"
	"neJok/solution/config"
	"neJok/solution/internal/model"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStats_Success(t *testing.T) {
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

	advertiserID := "3fa85f64-5717-4562-b3fc-2c963f66afa6"
	campaign := model.CampaignCreate{
		AdTitle:           "Sample Ad",
		AdText:            "This is a sample ad text",
		StartDate:         new(int32),
		EndDate:           new(int32),
		ImpressionsLimit:  new(int),
		ClicksLimit:       new(int),
		CostPerImpression: new(float64),
		CostPerClick:      new(float64),
		Targeting: model.CampaignTargeting{
			Gender: "ALL",
		},
	}

	*campaign.StartDate = 0
	*campaign.EndDate = 20250222
	*campaign.ImpressionsLimit = 1
	*campaign.ClicksLimit = 1
	*campaign.CostPerImpression = 0.5
	*campaign.CostPerClick = 0.2

	campaignData, err := json.Marshal(campaign)
	if err != nil {
		t.Fatalf("Error marshalling data: %v", err)
	}

	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("POST", "/advertisers/"+advertiserID+"/campaigns", bytes.NewReader(campaignData))
	req2.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w2, req2)

	assert.Equal(t, 201, w2.Code)

	clientID := "3fa85f64-5717-4562-b3fc-2c963f66afa6"
	client := []map[string]interface{}{
		{
			"client_id": clientID,
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

	w3 := httptest.NewRecorder()
	req3, _ := http.NewRequest("POST", "/clients/bulk", bytes.NewReader(clientData))
	req3.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w3, req3)
	assert.Equal(t, 201, w3.Code)

	w4 := httptest.NewRecorder()
	req4, _ := http.NewRequest("GET", "/ads?client_id="+clientID, bytes.NewReader(campaignData))
	req4.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w4, req4)

	assert.Equal(t, 200, w4.Code)

	var response model.AdsResponse
	err = json.Unmarshal(w4.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Error unmarshalling response: %v", err)
	}

	adsRequest := model.AdsRequest{
		ClientID: clientID,
	}
	w5 := httptest.NewRecorder()

	adsRequestData, err := json.Marshal(adsRequest)
	if err != nil {
		t.Fatalf("Error marshalling data: %v", err)
	}

	req5, _ := http.NewRequest("POST", "/ads/"+response.CampaignID.String()+"/click", bytes.NewReader(adsRequestData))
	req5.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w5, req5)

	assert.Equal(t, 204, w5.Code)

	w6 := httptest.NewRecorder()
	req6, _ := http.NewRequest("GET", "/stats/campaigns/"+response.CampaignID.String(), nil)
	req6.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w6, req6)

	assert.Equal(t, 200, w6.Code)

	var stats model.CampaignStats
	err = json.Unmarshal(w6.Body.Bytes(), &stats)
	if err != nil {
		t.Fatalf("Error unmarshalling response: %v", err)
	}

	assert.Equal(t, 1, int(stats.TotalClicks))
	assert.Equal(t, 1, int(stats.TotalViews))

	w7 := httptest.NewRecorder()
	req7, _ := http.NewRequest("GET", "/stats/campaigns/"+response.CampaignID.String()+"/daily", nil)
	req7.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w7, req7)

	assert.Equal(t, 200, w7.Code)

	var dailyStats []model.CampaignStatsDaily
	err = json.Unmarshal(w7.Body.Bytes(), &dailyStats)
	if err != nil {
		t.Fatalf("Error unmarshalling response: %v", err)
	}

	assert.Equal(t, 1, int(dailyStats[0].TotalClicks))
	assert.Equal(t, 0, int(dailyStats[0].Date))
}
