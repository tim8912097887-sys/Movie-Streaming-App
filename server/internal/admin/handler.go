package admin

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tim8912097887-sys/server/internal/configs"
	"github.com/tim8912097887-sys/server/internal/shared"
	"github.com/tim8912097887-sys/server/internal/shared/response"
	"github.com/tim8912097887-sys/server/internal/shared/validation"
)

type AdminService interface {
	CreateRating(ctx context.Context, userId string,tokenVersion int, movieId string, adminReview string) error
}

type Handler struct{
	logger *slog.Logger
	envConfigs configs.Configs
	adminService AdminService
}

type AdminHandlerConfig struct {
	Logger *slog.Logger
	EnvConfigs configs.Configs
	AdminService AdminService
}

func NewAdminHandler(adminHandlerConfig AdminHandlerConfig) *Handler {
	return &Handler{
		logger: adminHandlerConfig.Logger,
		envConfigs: adminHandlerConfig.EnvConfigs,
		adminService: adminHandlerConfig.AdminService,
	}
}
func (h *Handler) RegisterRoutes(r *gin.RouterGroup,accessTokenMiddleware gin.HandlerFunc) {
	r.POST("/movies/:movie_id/ratings", accessTokenMiddleware,h.CreateRating)
}

func (h *Handler) CreateRating(c *gin.Context) {

	userId,exist := c.Get("user_id")
	if !exist {
		c.JSON(http.StatusUnauthorized, response.NewErrorResponse("UNAUTHORIZED", "Unauthorized"))
		return
	}
	userIdStr, ok := userId.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("SERVER_ERROR", "Internal Type Assertion Error"))
		return
	}

	tokenVersion,exist := c.Get("token_version")
	if !exist {
		c.JSON(http.StatusUnauthorized, response.NewErrorResponse("UNAUTHORIZED", "Unauthorized"))
		return
	}
	tokenVersionInt, ok := tokenVersion.(int)
	if !ok {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("SERVER_ERROR", "Internal Type Assertion Error"))
		return
	}

	validatedParams, err := validation.Validate(CreateRatingParams{MovieId: c.Param("movie_id")})
	if err != nil {
		h.logger.Error("Failed to bind and validate request", slog.Any("error", err))
		c.JSON(http.StatusBadRequest,response.NewErrorResponse("VALIDATION_ERROR",err.Error()))
		return
	}
	movieId := validatedParams.MovieId

	var adminReview CreateRatingBody
	adminReview, err = validation.ValidateRequestBody[CreateRatingBody](c.Request)
	if err != nil {
		h.logger.Error("Failed to bind and validate request", slog.Any("error", err))
		c.JSON(http.StatusBadRequest,response.NewErrorResponse("VALIDATION_ERROR",err.Error()))
		return
	}

	err = h.adminService.CreateRating(c.Request.Context(), userIdStr,tokenVersionInt, movieId, adminReview.AdminReview)
    // Handle business errors
	if err == shared.ErrMovieNotFound {
		c.JSON(http.StatusNotFound, response.NewErrorResponse("MOVIE_NOT_FOUND", err.Error()))
		return
	}
	if err == shared.ErrUserNotFound {
		c.JSON(http.StatusNotFound, response.NewErrorResponse("USER_NOT_FOUND", err.Error()))
		return
	}
	if err == shared.ErrTokenVersionMismatch {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("TOKEN_VERSION_MISMATCH", err.Error()))
		return
	}
	if err == shared.ErrForbidden {
		c.JSON(http.StatusForbidden, response.NewErrorResponse("FORBIDDEN", err.Error()))
		return
	}
	if err != nil {
		h.logger.Error("Failed to create rating", slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("SERVER_ERROR", err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.NewSuccessResponse("Success"))
}