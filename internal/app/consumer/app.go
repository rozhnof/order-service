package consumer

import (
	"context"
	"fmt"
	"log/slog"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rozhnof/order-service/internal/app"
	"github.com/rozhnof/order-service/internal/pkg/config"
	"github.com/rozhnof/order-service/internal/pkg/mail"
	"github.com/rozhnof/order-service/internal/pkg/postgres"
	"github.com/rozhnof/order-service/internal/pkg/rabbitmq"
	"github.com/rozhnof/order-service/internal/services"
)

type ConsumerApp struct {
	logger                 *slog.Logger
	createdOrderConsumer   rabbitmq.Consumer[services.CreatedOrderMessage]
	processedOrderConsumer rabbitmq.Consumer[services.ProcessedOrderMessage]
	notificationConsumer   rabbitmq.Consumer[services.NotificationMessage]
	mailSender             mail.Sender
}

type Config struct {
	Logger   config.Logger `yaml:"logging" env-required:"true"`
	RabbitMQ config.RabbitMQ
	Mail     config.Mail
}

func NewConsumerApp(ctx context.Context, cfg Config, ch *amqp.Channel, logger *slog.Logger, db postgres.Database) (ConsumerApp, error) {
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

	mailSender, err := mail.NewSender(cfg.Mail.Email, cfg.Mail.Password, cfg.Mail.Address, cfg.Mail.Port)
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
	go a.createdOrderConsumer.ConsumeMessages(ctx, a.HandleCreatedOrderMessage)
	go a.processedOrderConsumer.ConsumeMessages(ctx, a.HandleProcessedOrderMessage)
	go a.notificationConsumer.ConsumeMessages(ctx, a.HandleNotificationMessage)
}

func (a *ConsumerApp) HandleCreatedOrderMessage(ctx context.Context, msg services.CreatedOrderMessage) {
	fmt.Println("cunsumed message from created_order queue:", msg)
}

func (a *ConsumerApp) HandleProcessedOrderMessage(ctx context.Context, msg services.ProcessedOrderMessage) {
	fmt.Println("cunsumed message from processed_order queue", msg)
}

func (a *ConsumerApp) HandleNotificationMessage(ctx context.Context, msg services.NotificationMessage) {
	fmt.Println("cunsumed message from notification queue", msg)

	mailMsg := mail.Message{
		Receiver: msg.ClientEmail,
		Subject:  msg.Subject,
		Body:     msg.Message,
	}

	if err := a.mailSender.SendMessage(mailMsg); err != nil {
		a.logger.Error("failed send email message", slog.String("error", err.Error()))
	}
}
