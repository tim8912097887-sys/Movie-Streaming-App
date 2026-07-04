package api

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type ApiConfig struct {
	Addr string
	Logger *slog.Logger
}

type Api struct {
	Config ApiConfig
}

func (a *Api) Mount() http.Handler {
	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	return router
}

func (a *Api) Run(ctx context.Context, h http.Handler, shutdownTimeout time.Duration) error {
	server := &http.Server{
		Addr:    a.Config.Addr,
		Handler: h,
		ReadTimeout:       5 * time.Second,
        ReadHeaderTimeout: 2 * time.Second,
        WriteTimeout:      10 * time.Second,
        IdleTimeout:       120 * time.Second,
	}

	// Channel to notify when the server is initialized failure
	serverErrorCh := make(chan error, 1)
	// Start the server with goroutine
	go func() {
		a.Config.Logger.Info("starting server",slog.String("address", a.Config.Addr))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.Config.Logger.Error("failed to start server",slog.Any("error", err))
			serverErrorCh <- err
		}
	}()

	select {
		case <-ctx.Done():
			a.Config.Logger.Info("shutting down the server",slog.String("reason", ctx.Err().Error()))
		case err := <-serverErrorCh:
			return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		a.Config.Logger.Error("failed to shut down the server",slog.Any("error", err))
		if closeErr := server.Close(); closeErr != nil {
			a.Config.Logger.Error("failed to close the server",slog.Any("error", err))
			return errors.Join(err,closeErr)
		}
		return err
	}

	a.Config.Logger.Info("server shut down gracefully")
	return nil

}