package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/tim8912097887-sys/server/cmd/api"
	"github.com/tim8912097887-sys/server/internal/configs"
)

func main() {

	handlerOpts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}
	logger := slog.New(slog.NewJSONHandler(os.Stderr, handlerOpts))
	slog.SetDefault(logger)

	envConfigs := configs.InitConfigs()

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()
	apiConfig := api.ApiConfig{Addr: envConfigs.Addr, Logger: logger}
	api := api.Api{Config: apiConfig}
	api.Run(ctx, api.Mount(), 8*time.Second)
}