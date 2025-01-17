package app

import (
	"errors"

	amqp "github.com/rabbitmq/amqp091-go"
)

type QueueSettings struct {
	Name       string
	Durable    bool
	AutoDelete bool
	Exclusive  bool
	NoWait     bool
	Args       amqp.Table
	AutoAck    bool
	NoLocal    bool
}

var (
	CreatedOrderQueue   QueueSettings
	ProcessedOrderQueue QueueSettings
	NotificationQueue   QueueSettings
)

func InitQueues(ch *amqp.Channel) error {
	CreatedOrderQueue = QueueSettings{
		Name:      "created_order",
		Durable:   true,
		Exclusive: false,
		NoWait:    false,
		Args:      nil,
		AutoAck:   false,
		NoLocal:   false,
	}

	ProcessedOrderQueue = QueueSettings{
		Name:      "processed_order",
		Durable:   true,
		Exclusive: false,
		NoWait:    false,
		Args:      nil,
		AutoAck:   false,
		NoLocal:   false,
	}

	NotificationQueue = QueueSettings{
		Name:      "notification",
		Durable:   true,
		Exclusive: false,
		NoWait:    false,
		Args:      nil,
		AutoAck:   false,
		NoLocal:   false,
	}

	var (
		_, createdOrderErr   = declareQueue(ch, CreatedOrderQueue)
		_, processedOrderErr = declareQueue(ch, ProcessedOrderQueue)
		_, notificationErr   = declareQueue(ch, NotificationQueue)
	)

	return errors.Join(createdOrderErr, processedOrderErr, notificationErr)
}

func declareQueue(ch *amqp.Channel, qs QueueSettings) (amqp.Queue, error) {
	return ch.QueueDeclare(
		qs.Name,
		qs.Durable,
		qs.AutoDelete,
		qs.Exclusive,
		qs.NoWait,
		qs.Args,
	)
}
