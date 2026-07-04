package users

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tim8912097887-sys/server/internal/shared/response"
	"github.com/tim8912097887-sys/server/internal/shared/validation"
)

type UserService interface {
    createUser(ctx context.Context, user CreateUserSchema) (UserDTO, error)
    loginUser(ctx context.Context, user LoginUserSchema) (LoginServiceResponse, error)
	getAllUsers(ctx context.Context) ([]UserDTO, error)
}

type UserHandlerConfig struct {
	UserService UserService
	Logger *slog.Logger
}

type handler struct {
	userService UserService
	logger *slog.Logger
}

func NewUserHandler(userHandlerConfig UserHandlerConfig) *handler {
	return &handler{
		userService: userHandlerConfig.UserService,
		logger: userHandlerConfig.Logger,
	}
}

func (h *handler) RegisterRoutes(r *gin.RouterGroup) {
	r.POST("/register",h.CreateUser)
	r.GET("/",h.GetUsers)
	r.POST("/login",h.LoginUser)
}

func (h *handler) CreateUser(c *gin.Context) {

	var createdUser UserDTO
	var user CreateUserSchema
	var err error
	user, err = validation.BindAndValidate[CreateUserSchema](c.Request)
	if err != nil {
		h.logger.Error("Failed to bind and validate request", slog.Any("error", err))
		c.JSON(http.StatusBadRequest,response.NewErrorResponse("VALIDATION_ERROR",err.Error()))
		return
	}

	createdUser, err = h.userService.createUser(c.Request.Context(), user)

	if err != nil {
		h.logger.Error("Failed to create user", slog.Any("error", err))
		// Handle business errors
		if err == ErrUserAlreadyExists {
			c.JSON(http.StatusBadRequest, response.NewErrorResponse("USER_ALREADY_EXISTS", err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("SERVER_ERROR", err.Error()))
		return
	}

	c.JSON(http.StatusCreated, response.NewSuccessResponse(createdUser))
}

func (h *handler) LoginUser(c *gin.Context) {

	var loginResponse LoginHandlerResponse
	var loginServiceResponse LoginServiceResponse
	var user LoginUserSchema
	var err error
	user, err = validation.BindAndValidate[LoginUserSchema](c.Request)
	if err != nil {
		h.logger.Error("Failed to bind and validate request", slog.Any("error", err))
		c.JSON(http.StatusBadRequest,response.NewErrorResponse("VALIDATION_ERROR",err.Error()))
		return
	}

	loginServiceResponse, err = h.userService.loginUser(c.Request.Context(), user)

	if err != nil {
		h.logger.Error("Failed to login user", slog.Any("error", err))
		// Handle business errors
		if err == ErrInvalidCredentials {
			c.JSON(http.StatusBadRequest, response.NewErrorResponse("INVALID_CREDENTIALS", err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("SERVER_ERROR", err.Error()))
		return
	}

	loginResponse.AccessToken = loginServiceResponse.AccessToken
    // Set httpOnly cookie for refresh token
	c.SetCookie("refresh_token",loginServiceResponse.RefreshToken,7*24*60*60,"/api/v1/users","",false,true)
	c.JSON(http.StatusOK, response.NewSuccessResponse(loginResponse))
}

// For debug
func (h *handler) GetUsers(c *gin.Context) {

	users, err := h.userService.getAllUsers(c.Request.Context())
	if err != nil {
		h.logger.Error("Failed to get users", slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("SERVER_ERROR", err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.NewSuccessResponse(users))
}