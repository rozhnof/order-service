package publisher

import (
	"context"
	"log/slog"

	"github.com/gin-gonic/gin"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rozhnof/order-service/internal/app"
	"github.com/rozhnof/order-service/internal/handlers"
	"github.com/rozhnof/order-service/internal/pkg/config"
	"github.com/rozhnof/order-service/internal/pkg/postgres"
	"github.com/rozhnof/order-service/internal/pkg/rabbitmq"
	"github.com/rozhnof/order-service/internal/pkg/server"
	"github.com/rozhnof/order-service/internal/repository"
	"github.com/rozhnof/order-service/internal/services"
)

type PublisherApp struct {
	logger     *slog.Logger
	httpServer *server.HTTPServer
}

type Config struct {
	Logger   config.Logger `yaml:"logging" env-required:"true"`
	Server   config.Server `yaml:"server" env-required:"true"`
	RabbitMQ config.RabbitMQ
	Postgres config.Postgres
}

func NewPublisherApp(ctx context.Context, cfg Config, ch *amqp.Channel, logger *slog.Logger, db postgres.Database) (PublisherApp, error) {
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
		orderHandler = handlers.NewOrderHandler(orderService, logger)
	)

	router := gin.New()
	router.Use(
		LogMiddleware(logger),
	)
	InitRoutes(router, orderHandler)

	httpServer := server.NewHTTPServer(ctx, cfg.Server.Address, router, logger)

	return PublisherApp{
		logger:     logger,
		httpServer: httpServer,
	}, nil
}

func (a *PublisherApp) Run(ctx context.Context) error {
	return a.httpServer.Run(ctx)
}
