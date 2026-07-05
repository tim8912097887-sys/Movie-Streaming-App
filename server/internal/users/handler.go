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
    loginUser(ctx context.Context, user LoginUserSchema) (TokenResponse, error)
	logoutUser(ctx context.Context, userId string,tokenVersion int) error
	refreshToken(ctx context.Context, userId string,tokenVersion int) (TokenResponse, error)
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

func (h *handler) RegisterRoutes(r *gin.RouterGroup,refreshTokenMiddleware gin.HandlerFunc) {
	r.POST("/register",h.CreateUser)
	r.GET("/",h.GetUsers)
	r.POST("/login",h.LoginUser)
	r.POST("/logout",refreshTokenMiddleware,h.LogoutUser)
	r.POST("/refresh",refreshTokenMiddleware,h.RefreshToken)
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

	var accessTokenResponse AccessTokenResponse
	var tokenResponse TokenResponse 
	var user LoginUserSchema
	var err error
	user, err = validation.BindAndValidate[LoginUserSchema](c.Request)
	if err != nil {
		h.logger.Error("Failed to bind and validate request", slog.Any("error", err))
		c.JSON(http.StatusBadRequest,response.NewErrorResponse("VALIDATION_ERROR",err.Error()))
		return
	}

	tokenResponse, err = h.userService.loginUser(c.Request.Context(), user)

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

	accessTokenResponse.AccessToken = tokenResponse.AccessToken
    // Set httpOnly cookie for refresh token
	setCookie(c,tokenResponse.RefreshToken,7*24*60*60)
	c.JSON(http.StatusOK, response.NewSuccessResponse(accessTokenResponse))
}

func (h *handler) LogoutUser(c *gin.Context) {
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

	err := h.userService.logoutUser(c.Request.Context(), userIdStr,tokenVersionInt)
	if err != nil {
		h.logger.Error("Failed to logout user", slog.Any("error", err))
		// Handle business errors
		if err == ErrUserNotFound {
			// Clear cookie
		    setCookie(c,"",-1)
			c.JSON(http.StatusBadRequest, response.NewErrorResponse("USER_NOT_FOUND", err.Error()))
			return
		}
		if err == ErrTokenVersionMismatch {
			// Clear cookie
		    setCookie(c,"",-1)
			c.JSON(http.StatusBadRequest, response.NewErrorResponse("TOKEN_VERSION_MISMATCH", err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("SERVER_ERROR", err.Error()))
		return
	}
	setCookie(c,"",-1)
	c.JSON(http.StatusOK, response.NewSuccessResponse("Logout successfully"))
}

func (h *handler) RefreshToken(c *gin.Context) {
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

	var tokenResponse TokenResponse
	var accessTokenResponse AccessTokenResponse

	tokenResponse,err := h.userService.refreshToken(c.Request.Context(), userIdStr,tokenVersionInt)
	if err != nil {
		h.logger.Error("Failed to logout user", slog.Any("error", err))
		
		// Handle business errors
		if err == ErrUserNotFound {
			// Clear cookie
		    setCookie(c,"",-1)
			c.JSON(http.StatusBadRequest, response.NewErrorResponse("USER_NOT_FOUND", err.Error()))
			return
		}
		if err == ErrTokenVersionMismatch {
			// Clear cookie
		    setCookie(c,"",-1)
			c.JSON(http.StatusBadRequest, response.NewErrorResponse("TOKEN_VERSION_MISMATCH", err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("SERVER_ERROR", err.Error()))
		return
	}

	accessTokenResponse.AccessToken = tokenResponse.AccessToken
	// Set httpOnly cookie for refresh token
	setCookie(c,tokenResponse.RefreshToken,7*24*60*60)
	c.JSON(http.StatusOK, response.NewSuccessResponse(accessTokenResponse))
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

func setCookie(c *gin.Context, value string,expiredTime int) {
	c.SetCookie("refresh_token", value, expiredTime, "/api/v1/users", "", false, true)
}