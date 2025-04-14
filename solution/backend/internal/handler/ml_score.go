package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"neJok/solution/internal/model"
	"neJok/solution/internal/service"
)

type MLScoreHandler struct {
	mlScoreSvc    *service.MLScoreService
	clientSvc     *service.ClientService
	advertiserSvc *service.AdvertiserService
	actCacheSvc   *service.ActCacheService
}

func NewMLScoreHandler(mlScoreSvc *service.MLScoreService, clientSvc *service.ClientService, advertiserSvc *service.AdvertiserService, actCacheSvc *service.ActCacheService) *MLScoreHandler {
	return &MLScoreHandler{mlScoreSvc, clientSvc, advertiserSvc, actCacheSvc}
}

// CreateOrUpdate godoc
// @Summary Создать или обновить ML-оценку
// @Description Создает или обновляет ML-оценку для клиента и рекламодателя. Я бы хотел возвращать тут 204, но по заданию в спецификации стоит 200.
// @Tags ML-оценки
// @Accept json
// @Produce json
// @Param mlScore body model.MLScore true "Модель ML-оценки"
// @Success 200 "Операция успешно выполнена"
// @Failure 400 {object} model.ErrorResponse "Неверные данные или отсутствуют клиент или рекламодатель"
// @Failure 404 {object} model.ErrorResponse "Клиент или рекламодатель не найдены"
// @Failure 500 {object} model.ErrorResponse "Ошибка сервера"
// @Router /ml-scores [post]
func (h *MLScoreHandler) CreateOrUpdate(c *gin.Context) {
	req := model.MLScore{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, BuildError(err.Error()))
		return
	}

	_, err := h.clientSvc.GetByID(req.ClientID)
	if err != nil {
		c.JSON(404, BuildError("Client not found"))
		return
	}

	_, err = h.advertiserSvc.GetByID(req.AdvertiserID)
	if err != nil {
		c.JSON(404, BuildError("Advertiser not found"))
		return
	}

	err = h.mlScoreSvc.CreateOrUpdate(req)
	if err != nil {
		HandleError500(c, err)
		return
	}

	err = h.actCacheSvc.SetInt(fmt.Sprintf("%s:%s", req.AdvertiserID.String(), req.ClientID.String()), *req.Score)
	if err != nil {
		HandleError500(c, err)
		return
	}

	redisKey := fmt.Sprintf("top_ads:%s", req.ClientID.String())
	err = h.actCacheSvc.SetList(redisKey, make([]model.CampaignForUser, 0))
	if err != nil {
		HandleError500(c, err)
		return
	}

	c.Status(200)
}
