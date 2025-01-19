package services

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/rozhnof/order-service/internal/models"
)

type OrderRepository interface {
	Create(ctx context.Context, order models.Order) (models.Order, error)
	GetByID(ctx context.Context, orderID uuid.UUID) (models.Order, error)
	Update(ctx context.Context, order models.Order) error
	Delete(ctx context.Context, order models.Order) error
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

func (s OrderService) CreateOrder(ctx context.Context, clientEmail string) (models.Order, error) {
	order := models.NewOrder(clientEmail)

	order, err := s.repo.Create(ctx, order)
	if err != nil {
		return models.Order{}, err
	}

	createdOrderMessage := CreatedOrderMessage{
		Order: order,
	}
	if err := s.createdOrderSender.SendMessage(ctx, createdOrderMessage); err != nil {
		slog.Warn("failed send CreatedOrderMessage", slog.String("error", err.Error()))
	}

	processedOrderMessage := ProcessedOrderMessage{
		Order: order,
	}
	if err := s.processedOrderSender.SendMessage(ctx, processedOrderMessage); err != nil {
		slog.Warn("failed send ProcessedOrderMessage", slog.String("error", err.Error()))
	}

	notificationMessage := NotificationMessage{
		ClientEmail: order.ClientEmail,
		Subject:     "Create Order",
		Message:     "Order successfuly created",
	}
	if err := s.notificationSender.SendMessage(ctx, notificationMessage); err != nil {
		slog.Warn("failed send NotificationMessage", slog.String("error", err.Error()))
	}

	return order, nil
}
