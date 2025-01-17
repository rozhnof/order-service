package app

import (
	"errors"

	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	CreatedOrderQueue   = "created_order"
	ProcessedOrderQueue = "created_processed"
	NotificationQueue   = "notification"
)

func InitQueues(ch *amqp.Channel) error {
	var (
		_, createdOrderErr   = ch.QueueDeclare(CreatedOrderQueue, true, false, false, false, nil)
		_, processedOrderErr = ch.QueueDeclare(ProcessedOrderQueue, true, false, false, false, nil)
		_, notificationErr   = ch.QueueDeclare(NotificationQueue, true, false, false, false, nil)
	)

	return errors.Join(createdOrderErr, processedOrderErr, notificationErr)
}
