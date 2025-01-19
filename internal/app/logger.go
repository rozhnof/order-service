package app

import (
	"log/slog"
	"os"

	"github.com/rozhnof/order-service/internal/pkg/config"
)

func NewLogger(cfg config.Logger) (*slog.Logger, error) {
	var level slog.Leveler

	switch cfg.Level {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warning":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	}

	handler := slog.NewJSONHandler(
		os.Stdout,
		&slog.HandlerOptions{Level: level},
	)

	return slog.New(handler), nil
}
