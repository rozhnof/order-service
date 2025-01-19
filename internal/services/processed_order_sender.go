package services

import (
	"context"

	"github.com/rozhnof/order-service/internal/models"
)

type ProcessedOrderMessage struct {
	Order models.Order
}

type ProcessedOrderSender interface {
	SendMessage(ctx context.Context, msg ProcessedOrderMessage) error
}
