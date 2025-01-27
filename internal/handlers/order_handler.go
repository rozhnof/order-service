package handlers

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rozhnof/order-service/internal/models"
	"github.com/rozhnof/order-service/internal/repository"
	"github.com/rozhnof/order-service/internal/services"
)

type OrderHandler struct {
	service services.OrderService
	logger  *slog.Logger
}

func NewOrderHandler(service services.OrderService, logger *slog.Logger) OrderHandler {
	return OrderHandler{
		service: service,
		logger:  logger,
	}
}

type CreateOrderRequest struct {
	ClientEmail string `json:"client_email"`
}

type CreateOrderResponse struct {
	Order models.Order `json:"order"`
}

func (h OrderHandler) CreateOrder(c *gin.Context) {
	var request CreateOrderRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.String(http.StatusBadRequest, "invalid request")
		return
	}

	createdOrder, err := h.service.CreateOrder(c.Request.Context(), request.ClientEmail)
	if err != nil {
		if errors.Is(err, repository.ErrObjectNotFound) {
			c.String(http.StatusNotFound, err.Error())
			return
		}

		h.logger.Warn("failed create order", slog.String("error", err.Error()))

		c.Status(http.StatusInternalServerError)
		return
	}

	response := CreateOrderResponse{
		Order: createdOrder,
	}

	c.JSON(http.StatusOK, response)
}
