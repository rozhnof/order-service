package app

import (
	"github.com/gin-gonic/gin"
	"github.com/rozhnof/order-service/internal/handlers"
)

func InitRoutes(router gin.IRouter, orderHandler handlers.OrderHandler) {
	router.POST("/orders", orderHandler.CreateOrder)
}
