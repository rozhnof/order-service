package rabbitmq

import (
	"context"
	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Sender[T any] struct {
	channel   *amqp.Channel
	queueName string
}

func NewSender[T any](channel *amqp.Channel, queueName string) Sender[T] {
	return Sender[T]{
		channel:   channel,
		queueName: queueName,
	}
}

func (s Sender[T]) SendMessage(ctx context.Context, msg T) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	rabbitMsg := amqp.Publishing{
		ContentType: "application/json",
		Body:        data,
	}

	if err := s.channel.Publish("", s.queueName, false, false, rabbitMsg); err != nil {
		return err
	}

	return nil
}
