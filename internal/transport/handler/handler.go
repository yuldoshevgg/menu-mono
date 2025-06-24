package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/yuldoshevgg/menu-mono/internal/config"
	"github.com/yuldoshevgg/menu-mono/internal/core/repository"
)

type Handler struct {
	Cfg  *config.Config
	Repo repository.Store
}

type Response struct {
	Data interface{} `json:"data"`
}

func NewHandler(params Handler) Handler {
	return params
}

func (*Handler) handleResponse(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, Response{
		Data: data,
	})
}
