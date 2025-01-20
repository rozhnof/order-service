package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/rozhnof/order-service/internal/app"
	"github.com/rozhnof/order-service/internal/app/consumer"
	"github.com/rozhnof/order-service/internal/pkg/config"
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

	a, err := consumer.NewConsumerApp(ctx, cfg, ch, logger)
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
