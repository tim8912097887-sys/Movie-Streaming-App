package movies

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tim8912097887-sys/server/internal/shared"
	"github.com/tim8912097887-sys/server/internal/shared/response"
	"github.com/tim8912097887-sys/server/internal/shared/types"
	"github.com/tim8912097887-sys/server/internal/shared/validation"
)

type MovieService interface {
	GetMovies(ctx context.Context,paginationParams types.PaginationParams) ([]MovieDTO, int,error)
    GetUserMovie(ctx context.Context,userId string,tokenVersion int) ([]MovieDTO,error)
}

type MovieHandlerConfig struct {
	MovieService MovieService
	Logger *slog.Logger
}

type Handler struct{
	movieService MovieService
	logger *slog.Logger
}

func NewMovieHandler(movieHandlerConfig MovieHandlerConfig) *Handler {
	return &Handler{
		movieService: movieHandlerConfig.MovieService,
		logger: movieHandlerConfig.Logger,
	}
}

func (h *Handler) RegisterRoutes(r *gin.RouterGroup,accessMiddleware gin.HandlerFunc) {
	r.GET("",h.GetMovies)
	r.GET("/user",accessMiddleware,h.GetUserMovie)
}

func (h *Handler) GetMovies(c *gin.Context) {
     
	paginationParams, err := validation.ValidateQueryParams[types.PaginationParams](c)

	if err != nil {
		h.logger.Error("Failed to bind and validate request", slog.Any("error", err))
		c.JSON(http.StatusBadRequest,response.NewErrorResponse("VALIDATION_ERROR",err.Error()))
		return
	}

	movies, totalPages, err := h.movieService.GetMovies(c.Request.Context(), paginationParams)

	// Timeout or cancel error
	if errors.Is(err, context.DeadlineExceeded) {
		c.JSON(http.StatusGatewayTimeout,
		response.NewErrorResponse("REQUEST_TIMEOUT", "Request timed out"))
		return
	}

	if errors.Is(err, context.Canceled) {
		h.logger.Info("Request canceled", slog.Any("error", err))
		return 
	}

	if err != nil {
		h.logger.Error("Failed to get movies", slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("SERVER_ERROR", err.Error()))
		return
	}

	paginationResponse := PaginationResponse{
		Movies: movies,
		TotalPage: totalPages,
		CurrentPage: paginationParams.Page,
	}

	c.JSON(http.StatusOK, response.NewSuccessResponse(paginationResponse))
}

func (h *Handler) GetUserMovie(c *gin.Context) {
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

	movies, err := h.movieService.GetUserMovie(c.Request.Context(), userIdStr,tokenVersionInt)

	// Handle business errors
	if err == shared.ErrTokenVersionMismatch {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("TOKEN_VERSION_MISMATCH", err.Error()))
		return
	}
	if err == shared.ErrUserNotFound {
		c.JSON(http.StatusNotFound, response.NewErrorResponse("USER_NOT_FOUND", err.Error()))
		return
	}

	// Timeout or cancel error
	if errors.Is(err, context.DeadlineExceeded) {
		c.JSON(http.StatusGatewayTimeout,
		response.NewErrorResponse("REQUEST_TIMEOUT", "Request timed out"))
		return
	}

	if errors.Is(err, context.Canceled) {
		h.logger.Info("Request canceled", slog.Any("error", err))
		return 
	}

	if err != nil {
		h.logger.Error("Failed to get movies", slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("SERVER_ERROR", err.Error()))
		return
	}

	paginationResponse := PaginationResponse{
		Movies: movies,
		TotalPage: 1,
		CurrentPage: 1,
	}

	c.JSON(http.StatusOK, response.NewSuccessResponse(paginationResponse))
}