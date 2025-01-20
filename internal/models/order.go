package models

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	ID          uuid.UUID  `json:"id"`
	ClientEmail string     `json:"client_email"`
	CreatedAt   time.Time  `json:"created_at"`
	ConfirmedAt *time.Time `json:"confirmed_at"`
	SendedAt    *time.Time `json:"sended_at"`
	DeliveredAt *time.Time `json:"delivered_at"`
}

func NewOrder(clientEmail string) Order {
	return Order{
		ID:          uuid.New(),
		ClientEmail: clientEmail,
		CreatedAt:   time.Now(),
	}
}
