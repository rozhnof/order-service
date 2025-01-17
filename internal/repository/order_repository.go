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
	const query = `INSERT INTO orders(id) VALUES ($1)`

	if _, err := r.db.Exec(ctx, query, order.ID); err != nil {
		return models.Order{}, err
	}

	return order, nil
}
