package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/rozhnof/order-service/internal/app"
	"github.com/rozhnof/order-service/internal/app/publisher"
	"github.com/rozhnof/order-service/internal/pkg/config"
	"github.com/rozhnof/order-service/internal/pkg/postgres"
	"github.com/rozhnof/order-service/internal/pkg/rabbitmq"
)

const (
	EnvConfigPath = "CONFIG_PATH"
)

func main() {
	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	cfg, err := config.NewConfig[publisher.Config](os.Getenv(EnvConfigPath))
	if err != nil {
		slog.Error("init config failed", slog.String("error", err.Error()))
		return
	}

	logger, err := app.NewLogger(cfg.Logger)
	if err != nil {
		slog.Error("init logger failed", slog.String("error", err.Error()))
		return
	}

	rabbitConnection, err := rabbitmq.NewConnection(cfg.RabbitMQ.ConnectionURL())
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

	postgresDatabase, err := postgres.NewDatabase(ctx, cfg.Postgres.ConnectionURL())
	if err != nil {
		logger.Error("init postgres failed", slog.String("error", err.Error()))
		return
	}
	defer postgresDatabase.Close()
	logger.Info("init postgres success")

	a, err := publisher.NewPublisherApp(ctx, cfg, ch, logger, postgresDatabase)
	if err != nil {
		logger.Error("init app failed", slog.String("error", err.Error()))
		return
	}
	logger.Info("init app success")

	logger.Info("run app")
	if err := a.Run(ctx); err != nil {
		logger.Error("app error", slog.String("error", err.Error()))
		return
	}

	logger.Error("shutdown app")
}
