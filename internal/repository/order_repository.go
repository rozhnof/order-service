package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/rozhnof/order-service/internal/models"
	"github.com/rozhnof/order-service/internal/pkg/postgres"
)

var (
	ErrObjectNotFound = errors.New("object not found")
	ErrDuplicate      = errors.New("object is duplicate")
)

type OrderRepository struct {
	db postgres.Database
}

func NewOrderRepository(db postgres.Database) OrderRepository {
	return OrderRepository{
		db: db,
	}
}

func (r OrderRepository) Create(ctx context.Context, order models.Order) (models.Order, error) {
	var clientID uuid.UUID
	if err := r.db.QueryRow(ctx, createClientQuery, order.ClientEmail).Scan(&clientID); err != nil {
		return models.Order{}, err
	}

	if err := r.db.QueryRow(ctx, createOrderQuery, clientID).Scan(&order.ID); err != nil {
		return models.Order{}, err
	}

	return order, nil
}

func (r OrderRepository) GetByID(ctx context.Context, orderID uuid.UUID) (models.Order, error) {
	var (
		order    models.Order
		clientID uuid.UUID
	)
	if err := r.db.QueryRow(ctx, getOrderQuery, orderID).Scan(
		&clientID,
		&order.CreatedAt,
		&order.ConfirmedAt,
		&order.SendedAt,
		&order.DeliveredAt,
	); err != nil {
		return models.Order{}, err
	}

	var clientEmail string
	if err := r.db.QueryRow(ctx, getClientQuery, clientID).Scan(&clientEmail); err != nil {
		return models.Order{}, err
	}

	return order, nil
}

func (r OrderRepository) Update(ctx context.Context, order models.Order) error {
	var clientID uuid.UUID
	if err := r.db.QueryRow(ctx, createClientQuery, order.ClientEmail).Scan(&clientID); err != nil {
		return err
	}

	if _, err := r.db.Exec(ctx, updateOrderQuery, clientID, order.ConfirmedAt, order.SendedAt, order.DeliveredAt); err != nil {
		return err
	}

	return nil
}

func (r OrderRepository) Delete(ctx context.Context, order models.Order) error {
	if _, err := r.db.Exec(ctx, deleteOrderQuery, order.ID); err != nil {
		return err
	}

	return nil
}
