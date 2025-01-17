package services

import "context"

type NotificationMessage struct {
	ClientEmail string
	Subject     string
	Message     string
}

type NotificationSender interface {
	SendMessage(ctx context.Context, msg NotificationMessage) error
}
