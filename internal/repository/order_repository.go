package repository

import (
	"context"
	"errors"

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
	if _, err := r.db.Exec(ctx, createOrderQuery, order.ID); err != nil {
		return models.Order{}, err
	}

	return order, nil
}
