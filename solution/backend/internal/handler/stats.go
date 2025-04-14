package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"neJok/solution/internal/model"
	"neJok/solution/internal/service"
)

type StatsHandler struct {
	advertiserSvc *service.AdvertiserService
	campaignSvc   *service.CampaignService
	adsHistorySvc *service.AdsHistoryService
	actCacheSvc   *service.ActCacheService
}

func NewStatsHandler(advertiserSvc *service.AdvertiserService, campaignSvc *service.CampaignService, adsHistory *service.AdsHistoryService, actCacheSvc *service.ActCacheService) *StatsHandler {
	return &StatsHandler{advertiserSvc, campaignSvc, adsHistory, actCacheSvc}
}

// GetCampaignStats godoc
// @Summary Получить статистику по кампании
// @Description Возвращает агрегированную статистику по указанной кампании
// @Tags Статистика
// @Accept json
// @Produce json
// @Param campaignID path string true "Идентификатор кампании" Format(uuid)
// @Success 200 {object} model.CampaignStats "Агрегированная статистика по кампании"
// @Failure 400 {object} model.ErrorResponse "Неверный идентификатор кампании"
// @Failure 404 {object} model.ErrorResponse "Кампания не найдена"
// @Failure 500 {object} model.ErrorResponse "Ошибка сервера"
// @Router /stats/campaigns/{campaignID} [get]
func (h *StatsHandler) GetCampaignStats(c *gin.Context) {
	campaignID, err := uuid.Parse(c.Param("campaignID"))
	if err != nil {
		c.JSON(400, BuildError(err.Error()))
		return
	}

	_, err = h.campaignSvc.GetByID(campaignID)
	if err != nil {
		c.JSON(404, BuildError("Campaign not found"))
		return
	}

	stats, err := h.adsHistorySvc.GetAggregatedCampaignStats(campaignID)
	if err != nil {
		HandleError500(c, err)
		return
	}
	c.JSON(200, stats)
}

// GetAdvertiserStats godoc
// @Summary Получить статистику по рекламодателю
// @Description Возвращает агрегированную статистику по указанному рекламодателю
// @Tags Статистика
// @Accept json
// @Produce json
// @Param advertiserID path string true "Идентификатор рекламодателя" Format(uuid)
// @Success 200 {object} model.CampaignStats "Агрегированная статистика по рекламодателю"
// @Failure 400 {object} model.ErrorResponse "Неверный идентификатор рекламодателя"
// @Failure 404 {object} model.ErrorResponse "Рекламодатель не найден"
// @Failure 500 {object} model.ErrorResponse "Ошибка сервера"
// @Router /stats/advertisers/{advertiserID}/campaigns [get]
func (h *StatsHandler) GetAdvertiserStats(c *gin.Context) {
	advertiserID, err := uuid.Parse(c.Param("advertiserID"))
	if err != nil {
		c.JSON(400, BuildError(err.Error()))
		return
	}

	_, err = h.advertiserSvc.GetByID(advertiserID)
	if err != nil {
		c.JSON(404, BuildError("Advertiser not found"))
		return
	}

	stats, err := h.adsHistorySvc.GetAggregatedAdvertiserStats(advertiserID)
	if err != nil {
		HandleError500(c, err)
		return
	}
	c.JSON(200, stats)
}

// GetCampaignDailyStats godoc
// @Summary Получить дневную статистику по кампании
// @Description Возвращает агрегированную статистику по кампаниям за каждый день от даты начала до текущего дня
// @Tags Статистика
// @Accept json
// @Produce json
// @Param campaignID path string true "Идентификатор кампании" Format(uuid)
// @Success 200 {array} model.CampaignStatsDaily "Агрегированная статистика по кампании за каждый день"
// @Failure 400 {object} model.ErrorResponse "Неверный запрос"
// @Failure 404 {object} model.ErrorResponse "Кампания не найдена"
// @Failure 500 {object} model.ErrorResponse "Ошибка сервера"
// @Router /stats/campaigns/{campaignID}/daily [get]
func (h *StatsHandler) GetCampaignDailyStats(c *gin.Context) {
	campaignID, err := uuid.Parse(c.Param("campaignID"))
	if err != nil {
		c.JSON(400, BuildError(err.Error()))
		return
	}

	campaign, err := h.campaignSvc.GetByID(campaignID)
	if err != nil {
		c.JSON(404, BuildError("Campaign not found"))
		return
	}

	currentDate, err := h.actCacheSvc.GetInt("current_day")
	if err != nil {
		HandleError500(c, err)
		return
	}

	if currentDate < int(*campaign.StartDate) {
		c.JSON(200, make([]model.CampaignStatsDaily, 0))
		return
	}

	endDate := *campaign.EndDate
	if currentDate < int(endDate) {
		endDate = int32(currentDate)
	}

	stats, err := h.adsHistorySvc.GetAggregatedCampaignDailyStats(campaignID, *campaign.StartDate, endDate)
	if err != nil {
		HandleError500(c, err)
		return
	}
	c.JSON(200, stats)
}

// GetAdvertiserDailyStats godoc
// @Summary Получить дневную статистику по рекламодателю
// @Description Возвращает агрегированную статистику по рекламодателю за каждый день, начиная с даты создания рекламодателя
// @Tags Статистика
// @Accept json
// @Produce json
// @Param advertiserID path string true "Идентификатор рекламодателя" Format(uuid)
// @Success 200 {array} model.CampaignStatsDaily "Агрегированная статистика по рекламодателю за каждый день"
// @Failure 400 {object} model.ErrorResponse "Неверный идентификатор рекламодателя"
// @Failure 404 {object} model.ErrorResponse "Рекламодатель не найден"
// @Failure 500 {object} model.ErrorResponse "Ошибка сервера"
// @Router /stats/advertisers/{advertiserID}/campaigns/daily [get]
func (h *StatsHandler) GetAdvertiserDailyStats(c *gin.Context) {
	advertiserID, err := uuid.Parse(c.Param("advertiserID"))
	if err != nil {
		c.JSON(400, BuildError(err.Error()))
		return
	}

	advertiser, err := h.advertiserSvc.GetByID(advertiserID)
	if err != nil {
		c.JSON(404, BuildError("Advertiser not found"))
		return
	}

	currentDate, err := h.actCacheSvc.GetInt("current_day")
	if err != nil {
		HandleError500(c, err)
		return
	}

	stats, err := h.adsHistorySvc.GetAggregatedAdvertiserDailyStats(advertiserID, advertiser.CreatedAt, int32(currentDate))
	if err != nil {
		HandleError500(c, err)
		return
	}
	c.JSON(200, stats)
}
