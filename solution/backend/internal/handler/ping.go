package handler

import (
	"github.com/gin-gonic/gin"
	"neJok/solution/internal/model"
)

type PingHandler struct{}

func NewPingHandler() *PingHandler {
	return &PingHandler{}
}

// Ping godoc
// @Summary Проверка работоспособности сервера
// @Description Проверка, что сервер работает и отвечает корректно
// @Tags ping
// @Accept json
// @Produce json
// @Success 200 {object} model.PingResponse "Ответ с состоянием сервера"
// @Router /ping [get]
func (h *PingHandler) Ping(c *gin.Context) {
	c.JSON(200, model.PingResponse{Status: "ok"})
}
