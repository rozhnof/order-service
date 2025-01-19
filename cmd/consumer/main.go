package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/rozhnof/order-service/internal/app"
	"github.com/rozhnof/order-service/internal/app/consumer"
	"github.com/rozhnof/order-service/internal/pkg/config"
	"github.com/rozhnof/order-service/internal/pkg/postgres"
	"github.com/rozhnof/order-service/internal/pkg/rabbitmq"
)

const (
	EnvConfigPath = "CONFIG_PATH"
)

func main() {
	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	cfg, err := config.NewConfig[consumer.Config](os.Getenv(EnvConfigPath))
	if err != nil {
		slog.Error("init config failed", slog.String("error", err.Error()))
		return
	}

	logger, err := app.NewLogger(cfg.Logger)
	if err != nil {
		slog.Error("init logger failed", slog.String("error", err.Error()))
		return
	}

	rabbitURL := fmt.Sprintf(rabbitmq.URL, cfg.RabbitMQ.User, cfg.RabbitMQ.Password, cfg.RabbitMQ.Address, cfg.RabbitMQ.Port)
	rabbitConnection, err := rabbitmq.NewConnection(rabbitURL)
	if err != nil {
		logger.Error("init rabbitmq connection failed", slog.String("error", err.Error()))
		return
	}
	defer func() {
		if err := rabbitConnection.Close(); err != nil {
			logger.Error("close rabbitmq connection failed", slog.String("error", err.Error()))
		}
	}()
	logger.Info("init rabbitmq connection success")

	ch, err := rabbitConnection.Channel()
	if err != nil {
		logger.Error("init rabbitmq channel failed", slog.String("error", err.Error()))
		return
	}
	defer func() {
		if err := ch.Close(); err != nil {
			logger.Error("close rabbitmq channel failed", slog.String("error", err.Error()))
		}
	}()
	logger.Info("init rabbitmq channel success")

	postgresURL := fmt.Sprintf(postgres.URL, cfg.Postgres.User, cfg.Postgres.Password, cfg.Postgres.Address, cfg.Postgres.Port, cfg.Postgres.DB)
	postgresDatabase, err := postgres.NewDatabase(ctx, postgresURL)
	if err != nil {
		logger.Error("init postgres failed", slog.String("error", err.Error()))
		return
	}
	defer postgresDatabase.Close()
	logger.Info("init postgres success")

	a, err := consumer.NewConsumerApp(ctx, cfg, ch, logger, postgresDatabase)
	if err != nil {
		logger.Error("init app failed", slog.String("error", err.Error()))
		return
	}
	logger.Info("init app success")

	logger.Info("start app")
	a.Start(ctx)

	<-ctx.Done()
	logger.Info("stop app")
}
