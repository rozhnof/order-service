package rabbitmq

import (
	"context"
	"encoding/json"
	"log/slog"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer[T any] struct {
	channel   *amqp.Channel
	queueName string
	logger    *slog.Logger
}

func NewConsumer[T any](channel *amqp.Channel, queueName string, logger *slog.Logger) (Consumer[T], error) {
	return Consumer[T]{
		channel:   channel,
		queueName: queueName,
		logger:    logger,
	}, nil
}

func (c Consumer[T]) ConsumeMessages(ctx context.Context, f func(context.Context, T)) error {
	messages, err := c.channel.Consume(c.queueName, "", false, false, false, false, nil)
	if err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		case m := <-messages:
			var msg T
			if err := json.Unmarshal(m.Body, &msg); err != nil {
				c.logger.Error("failed unmarshal rabbitmq message", slog.String("error", err.Error()))
			}

			f(ctx, msg)
		}
	}
}
