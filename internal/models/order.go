package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

const (
	Confirmed = iota
	Sended
	Delivered
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

func (o *Order) SetStatus(status int) error {
	if status == Confirmed {
		return o.Confirm()
	}

	if status == Sended {
		return o.Send()
	}

	if status == Delivered {
		return o.Deliver()
	}

	return errors.New("invalid status")
}

func (o *Order) Confirm() error {
	if o.DeliveredAt != nil {
		return errors.New("confirm error: order already delivered")
	}

	if o.SendedAt != nil {
		return errors.New("confirm error: order already sended")
	}

	now := time.Now()
	o.ConfirmedAt = &now

	return nil
}

func (o *Order) Send() error {
	if o.DeliveredAt != nil {
		return errors.New("confirm error: order already delivered")
	}

	now := time.Now()
	o.SendedAt = &now

	return nil
}

func (o *Order) Deliver() error {
	now := time.Now()
	o.DeliveredAt = &now

	return nil
}
