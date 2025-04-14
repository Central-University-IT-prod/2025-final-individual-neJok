package handler

import (
	"github.com/gin-gonic/gin"
	"neJok/solution/internal/model"
	"net/http"
)

func BuildError(message string) model.ErrorResponse {
	return model.ErrorResponse{Status: "error", Message: message}
}

func HandleError500(c *gin.Context, err error) {
	_ = c.Error(err)
	c.JSON(http.StatusInternalServerError, BuildError(err.Error()))
}
