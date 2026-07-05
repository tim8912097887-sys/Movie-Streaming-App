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
	"github.com/tim8912097887-sys/server/internal/db"
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

	// Initialize db
	dbClient, err := db.ConnectDB(ctx, envConfigs.DbUrl)
	if err != nil {
		logger.Error("Failed to connect to db", slog.Any("error", err))
		os.Exit(1)
		return
	}
	logger.Info("Connected to db", slog.String("name", envConfigs.DbName))
	
	defer func() {
        logger.Info("Disconnecting from MongoDB...")
        if err := dbClient.Disconnect(context.Background()); err != nil {
            logger.Error("Failed to cleanly disconnect MongoDB", slog.Any("error", err))
        }
    }()
	
	apiConfig := api.ApiConfig{Logger: logger, EnvConfigs: envConfigs}
	api := api.Api{Config: apiConfig}
	api.Run(ctx, api.Mount(dbClient), 8*time.Second)
}