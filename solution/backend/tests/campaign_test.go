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

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateAdvertiser(t *testing.T) {
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

func TestCreateCampaign_Success(t *testing.T) {
	cfg := config.LoadConfig()
	router := app.CreateRouter(cfg)

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

	*campaign.StartDate = 20250221
	*campaign.EndDate = 20250222
	*campaign.ImpressionsLimit = 1000
	*campaign.ClicksLimit = 500
	*campaign.CostPerImpression = 0.5
	*campaign.CostPerClick = 0.2

	campaignData, err := json.Marshal(campaign)
	if err != nil {
		t.Fatalf("Error marshalling data: %v", err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/advertisers/"+advertiserID+"/campaigns", bytes.NewReader(campaignData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, 201, w.Code)
}

func TestCreateCampaign_BadRequest_InvalidDateRange(t *testing.T) {
	cfg := config.LoadConfig()
	router := app.CreateRouter(cfg)

	advertiserID := uuid.New().String()
	campaign := model.CampaignCreate{
		AdTitle:           "Invalid Date Ad",
		AdText:            "This ad has invalid date range",
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

	*campaign.StartDate = 20250222
	*campaign.EndDate = 20250221 // End date is earlier than start date

	campaignData, err := json.Marshal(campaign)
	if err != nil {
		t.Fatalf("Error marshalling data: %v", err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/advertisers/"+advertiserID+"/campaigns", bytes.NewReader(campaignData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
}

func TestCreateCampaign_BadRequest_AgeRangeInvalid(t *testing.T) {
	cfg := config.LoadConfig()
	router := app.CreateRouter(cfg)

	advertiserID := uuid.New().String()
	campaign := model.CampaignCreate{
		AdTitle:           "Age Range Invalid Ad",
		AdText:            "This ad has an invalid age range",
		StartDate:         new(int32),
		EndDate:           new(int32),
		ImpressionsLimit:  new(int),
		ClicksLimit:       new(int),
		CostPerImpression: new(float64),
		CostPerClick:      new(float64),
		Targeting: model.CampaignTargeting{
			Gender:  "ALL",
			AgeFrom: new(int),
			AgeTo:   new(int),
		},
	}

	*campaign.StartDate = 20250221
	*campaign.EndDate = 20250222
	*campaign.ImpressionsLimit = 1000
	*campaign.ClicksLimit = 500
	*campaign.CostPerImpression = 0.5
	*campaign.CostPerClick = 0.2
	*campaign.Targeting.AgeFrom = 25
	*campaign.Targeting.AgeTo = 18 // Invalid age range

	campaignData, err := json.Marshal(campaign)
	if err != nil {
		t.Fatalf("Error marshalling data: %v", err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/advertisers/"+advertiserID+"/campaigns", bytes.NewReader(campaignData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
}

func TestCreateCampaign_BadRequest_AdvertiserNotFound(t *testing.T) {
	cfg := config.LoadConfig()
	router := app.CreateRouter(cfg)

	advertiserID := uuid.New().String() // Assuming this ID doesn't exist
	campaign := model.CampaignCreate{
		AdTitle:           "Ad for Non-Existing Advertiser",
		AdText:            "This ad is for an advertiser that doesn't exist",
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

	*campaign.StartDate = 20250221
	*campaign.EndDate = 20250222
	*campaign.ImpressionsLimit = 1000
	*campaign.ClicksLimit = 500
	*campaign.CostPerImpression = 0.5
	*campaign.CostPerClick = 0.2

	campaignData, err := json.Marshal(campaign)
	if err != nil {
		t.Fatalf("Error marshalling data: %v", err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/advertisers/"+advertiserID+"/campaigns", bytes.NewReader(campaignData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, 404, w.Code)
}

func TestGetMany_Success(t *testing.T) {
	cfg := config.LoadConfig()
	router := app.CreateRouter(cfg)

	advertiserID := "3fa85f64-5717-4562-b3fc-2c963f66afa6"
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/advertisers/"+advertiserID+"/campaigns", nil)
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var returnedCampaigns []model.Campaign
	err := json.Unmarshal(w.Body.Bytes(), &returnedCampaigns)
	if err != nil {
		t.Fatalf("Error unmarshalling response: %v", err)
	}
	assert.Equal(t, 2, len(returnedCampaigns))
}

func TestGetMany_BadRequest_InvalidPageParam(t *testing.T) {
	cfg := config.LoadConfig()
	router := app.CreateRouter(cfg)

	advertiserID := uuid.New().String()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/advertisers/"+advertiserID+"/campaigns?size=10&page=-1", nil)
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
}

func TestGetMany_BadRequest_InvalidSizeParam(t *testing.T) {
	cfg := config.LoadConfig()
	router := app.CreateRouter(cfg)

	advertiserID := uuid.New().String()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/advertisers/"+advertiserID+"/campaigns?size=abc&page=0", nil)
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
}

func TestGetMany_AdvertiserNotFound(t *testing.T) {
	cfg := config.LoadConfig()
	router := app.CreateRouter(cfg)

	advertiserID := uuid.New().String() // Assuming this ID doesn't exist

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/advertisers/"+advertiserID+"/campaigns?size=10&page=0", nil)
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// Ensure we get a 404 error if the advertiser is not found
	assert.Equal(t, 404, w.Code)
}

func TestGetMany_Success_EmptyCampaigns(t *testing.T) {
	cfg := config.LoadConfig()
	router := app.CreateRouter(cfg)

	advertiserID := "3fa85f64-5717-4562-b3fc-2c963f66afa6"

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/advertisers/"+advertiserID+"/campaigns?size=0&page=0", nil)
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var campaigns []model.Campaign
	err := json.Unmarshal(w.Body.Bytes(), &campaigns)
	if err != nil {
		t.Fatalf("Error unmarshalling response: %v", err)
	}
	assert.Equal(t, make([]model.Campaign, 0), campaigns)
}

func TestGetMany_Success_Pagination(t *testing.T) {
	cfg := config.LoadConfig()
	router := app.CreateRouter(cfg)

	advertiserID := "3fa85f64-5717-4562-b3fc-2c963f66afa6"
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/advertisers/"+advertiserID+"/campaigns?size=2&page=1", nil)
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var returnedCampaigns []model.Campaign
	err := json.Unmarshal(w.Body.Bytes(), &returnedCampaigns)
	if err != nil {
		t.Fatalf("Error unmarshalling response: %v", err)
	}

	assert.Equal(t, len(returnedCampaigns), 0)
}

func TestDeleteOne_Success(t *testing.T) {
	cfg := config.LoadConfig()
	router := app.CreateRouter(cfg)

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

	*campaign.StartDate = 20250221
	*campaign.EndDate = 20250222
	*campaign.ImpressionsLimit = 1000
	*campaign.ClicksLimit = 500
	*campaign.CostPerImpression = 0.5
	*campaign.CostPerClick = 0.2

	campaignData, err := json.Marshal(campaign)
	if err != nil {
		t.Fatalf("Error marshalling data: %v", err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/advertisers/"+advertiserID+"/campaigns", bytes.NewReader(campaignData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, 201, w.Code)

	var result model.Campaign
	err = json.Unmarshal(w.Body.Bytes(), &result)
	if err != nil {
		t.Fatalf("Error unmarshalling response: %v", err)
	}

	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("DELETE", "/advertisers/"+advertiserID+"/campaigns/"+result.CampaignID.String(), nil)
	req2.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w2, req2)

	// Ensure the campaign was deleted successfully (status code 204)
	assert.Equal(t, 204, w2.Code)
}

func TestDeleteOne_BadRequest_InvalidAdvertiserID(t *testing.T) {
	cfg := config.LoadConfig()
	router := app.CreateRouter(cfg)

	advertiserID := "invalid-uuid"
	campaignID := uuid.New().String()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/advertisers/"+advertiserID+"/campaigns/"+campaignID, nil)
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// Ensure we get a 400 error for invalid advertiser ID
	assert.Equal(t, 400, w.Code)
}

func TestDeleteOne_BadRequest_InvalidCampaignID(t *testing.T) {
	cfg := config.LoadConfig()
	router := app.CreateRouter(cfg)

	advertiserID := uuid.New().String()
	campaignID := "invalid-uuid"

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/advertisers/"+advertiserID+"/campaigns/"+campaignID, nil)
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// Ensure we get a 400 error for invalid campaign ID
	assert.Equal(t, 400, w.Code)
}

func TestDeleteOne_NotFound(t *testing.T) {
	cfg := config.LoadConfig()
	router := app.CreateRouter(cfg)

	advertiserID := uuid.New().String()
	campaignID := uuid.New().String()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/advertisers/"+advertiserID+"/campaigns/"+campaignID, nil)
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// Ensure we get a 404 error when the campaign is not found
	assert.Equal(t, 404, w.Code)
}
