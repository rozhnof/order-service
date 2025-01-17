package services

import "context"

type CreatedOrderMessage string

type CreatedOrderSender interface {
	SendMessage(ctx context.Context, msg CreatedOrderMessage) error
}
