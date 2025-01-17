package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rozhnof/order-service/internal/app"
	"github.com/rozhnof/order-service/internal/pkg/postgres"
	"github.com/rozhnof/order-service/internal/services"
)

type ConsumerApp struct {
	logger                 *slog.Logger
	createdOrderMessages   <-chan amqp.Delivery
	processedOrderMessages <-chan amqp.Delivery
	notificationMessages   <-chan amqp.Delivery
}

func NewConsumerApp(ctx context.Context, ch *amqp.Channel, logger *slog.Logger, db postgres.Database) (ConsumerApp, error) {
	if err := app.InitQueues(ch); err != nil {
		return ConsumerApp{}, err
	}

	createdOrderMessages, err := Consume(ch, app.CreatedOrderQueue)
	if err != nil {
		return ConsumerApp{}, err
	}

	processedOrderMessages, err := Consume(ch, app.ProcessedOrderQueue)
	if err != nil {
		return ConsumerApp{}, err
	}

	notificationMessages, err := Consume(ch, app.NotificationQueue)
	if err != nil {
		return ConsumerApp{}, err
	}

	return ConsumerApp{
		logger:                 logger,
		createdOrderMessages:   createdOrderMessages,
		processedOrderMessages: processedOrderMessages,
		notificationMessages:   notificationMessages,
	}, nil
}

func Consume(ch *amqp.Channel, qs app.QueueSettings) (<-chan amqp.Delivery, error) {
	return ch.Consume(qs.Name, "", qs.AutoAck, qs.Exclusive, qs.NoLocal, qs.NoWait, qs.Args)
}

func (a *ConsumerApp) Start(ctx context.Context) {
	go RunProcess(ctx, a.logger, a.createdOrderMessages, ProcessCreatedOrderMessage)
	go RunProcess(ctx, a.logger, a.processedOrderMessages, ProcessProcessedOrderMessage)
	go RunProcess(ctx, a.logger, a.notificationMessages, ProcessNotificationMessage)
}

func RunProcess[T any](ctx context.Context, logger *slog.Logger, messages <-chan amqp.Delivery, f func(m T)) {
	for {
		select {
		case <-ctx.Done():
			return
		case m := <-messages:
			var msg T
			if err := json.Unmarshal(m.Body, &msg); err != nil {
				logger.Error("failed unmarshal rabbitmq message", slog.String("error", err.Error()))
			}

			f(msg)
		}
	}
}

func ProcessCreatedOrderMessage(msg services.CreatedOrderMessage) {
	fmt.Println("cunsumed CreatedOrderMessage", msg)
}

func ProcessProcessedOrderMessage(msg services.ProcessedOrderMessage) {
	fmt.Println("cunsumed ProcessedOrderMessage", msg)
}

func ProcessNotificationMessage(msg services.NotificationMessage) {
	fmt.Println("cunsumed NotificationMessage", msg)
}
