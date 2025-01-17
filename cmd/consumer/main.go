package main

import (
	"context"
	"log/slog"
	"os/signal"
	"syscall"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rozhnof/order-service/internal/app"
	"github.com/rozhnof/order-service/internal/app/consumer"
	"github.com/rozhnof/order-service/internal/pkg/postgres"
)

const (
	rabbitURL   = "amqp://user:password@localhost:5672/"
	postgresURL = "postgres://user:password@localhost:5432/order_db?sslmode=disable"
)

func main() {
	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	logger := app.NewLogger()

	conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		logger.Error("init rabbitmq connection failed", slog.String("error", err.Error()))
		return
	}
	defer func() {
		if err := conn.Close(); err != nil {
			logger.Error("close rabbitmq connection failed", slog.String("error", err.Error()))
		}
	}()
	logger.Info("init rabbitmq connection success")

	ch, err := conn.Channel()
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

	postgresDatabase, err := postgres.NewDatabase(ctx, postgresURL)
	if err != nil {
		logger.Error("init postgres failed", slog.String("error", err.Error()))
		return
	}
	defer postgresDatabase.Close()
	logger.Info("init postgres success")

	a, err := consumer.NewConsumerApp(ctx, ch, logger, postgresDatabase)
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
