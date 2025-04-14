package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"neJok/solution/internal/model"
	"neJok/solution/internal/service"
	"path/filepath"
	"strconv"
	"strings"
)

type CampaignHandler struct {
	campaignSvc   *service.CampaignService
	advertiserSvc *service.AdvertiserService
	actCacheSvc   *service.ActCacheService
	s3Svc         *service.S3Service
	openAISvc     *service.OpenAIService
}

func NewCampaignHandler(campaignSvc *service.CampaignService, advertiserSvc *service.AdvertiserService, actCacheSvc *service.ActCacheService, s3Svc *service.S3Service, openAISvc *service.OpenAIService) *CampaignHandler {
	return &CampaignHandler{campaignSvc, advertiserSvc, actCacheSvc, s3Svc, openAISvc}
}

// Create godoc
// @Summary Создать кампанию
// @Description Создает новую кампанию для рекламодателя, с проверкой параметров, таких как даты, таргетинг, изображение.
// @Tags Кампании
// @Accept application/json, multipart/form-data
// @Produce json
// @Param advertiserID path string true "Идентификатор рекламодателя" Format(uuid)
// @Param image_file formData file false "Изображение для кампании (Доступно только с multipart/form-data)"
// @Param campaign body model.CampaignCreate true "Данные кампании"
// @Success 201 {object} model.Campaign "Созданная кампания"
// @Failure 400 {object} model.ErrorResponse "Ошибки в данных запроса (неверные даты, некорректный таргетинг, ошибка с изображением)"
// @Failure 422 {object} model.ErrorResponse "Текст или название рекламной кампании не прошли модерацию"
// @Failure 404 {object} model.ErrorResponse "Рекламодатель не найден"
// @Failure 500 {object} model.ErrorResponse "Ошибка сервера"
// @Router /advertisers/{advertiserID}/campaigns [post]
func (h *CampaignHandler) Create(c *gin.Context) {
	advertiserID, err := uuid.Parse(c.Param("advertiserID"))
	if err != nil {
		c.JSON(400, BuildError(err.Error()))
		return
	}

	req := model.CampaignCreate{}
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(400, BuildError(err.Error()))
		return
	}

	if (req.Targeting.AgeFrom != nil && req.Targeting.AgeTo != nil) && (*req.Targeting.AgeFrom > *req.Targeting.AgeTo) {
		c.JSON(400, BuildError("Age from too high"))
		return
	}

	_, err = h.advertiserSvc.GetByID(advertiserID)
	if err != nil {
		c.JSON(404, BuildError("Advertiser not found"))
		return
	}
	req.Targeting.SetDefaults()

	currentDate, err := h.actCacheSvc.GetInt("current_day")
	if err != nil {
		HandleError500(c, err)
		return
	}

	if *req.StartDate < int32(currentDate) || *req.EndDate < int32(currentDate) {
		c.JSON(400, BuildError("Bad date values"))
		return
	}

	campaignID := uuid.New()
	if req.ImageFile != nil && req.ImageURL != nil {
		c.JSON(400, BuildError("Both file and URL are provided. Only one is allowed"))
		return
	}

	moderation, err := h.actCacheSvc.GetStr("moderation")
	if err == nil && moderation == "true" {
		moderated, err := h.openAISvc.ModerateText(req.AdText + "\n" + req.AdTitle)
		if err != nil {
			c.JSON(500, BuildError("Failed to moderate text\n"+err.Error()))
			return
		}
		if !moderated {
			c.JSON(422, BuildError("Your text did not pass our moderation :("))
			return
		}
	}

	if req.ImageFile != nil {
		if req.ImageFile.Size > 25*1024*1024 { // 25 MB limit
			c.JSON(400, BuildError("File size exceeds 25 MB limit"))
			return
		}

		ext := strings.ToLower(filepath.Ext(req.ImageFile.Filename))
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".gif" && ext != ".webp" {
			c.JSON(400, BuildError("Invalid file type"))
			return
		}

		file, err := req.ImageFile.Open()
		if err != nil {
			c.JSON(400, BuildError("Failed to open image file"))
			return
		}
		defer file.Close()

		fileName := fmt.Sprintf("%s%s", campaignID.String(), ext)
		fileURL, err := h.s3Svc.UploadToS3(file, fileName)
		if err != nil {
			c.JSON(500, BuildError("Failed to upload image"))
			return
		}
		req.ImageURL = &fileURL
		req.FileName = &fileName
	}

	campaign, err := h.campaignSvc.Add(advertiserID, req, campaignID)
	if err != nil {
		HandleError500(c, err)
		return
	}
	go func() {
		h.actCacheSvc.DeleteKeysByPrefix("top_ads:")
	}()
	c.JSON(201, campaign)
}

// GetMany godoc
// @Summary Получить список кампаний рекламодателя
// @Description Возвращает список кампаний для указанного рекламодателя с возможностью пагинации.
// @Tags Кампании
// @Accept json
// @Produce json
// @Param advertiserID path string true "Идентификатор рекламодателя" Format(uuid)
// @Param size query int false "Количество кампаний на страницу" default(10)
// @Param page query int false "Номер страницы начиная с 0" default(0)
// @Success 200 {array} model.Campaign "Список кампаний"
// @Failure 400 {object} model.ErrorResponse "Ошибки в параметрах запроса"
// @Failure 404 {object} model.ErrorResponse "Рекламодатель не найден"
// @Failure 500 {object} model.ErrorResponse "Ошибка сервера"
// @Router /advertisers/{advertiserID}/campaigns [get]
func (h *CampaignHandler) GetMany(c *gin.Context) {
	advertiserID, err := uuid.Parse(c.Param("advertiserID"))
	if err != nil {
		c.JSON(400, BuildError(err.Error()))
		return
	}

	var req model.CampaignGetManyRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(400, BuildError(err.Error()))
		return
	}

	_, err = h.advertiserSvc.GetByID(advertiserID)
	if err != nil {
		c.JSON(404, BuildError("Advertiser not found"))
		return
	}
	req.Limit, _ = strconv.Atoi(c.DefaultQuery("size", "10"))

	total, campaigns, err := h.campaignSvc.GetMany(advertiserID, req.Limit, req.Offset*req.Limit)
	if err != nil {
		HandleError500(c, err)
		return
	}

	if req.Limit == 0 {
		campaigns = []model.Campaign{}
	}

	c.Header("X-Total-Count", strconv.Itoa(int(total)))
	if len(campaigns) == 0 {
		c.JSON(200, []model.Campaign{})
		return
	}
	c.JSON(200, campaigns)
}

// GetOne godoc
// @Summary Получить информацию о кампании по ID
// @Description Возвращает информацию о кампании для указанного рекламодателя по его ID и ID кампании
// @Tags Кампании
// @Accept json
// @Produce json
// @Param advertiserID path string true "Идентификатор рекламодателя" Format(uuid)
// @Param campaignID path string true "Идентификатор кампании" Format(uuid)
// @Success 200 {object} model.Campaign "Информация о кампании"
// @Failure 400 {object} model.ErrorResponse "Неверные идентификаторы рекламодателя или кампании"
// @Failure 404 {object} model.ErrorResponse "Кампания не найдена"
// @Failure 500 {object} model.ErrorResponse "Ошибка сервера"
// @Router /advertisers/{advertiserID}/campaigns/{campaignID} [get]
func (h *CampaignHandler) GetOne(c *gin.Context) {
	advertiserID, err := uuid.Parse(c.Param("advertiserID"))
	if err != nil {
		c.JSON(400, BuildError(err.Error()))
		return
	}

	campaignID, err := uuid.Parse(c.Param("campaignID"))
	if err != nil {
		c.JSON(400, BuildError(err.Error()))
		return
	}

	campaign, err := h.campaignSvc.GetByID(campaignID)
	if err != nil || campaign.AdvertiserID != advertiserID {
		c.JSON(404, BuildError("Campaign not found"))
		return
	}
	c.JSON(200, campaign)
}

// DeleteOne godoc
// @Summary Удалить кампанию
// @Description Удаляет кампанию по ID для указанного рекламодателя
// @Tags Кампании
// @Accept json
// @Produce json
// @Param advertiserID path string true "Идентификатор рекламодателя" Format(uuid)
// @Param campaignID path string true "Идентификатор кампании" Format(uuid)
// @Success 204 "Кампания успешно удалена"
// @Failure 400 {object} model.ErrorResponse "Неверные идентификаторы рекламодателя или кампании"
// @Failure 404 {object} model.ErrorResponse "Кампания не найдена"
// @Failure 500 {object} model.ErrorResponse "Ошибка сервера"
// @Router /advertisers/{advertiserID}/campaigns/{campaignID} [delete]
func (h *CampaignHandler) DeleteOne(c *gin.Context) {
	advertiserID, err := uuid.Parse(c.Param("advertiserID"))
	if err != nil {
		c.JSON(400, BuildError(err.Error()))
		return
	}

	campaignID, err := uuid.Parse(c.Param("campaignID"))
	if err != nil {
		c.JSON(400, BuildError(err.Error()))
		return
	}

	campaign, err := h.campaignSvc.DeleteByID(advertiserID, campaignID)
	if err != nil {
		c.JSON(404, BuildError("Campaign not found"))
		return
	}

	go func() {
		h.actCacheSvc.DeleteKeysByPrefix("top_ads:")
		if campaign.FileName != nil {
			h.s3Svc.Delete(*campaign.FileName)
		}
	}()
	c.Status(204)
}

// UpdateOne godoc
// @Summary Обновить информацию о кампании
// @Description Обновляет данные кампании для указанного рекламодателя и кампании, с учетом проверок на невозможность изменения некоторых полей после запуска кампании. Нельзя передать одновременно image_url и image_file. Скорее всего со swagger не получиться отправить запрос с image_file, используйте для тестирования этого postman или альтернативы.
// @Tags Кампании
// @Accept application/json, multipart/form-data
// @Produce json
// @Param advertiserID path string true "Идентификатор рекламодателя" Format(uuid)
// @Param campaignID path string true "Идентификатор кампании" Format(uuid)
// @Param campaign body model.CampaignUpdate true "Обновленные данные кампании"
// @Param image_file formData file false "Изображение для кампании (Доступно только с multipart/form-data)"
// @Success 200 {object} model.Campaign "Обновленная кампания"
// @Failure 400 {object} model.ErrorResponse "Ошибки в данных запроса (например, некорректные даты или поля, которые нельзя изменить после запуска)"
// @Failure 422 {object} model.ErrorResponse "Текст или название рекламной кампании не прошли модерацию"
// @Failure 404 {object} model.ErrorResponse "Кампания не найдена"
// @Failure 500 {object} model.ErrorResponse "Ошибка сервера"
// @Router /advertisers/{advertiserID}/campaigns/{campaignID} [put]
func (h *CampaignHandler) UpdateOne(c *gin.Context) {
	advertiserID, err := uuid.Parse(c.Param("advertiserID"))
	if err != nil {
		c.JSON(400, BuildError(err.Error()))
		return
	}

	campaignID, err := uuid.Parse(c.Param("campaignID"))
	if err != nil {
		c.JSON(400, BuildError(err.Error()))
		return
	}

	req := model.CampaignUpdate{}
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(400, BuildError(err.Error()))
		return
	}

	campaign, err := h.campaignSvc.GetByID(campaignID)
	if err != nil || campaign.AdvertiserID != advertiserID {
		c.JSON(404, BuildError("Campaign not found"))
		return
	}

	currentDate, err := h.actCacheSvc.GetInt("current_day")
	if err != nil {
		HandleError500(c, err)
		return
	}

	int32CurrentDate := int32(currentDate)
	if int32CurrentDate >= *campaign.StartDate && (req.ImpressionsLimit != campaign.ImpressionsLimit || req.ClicksLimit != campaign.ClicksLimit || req.StartDate != campaign.StartDate || req.EndDate != campaign.EndDate) {
		c.JSON(400, BuildError("You can't edit some fields after the ad is launched."))
		return
	}

	if (*req.StartDate < int32CurrentDate) || (*req.EndDate < int32CurrentDate) {
		c.JSON(400, BuildError("Start date and end date cannot be less than current date."))
		return
	}

	if req.ImageFile != nil && req.ImageURL != nil {
		c.JSON(400, BuildError("Both file and URL are provided. Only one is allowed"))
		return
	}

	moderation, err := h.actCacheSvc.GetStr("moderation")
	if err == nil && moderation == "true" {
		moderated, err := h.openAISvc.ModerateText(req.AdText + "\n" + req.AdTitle)
		if err != nil {
			c.JSON(500, BuildError("Failed to moderate text\n"+err.Error()))
			return
		}
		if !moderated {
			c.JSON(422, BuildError("Your text did not pass our moderation :("))
			return
		}
	}

	if req.ImageFile != nil {
		if req.ImageFile.Size > 25*1024*1024 { // 25 MB limit
			c.JSON(400, BuildError("File size exceeds 25 MB limit"))
			return
		}

		ext := strings.ToLower(filepath.Ext(req.ImageFile.Filename))
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".gif" && ext != ".webp" {
			c.JSON(400, BuildError("Invalid file type"))
			return
		}

		file, err := req.ImageFile.Open()
		if err != nil {
			c.JSON(400, BuildError("Failed to open image file"))
			return
		}
		defer file.Close()

		fileName := fmt.Sprintf("%s%s", campaignID.String(), ext)
		fileURL, err := h.s3Svc.UploadToS3(file, fileName)
		if err != nil {
			c.JSON(500, BuildError("Failed to upload image"))
			return
		}
		req.ImageURL = &fileURL
		req.FileName = &fileName
	} else if campaign.FileName != nil {
		h.s3Svc.Delete(*campaign.FileName)
	}

	campaign, err = h.campaignSvc.UpdateByID(advertiserID, campaignID, req)
	if err != nil {
		c.JSON(400, BuildError(err.Error()))
		return
	}

	go func() {
		h.actCacheSvc.DeleteKeysByPrefix("top_ads:")
	}()
	c.JSON(200, campaign)
}
