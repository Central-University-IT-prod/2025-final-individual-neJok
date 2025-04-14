package handler

import (
	"github.com/gin-gonic/gin"
	"neJok/solution/internal/model"
	"neJok/solution/internal/service"
	"strconv"
)

type AIHandler struct {
	gigaChatSvc *service.GigaChatService
	actCacheSvc *service.ActCacheService
}

func NewAIHandler(gigaChatSvc *service.GigaChatService, actCacheSvc *service.ActCacheService) *AIHandler {
	return &AIHandler{gigaChatSvc, actCacheSvc}
}

// GenerateText godoc
// @Summary Сгенерировать текст для рекламной кампании
// @Description Создает текст для рекламной кампании по пожеланиям от пользователя
// @Tags AI
// @Accept json
// @Produce json
// @Param GenerateInfo body model.GenerateTextRequest true "Название товара и пожелания"
// @Success 200 {object} model.GenerateTextResponse "Обновленная кампания"
// @Failure 400 {object} model.ErrorResponse "Невалидные данные"
// @Failure 500 {object} model.GenerateTextResponse "Ошибка сервера"
// @Router /ai/text/generate [post]
func (h *AIHandler) GenerateText(c *gin.Context) {
	req := model.GenerateTextRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, BuildError(err.Error()))
		return
	}

	response, code := h.gigaChatSvc.GenerateText(req)
	c.JSON(code, response)
}

// SetModeration godoc
// @Summary Включить или выключить модерацию текстов
// @Description Включает или выключает модерацию текстов и названий рекламных кампаний
// @Tags AI
// @Accept json
// @Produce json
// @Param Moderation body model.ModerationRequest true "true - включить модерацию, false - выключить модерацию"
// @Success 200 {object} model.ModerationRequest "Обновленная кампания"
// @Failure 400 {object} model.ErrorResponse "Невалидные данные"
// @Failure 500 {object} model.ErrorResponse "Ошибка сервера"
// @Router /ai/text/moderation [post]
func (h *AIHandler) SetModeration(c *gin.Context) {
	req := model.ModerationRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, BuildError(err.Error()))
		return
	}

	h.actCacheSvc.SetStr("moderation", strconv.FormatBool(*req.Status), nil)
	c.JSON(200, req)
}
