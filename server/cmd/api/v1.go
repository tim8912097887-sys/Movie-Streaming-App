package api

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tim8912097887-sys/server/internal/auth"
	"github.com/tim8912097887-sys/server/internal/configs"
	"github.com/tim8912097887-sys/server/internal/shared/middlewares"
	"github.com/tim8912097887-sys/server/internal/shared/response"
	"github.com/tim8912097887-sys/server/internal/users"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type ApiConfig struct {
	Logger *slog.Logger
	EnvConfigs configs.Configs
}

type Api struct {
	Config ApiConfig
}

func (a *Api) Mount(dbClient *mongo.Client) http.Handler {
	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, response.NewSuccessResponse(map[string]string{
			"status": "OK",
		}))
	})

	v1Router :=router.Group("/api/v1")

	// Register user routes
	userRouter := v1Router.Group("/users")
	passwordService := auth.NewPasswordService()
	jwtService := auth.NewJWTService()
	refreshTokenMiddleware := middlewares.RefreshTokenMiddleware(jwtService, a.Config.EnvConfigs.RefreshTokenSecret)
	userCollection := dbClient.Database(a.Config.EnvConfigs.DbName).Collection("users")
	userRepository := users.NewUserRepository(userCollection)
	userServiceConfig := users.UserServiceConfig{PasswordService: passwordService, JWTService: jwtService, EnvConfigs: a.Config.EnvConfigs, Repository: userRepository}
	userService := users.NewUserService(userServiceConfig)
	userHandlerConfig := users.UserHandlerConfig{UserService: userService, Logger: a.Config.Logger}
	userHandler := users.NewUserHandler(userHandlerConfig)
	userHandler.RegisterRoutes(userRouter, refreshTokenMiddleware)
	
	return router
}

func (a *Api) Run(ctx context.Context, h http.Handler, shutdownTimeout time.Duration) error {
	server := &http.Server{
		Addr:    a.Config.EnvConfigs.Addr,
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
		a.Config.Logger.Info("starting server",slog.String("address", a.Config.EnvConfigs.Addr))
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