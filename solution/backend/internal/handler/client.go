package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"neJok/solution/internal/model"
	"neJok/solution/internal/service"
)

type ClientHandler struct {
	clientSvc *service.ClientService
}

func NewClientHandler(clientSvc *service.ClientService) *ClientHandler {
	return &ClientHandler{clientSvc}
}

// CreateOrUpdate godoc
// @Summary Создать или обновить клиента
// @Description Создает или обновляет информацию о клиентах, переданных в запросе
// @Tags Клиенты
// @Accept json
// @Produce json
// @Param clients body model.ClientCreateOrUpdateRequest true "Список клиентов для создания или обновления"
// @Success 201 {object} model.ClientCreateOrUpdateRequest "Информация о клиентах, которые были успешно созданы или обновлены"
// @Failure 400 {object} model.ErrorResponse "Невалидные данные или повторяющиеся идентификаторы клиентов"
// @Failure 500 {object} model.ErrorResponse "Ошибка сервера"
// @Router /clients/bulk [post]
func (h *ClientHandler) CreateOrUpdate(c *gin.Context) {
	req := model.ClientCreateOrUpdateRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, BuildError(err.Error()))
		return
	}

	clientIDs := make(map[uuid.UUID]bool)
	for _, client := range req {
		if _, ok := clientIDs[client.ClientID]; ok {
			c.JSON(400, BuildError("You can't transmit the same client_id"))
			return
		}

		clientIDs[client.ClientID] = true
	}

	err := h.clientSvc.CreateOrUpdate(req)
	if err != nil {
		HandleError500(c, err)
		return
	}

	c.JSON(201, req)
}

// GetByID godoc
// @Summary Получить информацию о клиенте по ID
// @Description Возвращает информацию о клиенте по указанному идентификатору
// @Tags Клиенты
// @Accept json
// @Produce json
// @Param clientID path string true "Идентификатор клиента" Format(uuid)
// @Success 200 {object} model.Client "Информация о клиенте"
// @Failure 400 {object} model.ErrorResponse "Неверный формат идентификатора клиента"
// @Failure 404 {object} model.ErrorResponse "Клиент не найден"
// @Failure 500 {object} model.ErrorResponse "Ошибка сервера"
// @Router /clients/{clientID} [get]
func (h *ClientHandler) GetByID(c *gin.Context) {
	clientID, err := uuid.Parse(c.Param("clientID"))
	if err != nil {
		c.JSON(400, BuildError(err.Error()))
		return
	}

	client, err := h.clientSvc.GetByID(clientID)
	if err != nil {
		c.JSON(404, BuildError("Client not found"))
		return
	}

	c.JSON(200, client)
}
