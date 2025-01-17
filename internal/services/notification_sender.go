package services

import "context"

type NotificationMessage string

type NotificationSender interface {
	SendMessage(ctx context.Context, msg NotificationMessage) error
}
