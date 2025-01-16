package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/rozhnof/order-service/internal/models"
)

var (
	ErrObjectNotFound = errors.New("object not found")
	ErrDuplicate      = errors.New("object is duplicate")
)

type OrderRepository struct{}

func NewOrderRepository() OrderRepository {
	return OrderRepository{}
}

func (r OrderRepository) Create(ctx context.Context, order models.Order) (models.Order, error) {
	return models.Order{
		ID: uuid.New(),
	}, nil
}
