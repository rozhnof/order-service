package consumer

import (
	"context"
	"fmt"
	"log/slog"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rozhnof/order-service/internal/app"
	"github.com/rozhnof/order-service/internal/pkg/mail"
	"github.com/rozhnof/order-service/internal/pkg/postgres"
	"github.com/rozhnof/order-service/internal/pkg/rabbitmq"
	"github.com/rozhnof/order-service/internal/services"
)

const (
	email    = "golang.auth.service@gmail.com"
	password = "jybh ayjb qosq kykn"
	address  = "smtp.gmail.com"
	port     = "587"
)

type ConsumerApp struct {
	logger                 *slog.Logger
	createdOrderConsumer   rabbitmq.Consumer[services.CreatedOrderMessage]
	processedOrderConsumer rabbitmq.Consumer[services.ProcessedOrderMessage]
	notificationConsumer   rabbitmq.Consumer[services.NotificationMessage]
	mailSender             mail.Sender
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

	mailSender, err := mail.NewSender(email, password, address, port)
	if err != nil {
		return ConsumerApp{}, err
	}

	return ConsumerApp{
		logger:                 logger,
		createdOrderConsumer:   createdOrderConsumer,
		processedOrderConsumer: processedOrderConsumer,
		notificationConsumer:   notificationConsumer,
		mailSender:             mailSender,
	}, nil
}

func (a *ConsumerApp) Start(ctx context.Context) {
	go a.createdOrderConsumer.ConsumeMessages(ctx, a.ProcessCreatedOrderMessage)
	go a.processedOrderConsumer.ConsumeMessages(ctx, a.ProcessProcessedOrderMessage)
	go a.notificationConsumer.ConsumeMessages(ctx, a.ProcessNotificationMessage)
}

func (a *ConsumerApp) ProcessCreatedOrderMessage(msg services.CreatedOrderMessage) {
	fmt.Println("cunsumed CreatedOrderMessage", msg)
}

func (a *ConsumerApp) ProcessProcessedOrderMessage(msg services.ProcessedOrderMessage) {
	fmt.Println("cunsumed ProcessedOrderMessage", msg)
}

func (a *ConsumerApp) ProcessNotificationMessage(msg services.NotificationMessage) {
	fmt.Println("cunsumed NotificationMessage", msg)

	mailMsg := mail.Message{
		Receiver: msg.ClientEmail,
		Subject:  msg.Subject,
		Body:     msg.Message,
	}

	if err := a.mailSender.SendMessage(mailMsg); err != nil {
		a.logger.Error("failed send email message", slog.String("error", err.Error()))
	}
}
