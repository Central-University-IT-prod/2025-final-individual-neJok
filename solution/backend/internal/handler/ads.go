package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"neJok/solution/internal/model"
	"neJok/solution/internal/service"
	campaignUtil "neJok/solution/pkg/campaign"
	"sort"
	"sync"
)

type AdsHandler struct {
	clientSvc     *service.ClientService
	actCacheSvc   *service.ActCacheService
	campaignSvc   *service.CampaignService
	adsHistorySvc *service.AdsHistoryService
	mlScoreSvc    *service.MLScoreService
}

func NewAdsHandler(clientSvc *service.ClientService, actCacheSvc *service.ActCacheService, campaignSvc *service.CampaignService, adsHistorySvc *service.AdsHistoryService, mlScoreSvc *service.MLScoreService) *AdsHandler {
	return &AdsHandler{clientSvc, actCacheSvc, campaignSvc, adsHistorySvc, mlScoreSvc}
}

func (h *AdsHandler) CalculateCampaigns(campaigns []model.CampaignForUser, stats model.CampaignsDBStats, clientID uuid.UUID) []model.CampaignForUser {
	var wg sync.WaitGroup
	var bestCampaigns []model.CampaignForUser
	for _, campaign := range campaigns {
		wg.Add(1)
		go func(campaign model.CampaignForUser) {
			defer wg.Done()
			dbScore, _ := h.actCacheSvc.GetInt(fmt.Sprintf("%s:%s", campaign.AdvertiserID, clientID))
			campaign.Score = float64(dbScore)
			score := campaignUtil.CalculateCampaignScore(campaign, stats.MaxScore, stats.MaxEndDate, stats.MinEndDate, stats.MaxCostPerImpression, stats.MaxCostPerClick)
			campaign.Score = score
			bestCampaigns = append(bestCampaigns, campaign)
		}(campaign)
	}
	wg.Wait()

	sort.Slice(bestCampaigns, func(i, j int) bool {
		return bestCampaigns[i].Score > bestCampaigns[j].Score
	})
	return bestCampaigns
}

func (h *AdsHandler) CheckCampaigns(campaigns []model.CampaignForUser, redisKey string, clientID uuid.UUID, currentDate int) (model.AdsResponse, error) {
	var maxScoreCampaign model.CampaignForUser

	for i, campaign := range campaigns {
		views, clicks, err := h.adsHistorySvc.GetViewsAndClicks(campaign.CampaignID)
		if err != nil {
			continue
		}

		if float64(views)+1 >= float64(campaign.ImpressionsLimit)*1.05 && campaign.ViewsCount == 0 {
			continue
		}

		if campaign.ViewsCount == 0 && int(clicks) < campaign.ClicksLimit {
			if i+1 >= len(campaigns) {
				h.actCacheSvc.SetList(redisKey, make([]model.CampaignForUser, 0))
			} else {
				h.actCacheSvc.SetList(redisKey, campaigns[i+1:])
			}

			h.adsHistorySvc.Add(campaign.CampaignID, campaign.AdvertiserID, clientID, campaign.CostPerImpression, "view", currentDate)
			return model.AdsResponse{
				CampaignID:   campaign.CampaignID,
				AdvertiserID: campaign.AdvertiserID,
				AdTitle:      campaign.AdTitle,
				AdText:       campaign.AdText,
				ImageURL:     campaign.ImageURL,
			}, nil
		} else if maxScoreCampaign.CampaignID == uuid.Nil || campaign.Score > maxScoreCampaign.Score {
			maxScoreCampaign = campaign
		}
	}

	h.actCacheSvc.SetList(redisKey, make([]model.CampaignForUser, 0))
	if maxScoreCampaign.CampaignID != uuid.Nil {
		return model.AdsResponse{
			CampaignID:   maxScoreCampaign.CampaignID,
			AdvertiserID: maxScoreCampaign.AdvertiserID,
			AdTitle:      maxScoreCampaign.AdTitle,
			AdText:       maxScoreCampaign.AdText,
			ImageURL:     maxScoreCampaign.ImageURL,
		}, nil
	}

	return model.AdsResponse{}, nil
}

// GetOne godoc
// @Summary Получить информацию о рекламе для клиента
// @Description Возвращает информацию о рекламе, соответствующей таргетингу клиента.
// @Tags Реклама
// @Accept json
// @Produce json
// @Param client_id query string true "Идентификатор клиента" Format(uuid)
// @Success 200 {object} model.AdsResponse "Информация о кампании"
// @Failure 400 {object} model.ErrorResponse "Некорректный идентификатор клиента"
// @Failure 404 {object} model.ErrorResponse "Клиент не найден или не найдена подходящая кампания"
// @Failure 500 {object} model.ErrorResponse "Ошибка сервера"
// @Router /ads [get]
func (h *AdsHandler) GetOne(c *gin.Context) {
	var req model.AdsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(400, BuildError(err.Error()))
		return
	}

	clientID, err := uuid.Parse(req.ClientID)
	if err != nil {
		c.JSON(400, BuildError("Invalid client ID"))
		return
	}

	client, err := h.clientSvc.GetByID(clientID)
	if err != nil {
		c.JSON(404, BuildError("Client not found"))
		return
	}

	currentDate, err := h.actCacheSvc.GetInt("current_day")
	if err != nil {
		HandleError500(c, err)
		return
	}

	redisKey := fmt.Sprintf("top_ads:%s", clientID.String())
	topAds, err := h.actCacheSvc.GetList(redisKey)
	if err == nil && len(topAds) > 0 {
		campaign, err := h.CheckCampaigns(topAds, redisKey, clientID, currentDate)
		if err != nil {
			HandleError500(c, err)
			return
		}
		if campaign.CampaignID != uuid.Nil {
			c.JSON(200, campaign)
			return
		}
	}
	campaigns, stats, err := h.campaignSvc.GetManyByTargeting(client.Gender, client.Age, client.Login, currentDate, clientID)
	if err != nil {
		HandleError500(c, err)
		return
	}

	if len(campaigns) == 0 {
		c.JSON(404, BuildError("No campaigns found"))
		return
	}

	stats.MaxScore, _ = h.mlScoreSvc.GetMax(clientID)

	campaigns = h.CalculateCampaigns(campaigns, stats, clientID)
	campaign, err := h.CheckCampaigns(campaigns, redisKey, clientID, currentDate)
	if err != nil {
		HandleError500(c, err)
		return
	}
	if campaign.CampaignID == uuid.Nil {
		c.JSON(404, BuildError("No campaigns found"))
		return
	}
	c.JSON(200, campaign)
}

// Click godoc
// @Summary Зарегистрировать клик по рекламе
// @Description Регистрирует клик клиента по рекламе, проверяя наличие предыдущих действий (просмотр и клик) и сроки действия кампании
// @Tags Реклама
// @Accept json
// @Produce json
// @Param campaignID path string true "Идентификатор кампании" Format(uuid)
// @Param ClientID body model.AdsRequest true "Идентификатор клиента"
// @Success 204 "Клик успешно зарегистрирован"
// @Failure 400 {object} model.ErrorResponse "Неверный формат идентификаторов клиента или кампании, или клиент уже кликнул"
// @Failure 404 {object} model.ErrorResponse "Клиент или кампания не найдены"
// @Failure 410 {object} model.ErrorResponse "Кампания не началась или завершена"
// @Failure 500 {object} model.ErrorResponse "Ошибка сервера"
// @Router /ads/{campaignID}/click [post]
func (h *AdsHandler) Click(c *gin.Context) {
	var req model.AdsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, BuildError(err.Error()))
		return
	}

	clientID, err := uuid.Parse(req.ClientID)
	if err != nil {
		c.JSON(400, BuildError("Invalid client ID"))
		return
	}

	campaignID, err := uuid.Parse(c.Param("campaignID"))
	if err != nil {
		c.JSON(400, BuildError(err.Error()))
		c.Abort()
		return
	}

	_, err = h.clientSvc.GetByID(clientID)
	if err != nil {
		c.JSON(404, BuildError("Client not found"))
		return
	}

	campaign, err := h.campaignSvc.GetByID(campaignID)
	if err != nil {
		c.JSON(404, BuildError("Campaign not found"))
		return
	}

	_, viewCreated := h.adsHistorySvc.GetOne(campaignID, clientID, "view")
	if viewCreated != nil {
		c.JSON(400, BuildError("Client view not found"))
		return
	}

	_, clickCreated := h.adsHistorySvc.GetOne(campaignID, clientID, "click")
	if clickCreated == nil {
		c.Status(204)
		return
	}

	currentDate, err := h.actCacheSvc.GetInt("current_day")
	if err != nil {
		HandleError500(c, err)
		return
	}
	if currentDate > int(*campaign.EndDate) || currentDate < int(*campaign.StartDate) {
		c.JSON(410, BuildError("The advertising campaign has not started yet"))
		return
	}

	err = h.adsHistorySvc.Add(campaignID, campaign.AdvertiserID, clientID, *campaign.CostPerClick, "click", currentDate)
	if err != nil {
		HandleError500(c, err)
		return
	}
	c.Status(204)
}
