package consumer

import (
	"context"
	"fmt"
	"log/slog"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rozhnof/order-service/internal/app"
	"github.com/rozhnof/order-service/internal/pkg/postgres"
	"github.com/rozhnof/order-service/internal/pkg/rabbitmq"
	"github.com/rozhnof/order-service/internal/services"
)

type ConsumerApp struct {
	logger                 *slog.Logger
	createdOrderConsumer   rabbitmq.Consumer[services.CreatedOrderMessage]
	processedOrderConsumer rabbitmq.Consumer[services.ProcessedOrderMessage]
	notificationConsumer   rabbitmq.Consumer[services.NotificationMessage]
}

func NewConsumerApp(ctx context.Context, ch *amqp.Channel, logger *slog.Logger, db postgres.Database) (ConsumerApp, error) {
	if err := app.InitQueues(ch); err != nil {
		return ConsumerApp{}, err
	}

	createdOrderConsumer, err := rabbitmq.NewConsumer[services.CreatedOrderMessage](ch, app.CreatedOrderQueue, logger)
	if err != nil {
		return ConsumerApp{}, err
	}

	processedOrderConsumer, err := rabbitmq.NewConsumer[services.ProcessedOrderMessage](ch, app.ProcessedOrderQueue, logger)
	if err != nil {
		return ConsumerApp{}, err
	}

	notificationConsumer, err := rabbitmq.NewConsumer[services.NotificationMessage](ch, app.NotificationQueue, logger)
	if err != nil {
		return ConsumerApp{}, err
	}

	return ConsumerApp{
		logger:                 logger,
		createdOrderConsumer:   createdOrderConsumer,
		processedOrderConsumer: processedOrderConsumer,
		notificationConsumer:   notificationConsumer,
	}, nil
}

func (a *ConsumerApp) Start(ctx context.Context) {
	go a.createdOrderConsumer.ConsumeMessages(ctx, ProcessCreatedOrderMessage)
	go a.processedOrderConsumer.ConsumeMessages(ctx, ProcessProcessedOrderMessage)
	go a.notificationConsumer.ConsumeMessages(ctx, ProcessNotificationMessage)
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
