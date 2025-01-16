package app

import (
	"context"
	"errors"
	"log/slog"

	"github.com/gin-gonic/gin"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rozhnof/order-service/internal/handlers"
	"github.com/rozhnof/order-service/internal/pkg/rabbitmq"
	"github.com/rozhnof/order-service/internal/pkg/server"
	"github.com/rozhnof/order-service/internal/repository"
	"github.com/rozhnof/order-service/internal/services"
)

const (
	createdOrderTopic   = "created_order"
	processedOrderTopic = "processed_order"
	notificationTopic   = "notification"
)

const (
	addr = ":8080"
)

type OrderProcessApp struct {
	logger     *slog.Logger
	httpServer *server.HTTPServer
}

func NewOrderProcessApp(ctx context.Context, ch *amqp.Channel, logger *slog.Logger) (OrderProcessApp, error) {
	if err := InitQueues(ch); err != nil {
		return OrderProcessApp{}, err
	}

	var (
		repo = repository.NewOrderRepository()
	)

	var (
		createdOrderSender   = rabbitmq.NewSender[services.CreatedOrderMessage](ch, createdOrderTopic)
		processedOrderSender = rabbitmq.NewSender[services.ProcessedOrderMessage](ch, processedOrderTopic)
		notificationSender   = rabbitmq.NewSender[services.NotificationMessage](ch, notificationTopic)
	)

	var (
		orderService = services.NewOrderService(repo, createdOrderSender, processedOrderSender, notificationSender)
		orderHandler = handlers.NewOrderHandler(orderService)
	)

	router := gin.New()
	InitRoutes(router, orderHandler)

	httpServer := server.NewHTTPServer(ctx, addr, router, logger)

	return OrderProcessApp{
		logger:     logger,
		httpServer: httpServer,
	}, nil
}

func (a *OrderProcessApp) Run(ctx context.Context) error {
	return a.httpServer.Run(ctx)
}

func InitQueues(ch *amqp.Channel) error {
	var (
		_, createdOrderErr   = ch.QueueDeclare(createdOrderTopic, true, false, false, false, nil)
		_, processedOrderErr = ch.QueueDeclare(processedOrderTopic, true, false, false, false, nil)
		_, notificationErr   = ch.QueueDeclare(notificationTopic, true, false, false, false, nil)
	)

	return errors.Join(createdOrderErr, processedOrderErr, notificationErr)
}
