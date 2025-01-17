package services

import "context"

type ProcessedOrderMessage string

type ProcessedOrderSender interface {
	SendMessage(ctx context.Context, msg ProcessedOrderMessage) error
}
