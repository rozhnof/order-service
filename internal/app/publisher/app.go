package publisher

import (
	"context"
	"log/slog"

	"github.com/gin-gonic/gin"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rozhnof/order-service/internal/app"
	"github.com/rozhnof/order-service/internal/handlers"
	"github.com/rozhnof/order-service/internal/pkg/postgres"
	"github.com/rozhnof/order-service/internal/pkg/rabbitmq"
	"github.com/rozhnof/order-service/internal/pkg/server"
	"github.com/rozhnof/order-service/internal/repository"
	"github.com/rozhnof/order-service/internal/services"
)

const (
	addr = ":8080"
)

type PublisherApp struct {
	logger     *slog.Logger
	httpServer *server.HTTPServer
}

func NewPublisherApp(ctx context.Context, ch *amqp.Channel, logger *slog.Logger, db postgres.Database) (PublisherApp, error) {
	if err := app.InitQueues(ch); err != nil {
		return PublisherApp{}, err
	}

	var (
		repo = repository.NewOrderRepository(db)
	)

	var (
		createdOrderSender   = rabbitmq.NewSender[services.CreatedOrderMessage](ch, app.CreatedOrderQueue)
		processedOrderSender = rabbitmq.NewSender[services.ProcessedOrderMessage](ch, app.ProcessedOrderQueue)
		notificationSender   = rabbitmq.NewSender[services.NotificationMessage](ch, app.NotificationQueue)
	)

	var (
		orderService = services.NewOrderService(repo, createdOrderSender, processedOrderSender, notificationSender)
		orderHandler = handlers.NewOrderHandler(orderService)
	)

	router := gin.New()
	InitRoutes(router, orderHandler)

	httpServer := server.NewHTTPServer(ctx, addr, router, logger)

	return PublisherApp{
		logger:     logger,
		httpServer: httpServer,
	}, nil
}

func (a *PublisherApp) Run(ctx context.Context) error {
	return a.httpServer.Run(ctx)
}
