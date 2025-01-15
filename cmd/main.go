package main

import (
	"context"
	"os/signal"
	"syscall"
)

const (
	EnvConfigPath = "CONFIG_PATH"
)

func main() {
	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	_ = ctx
}
