package handler

import (
	"github.com/gin-gonic/gin"
	"neJok/solution/internal/model"
	"neJok/solution/internal/service"
)

type TimeHandler struct {
	actCacheSvc *service.ActCacheService
}

func NewTimeHandler(actCacheSvc *service.ActCacheService) *TimeHandler {
	return &TimeHandler{actCacheSvc}
}

// Set godoc
// @Summary Установить текущую дату
// @Description Устанавливает текущую дату в кэш, только если она не меньше текущей даты
// @Tags Время
// @Accept json
// @Produce json
// @Param time body model.TimeSetRequest true "Запрос на установку текущей даты"
// @Success 200 {object} model.TimeSetRequest "Дата успешно установлена"
// @Failure 400 {object} model.ErrorResponse "Неверная дата"
// @Failure 500 {object} model.ErrorResponse "Ошибка сервера"
// @Router /time/advance [post]
func (h *TimeHandler) Set(c *gin.Context) {
	req := model.TimeSetRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, BuildError(err.Error()))
		return
	}

	currentDate, err := h.actCacheSvc.GetInt("current_day")
	if err != nil {
		HandleError500(c, err)
		return
	}

	int32CurrentDate := int32(currentDate)
	if *req.CurrentDate < int32CurrentDate {
		c.JSON(400, BuildError("Request current date cannot be less than now date"))
		return
	}

	err = h.actCacheSvc.SetInt("current_day", int(*req.CurrentDate))
	if err != nil {
		HandleError500(c, err)
		return
	}

	if *req.CurrentDate != int32CurrentDate {
		go func() {
			h.actCacheSvc.DeleteKeysByPrefix("top_ads:")
		}()
	}

	c.JSON(200, req)
}

// Get godoc
// @Summary Получить текущую дату
// @Description Возвращает текущую дату
// @Tags Время
// @Accept json
// @Produce json
// @Success 200 {object} model.TimeSetRequest "Дата успешно получена"
// @Failure 500 {object} model.ErrorResponse "Ошибка сервера"
// @Router /time/advance [get]
func (h *TimeHandler) Get(c *gin.Context) {
	currentDate, err := h.actCacheSvc.GetInt("current_day")
	if err != nil {
		HandleError500(c, err)
		return
	}

	currentDatePtr := int32(currentDate)
	c.JSON(200, model.TimeSetRequest{CurrentDate: &currentDatePtr})
}