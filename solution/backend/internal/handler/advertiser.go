package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"neJok/solution/internal/model"
	"neJok/solution/internal/service"
)

type AdvertiserHandler struct {
	advertiserSvc *service.AdvertiserService
	actCacheSvc   *service.ActCacheService
}

func NewAdvertiserHandler(advertiserSvc *service.AdvertiserService, actCacheSvc *service.ActCacheService) *AdvertiserHandler {
	return &AdvertiserHandler{advertiserSvc, actCacheSvc}
}

// CreateOrUpdate godoc
// @Summary Создать или обновить рекламодателя
// @Description Создает или обновляет информацию о рекламодателе с проверкой уникальности идентификатора
// @Tags Рекламодатели
// @Accept json
// @Produce json
// @Param advertisers body model.AdvertiserCreateOrUpdateRequest true "Список рекламодателей для создания или обновления"
// @Success 201 {object} model.AdvertiserCreateOrUpdateRequest "Информация о рекламодателях, которые были успешно созданы или обновлены"
// @Failure 400 {object} model.ErrorResponse "Невалидные данные или повторяющиеся идентификаторы рекламодателей"
// @Failure 500 {object} model.ErrorResponse "Ошибка сервера"
// @Router /advertisers/bulk [post]
func (h *AdvertiserHandler) CreateOrUpdate(c *gin.Context) {
	req := model.AdvertiserCreateOrUpdateRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, BuildError(err.Error()))
		return
	}

	advertiserIDs := make(map[uuid.UUID]bool)
	for _, advertiser := range req {
		if _, exists := advertiserIDs[advertiser.AdvertiserID]; exists {
			c.JSON(400, BuildError("You can't transmit the same advertiser_id"))
			return
		}
		advertiserIDs[advertiser.AdvertiserID] = true
	}

	currentDate, err := h.actCacheSvc.GetInt("current_day")
	if err != nil {
		HandleError500(c, err)
		return
	}

	err = h.advertiserSvc.CreateOrUpdate(req, int32(currentDate))
	if err != nil {
		HandleError500(c, err)
		return
	}

	c.JSON(201, req)
}

// GetByID godoc
// @Summary Получить информацию о рекламодателе по ID
// @Description Возвращает информацию о рекламодателе по указанному идентификатору
// @Tags Рекламодатели
// @Accept json
// @Produce json
// @Param advertiserID path string true "Идентификатор рекламодателя" Format(uuid)
// @Success 200 {object} model.Advertiser "Информация о рекламодателе"
// @Failure 400 {object} model.ErrorResponse "Неверный формат идентификатора рекламодателя"
// @Failure 404 {object} model.ErrorResponse "Рекламодатель не найден"
// @Failure 500 {object} model.ErrorResponse "Ошибка сервера"
// @Router /advertisers/{advertiserID} [get]
func (h *AdvertiserHandler) GetByID(c *gin.Context) {
	advertiserID, err := uuid.Parse(c.Param("advertiserID"))
	if err != nil {
		c.JSON(400, BuildError(err.Error()))
		c.Abort()
		return
	}

	advertiser, err := h.advertiserSvc.GetByID(advertiserID)
	if err != nil {
		c.JSON(404, BuildError("Advertiser not found"))
		return
	}

	c.JSON(200, advertiser)
}
