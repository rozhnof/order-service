package services

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/rozhnof/order-service/internal/models"
)

type OrderRepository interface {
	Create(ctx context.Context, order models.Order) (models.Order, error)
}

type OrderService struct {
	repo                 OrderRepository
	createdOrderSender   CreatedOrderSender
	processedOrderSender ProcessedOrderSender
	notificationSender   NotificationSender
}

func NewOrderService(
	repo OrderRepository,
	createdOrderSender CreatedOrderSender,
	processedOrderSender ProcessedOrderSender,
	notificationSender NotificationSender,
) OrderService {
	return OrderService{
		repo:                 repo,
		createdOrderSender:   createdOrderSender,
		processedOrderSender: processedOrderSender,
		notificationSender:   notificationSender,
	}
}

func (s OrderService) CreateOrder(ctx context.Context, order models.Order) (models.Order, error) {
	order.ID = uuid.New()

	order, err := s.repo.Create(ctx, order)
	if err != nil {
		return models.Order{}, err
	}

	if err := s.createdOrderSender.SendMessage(ctx, "CreatedOrderMessage"); err != nil {
		slog.Warn("failed send CreatedOrderMessage", slog.String("error", err.Error()))
	}

	if err := s.processedOrderSender.SendMessage(ctx, "ProcessedOrderMessage"); err != nil {
		slog.Warn("failed send ProcessedOrderMessage", slog.String("error", err.Error()))
	}

	if err := s.notificationSender.SendMessage(ctx, "NotificationMessage"); err != nil {
		slog.Warn("failed send NotificationMessage", slog.String("error", err.Error()))
	}

	return order, nil
}
