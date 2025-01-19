package services

import (
	"context"

	"github.com/rozhnof/order-service/internal/models"
)

type CreatedOrderMessage struct {
	Order models.Order
}

type CreatedOrderSender interface {
	SendMessage(ctx context.Context, msg CreatedOrderMessage) error
}
