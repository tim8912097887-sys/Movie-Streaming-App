package users_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/tim8912097887-sys/server/internal/auth"
	"github.com/tim8912097887-sys/server/internal/configs"
	"github.com/tim8912097887-sys/server/internal/shared"
	"github.com/tim8912097887-sys/server/internal/shared/response"
	"github.com/tim8912097887-sys/server/internal/shared/types"
	"github.com/tim8912097887-sys/server/internal/users"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func getCreateUserSchema(modifiers ...func(*types.CreateUserSchema)) types.CreateUserSchema {
	baseSchema := types.CreateUserSchema{
		Name:     "John Doe",
		Email:    "k6Vz4@example.com",
		Password: "password123",
		FavoriteGenres:   []types.Genres{
			{
				GenreID: 1,
				Name:    "Action",
			},
			{
				GenreID: 2,
				Name:    "Comedy",
			},
		},
	}

	for _, modifier := range modifiers {
		if modifier == nil {
			continue
		}
		modifier(&baseSchema)
	}

	return baseSchema
}

func getLoginUserSchema(modifiers ...func(*users.LoginUserSchema)) users.LoginUserSchema {
	baseSchema := users.LoginUserSchema{
		Email:    "k6Vz4@example.com",
		Password: "password123",
	}

	for _, modifier := range modifiers {
		if modifier == nil {
			continue
		}
		modifier(&baseSchema)
	}

	return baseSchema
}

func setupRouter(t *testing.T,h *users.Handler,middlware gin.HandlerFunc) *gin.Engine {
	t.Helper()
	r := gin.Default()
	usersGroup := r.Group("/api/v1/users")
	h.RegisterRoutes(usersGroup, middlware)
	return r
}

func decodeResponse[T any](t *testing.T,resp *http.Response) T {
	t.Helper()
	var payload T
	err := json.NewDecoder(resp.Body).Decode(&payload)
	if err != nil {
		t.Fatal(err)
	}
	return payload
}

func createUserRequest(t *testing.T, r *gin.Engine,payload types.CreateUserSchema) *http.Response {
	t.Helper()
	// Serialize payload
    body,err := json.Marshal(payload)
	if err != nil {
		t.Fatal(err)
	}
	// Construct request
	req := httptest.NewRequest(http.MethodPost,"/api/v1/users/register",bytes.NewReader(body))
	req.Header.Set("Content-Type","application/json")

	w := httptest.NewRecorder()
	// Make request
	r.ServeHTTP(w,req)
	return w.Result()
	
}

func loginUserRequest(t *testing.T, r *gin.Engine,payload users.LoginUserSchema) *http.Response {
	t.Helper()
	// Serialize payload
    body,err := json.Marshal(payload)
	if err != nil {
		t.Fatal(err)
	}
	// Construct request
	req := httptest.NewRequest(http.MethodPost,"/api/v1/users/login",bytes.NewReader(body))
	req.Header.Set("Content-Type","application/json")

	w := httptest.NewRecorder()
	// Make request
	r.ServeHTTP(w,req)
	return w.Result()
	
}

func wireupHandler(t *testing.T, repository *MockUserRepository, passwordService *MockPasswordService, jwtService *MockJWTService) *users.Handler {
	t.Helper()
	handlerOpts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}
	logger := slog.New(slog.NewJSONHandler(os.Stderr, handlerOpts))
	slog.SetDefault(logger)
	envConfigs := configs.InitConfigs()
	userService := users.NewUserService(users.UserServiceConfig{
		PasswordService: passwordService,
		JWTService:      jwtService,
		EnvConfigs:      envConfigs,
		Repository:      repository,
	})
	userHandler := users.NewUserHandler(users.UserHandlerConfig{UserService: userService, Logger: logger})
	return userHandler
}

func TestRegisterUserValidation(t *testing.T) {
	mockPasswordService := InitMockPasswordService()
	mockJWTService := InitMockJWTService()
	mockUserRepository := InitMockUserRepository()

	handler := wireupHandler(t, mockUserRepository, mockPasswordService, mockJWTService)

	tests := []struct {
		name         string
		modifier     func(*types.CreateUserSchema)
		expectedField string
		expectedTag  string
	}{
		{
			name: "name is required",
			modifier: func(s *types.CreateUserSchema) {
				s.Name = ""
			},
			expectedField: "Name",
			expectedTag:   "required",
		},
		{
			name: "name must be at least 3 characters",
			modifier: func(s *types.CreateUserSchema) {
				s.Name = "Jo"
			},
			expectedField: "Name",
			expectedTag:   "min",
		},
		{
			name: "name must not exceed 50 characters",
			modifier: func(s *types.CreateUserSchema) {
				s.Name = strings.Repeat("a", 51)
			},
			expectedField: "Name",
			expectedTag:   "max",
		},
		{
			name: "email is required",
			modifier: func(s *types.CreateUserSchema) {
				s.Email = ""
			},
			expectedField: "Email",
			expectedTag:   "required",
		},
		{
			name: "email must be a valid email format",
			modifier: func(s *types.CreateUserSchema) {
				s.Email = "invalid-email"
			},
			expectedField: "Email",
			expectedTag:   "email",
		},
		{
			name: "email must not exceed 60 characters",
			modifier: func(s *types.CreateUserSchema) {
				s.Email = strings.Repeat("a", 50) + "@example.com"
			},
			expectedField: "Email",
			expectedTag:   "max",
		},
		{
			name: "password is required",
			modifier: func(s *types.CreateUserSchema) {
				s.Password = ""
			},
			expectedField: "Password",
			expectedTag:   "required",
		},
		{
			name: "password must be at least 8 characters",
			modifier: func(s *types.CreateUserSchema) {
				s.Password = "short"
			},
			expectedField: "Password",
			expectedTag:   "min",
		},
		{
			name: "favorite genres are required",
			modifier: func(s *types.CreateUserSchema) {
				s.FavoriteGenres = nil
			},
			expectedField: "FavoriteGenres",
			expectedTag:   "required",
		},
		{
			name: "favorite genres items must satisfy nested validation",
			modifier: func(s *types.CreateUserSchema) {
				s.FavoriteGenres = []types.Genres{{}}
			},
			expectedField: "FavoriteGenres",
			expectedTag:   "required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payload := getCreateUserSchema(tt.modifier)
			resp := createUserRequest(t, setupRouter(t, handler, InitMockRefreshMiddleware(func(c *gin.Context) {})), payload)

			assertValidationError(t, resp, tt.expectedField, tt.expectedTag)
		})
	}
}

func TestLoginUserValidation(t *testing.T) {
	mockPasswordService := InitMockPasswordService()
	mockJWTService := InitMockJWTService()
	mockUserRepository := InitMockUserRepository()

	handler := wireupHandler(t, mockUserRepository, mockPasswordService, mockJWTService)

	tests := []struct {
		name         string
		modifier     func(*users.LoginUserSchema)
		expectedField string
		expectedTag  string
	} {
		{
			name: "email is required",
			modifier: func(s *users.LoginUserSchema) {
				s.Email = ""
			},
			expectedField: "Email",
			expectedTag:   "required",
		},
		{
			name: "email must be a valid email format",
			modifier: func(s *users.LoginUserSchema) {
				s.Email = "invalid-email"
			},
			expectedField: "Email",
			expectedTag:   "email",
		},
		{
			name: "password is required",
			modifier: func(s *users.LoginUserSchema) {
				s.Password = ""
			},
			expectedField: "Password",
			expectedTag:   "required",
		},
		{
			name: "password must be at least 8 characters",
			modifier: func(s *users.LoginUserSchema) {
				s.Password = "short"
			},
			expectedField: "Password",
			expectedTag:   "min",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payload := getLoginUserSchema(tt.modifier)
			resp := loginUserRequest(t, setupRouter(t, handler, InitMockRefreshMiddleware(func(c *gin.Context) {})), payload)

			assertValidationError(t, resp, tt.expectedField, tt.expectedTag)
		})
	}
}

func TestRegisterUserBusinessLogic(t *testing.T) {
	t.Run("returns created user on success", func(t *testing.T) {
		repository := InitMockUserRepository()
		repository.FindUserByEmailFunc = func(ctx context.Context, email string) (types.User, error) {
			return types.User{}, shared.ErrUserNotFound
		}
		repository.CreateUserFunc = func(ctx context.Context, user types.CreateUserSchema) (types.User, error) {
			return types.User{ID: bson.NewObjectID(), Name: user.Name, Email: user.Email}, nil
		}

		passwordService := InitMockPasswordService()
		passwordService.HashPasswordFunc = func(password string) (string, error) {
			return "hashed-password", nil
		}

		handler := wireupHandler(t, repository, passwordService, InitMockJWTService())
		resp := createUserRequest(t, setupRouter(t, handler, InitMockRefreshMiddleware(func(c *gin.Context) {})), getCreateUserSchema())

		if resp.StatusCode != http.StatusCreated {
			t.Fatalf("expected status code %d, got %d", http.StatusCreated, resp.StatusCode)
		}

		successResponse := decodeResponse[response.SuccessResponse](t, resp)
		if successResponse.State != "success" {
			t.Fatalf("expected success state, got %s", successResponse.State)
		}

		data, ok := successResponse.Data.(map[string]any)
		if !ok {
			t.Fatalf("expected response data to be a map, got %T", successResponse.Data)
		}
		if data["name"] != "John Doe" {
			t.Fatalf("expected created user name to be %q, got %#v", "John Doe", data["name"])
		}
		if data["email"] != "k6Vz4@example.com" {
			t.Fatalf("expected created user email to be %q, got %#v", "k6Vz4@example.com", data["email"])
		}
	})

	t.Run("returns bad request when user already exists", func(t *testing.T) {
		repository := InitMockUserRepository()
		repository.FindUserByEmailFunc = func(ctx context.Context, email string) (types.User, error) {
			return types.User{Email: email}, nil
		}

		handler := wireupHandler(t, repository, InitMockPasswordService(), InitMockJWTService())
		resp := createUserRequest(t, setupRouter(t, handler, InitMockRefreshMiddleware(func(c *gin.Context) {})), getCreateUserSchema())

		assertErrorResponse(t, resp, http.StatusBadRequest, "USER_ALREADY_EXISTS")
	})

	t.Run("returns internal server error when password hashing fails", func(t *testing.T) {
		repository := InitMockUserRepository()
		repository.FindUserByEmailFunc = func(ctx context.Context, email string) (types.User, error) {
			return types.User{}, shared.ErrUserNotFound
		}

		passwordService := InitMockPasswordService()
		passwordService.HashPasswordFunc = func(password string) (string, error) {
			return "", errors.New("hash failed")
		}

		handler := wireupHandler(t, repository, passwordService, InitMockJWTService())
		resp := createUserRequest(t, setupRouter(t, handler, InitMockRefreshMiddleware(func(c *gin.Context) {})), getCreateUserSchema())

		assertErrorResponse(t, resp, http.StatusInternalServerError, "SERVER_ERROR")
	})

	t.Run("returns internal server error when repository create fails", func(t *testing.T) {
		repository := InitMockUserRepository()
		repository.FindUserByEmailFunc = func(ctx context.Context, email string) (types.User, error) {
			return types.User{}, shared.ErrUserNotFound
		}
		repository.CreateUserFunc = func(ctx context.Context, user types.CreateUserSchema) (types.User, error) {
			return types.User{}, errors.New("create failed")
		}

		handler := wireupHandler(t, repository, InitMockPasswordService(), InitMockJWTService())
		resp := createUserRequest(t, setupRouter(t, handler, InitMockRefreshMiddleware(func(c *gin.Context) {})), getCreateUserSchema())

		assertErrorResponse(t, resp, http.StatusInternalServerError, "SERVER_ERROR")
	})

	t.Run("returns internal server error when repository lookup fails unexpectedly", func(t *testing.T) {
		repository := InitMockUserRepository()
		repository.FindUserByEmailFunc = func(ctx context.Context, email string) (types.User, error) {
			return types.User{}, errors.New("lookup failed")
		}

		handler := wireupHandler(t, repository, InitMockPasswordService(), InitMockJWTService())
		resp := createUserRequest(t, setupRouter(t, handler, InitMockRefreshMiddleware(func(c *gin.Context) {})), getCreateUserSchema())

		assertErrorResponse(t, resp, http.StatusInternalServerError, "SERVER_ERROR")
	})
}

func TestGetGenres(t *testing.T) {
	t.Run("returns genres on success", func(t *testing.T) {
		repository := InitMockUserRepository()
		repository.FindAllGenresFunc = func(ctx context.Context) ([]types.Genres, error) {
			return []types.Genres{{GenreID: 1, Name: "Action"}, {GenreID: 2, Name: "Comedy"}}, nil
		}

		handler := wireupHandler(t, repository, InitMockPasswordService(), InitMockJWTService())
		resp := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/api/v1/users/genres", nil)
		r := setupRouter(t, handler, InitMockRefreshMiddleware(func(c *gin.Context) {}))
		r.ServeHTTP(resp, req)

		if resp.Code != http.StatusOK {
			t.Fatalf("expected status code %d, got %d", http.StatusOK, resp.Code)
		}

		successResponse := decodeResponse[response.SuccessResponse](t, resp.Result())
		if successResponse.State != "success" {
			t.Fatalf("expected success state, got %s", successResponse.State)
		}

		data, ok := successResponse.Data.([]any)
		if !ok {
			t.Fatalf("expected response data to be a slice, got %T", successResponse.Data)
		}
		if len(data) != 2 {
			t.Fatalf("expected 2 genres in response data, got %d", len(data))
		}
	})

	t.Run("returns internal server error when repository lookup fails", func(t *testing.T) {
		repository := InitMockUserRepository()
		repository.FindAllGenresFunc = func(ctx context.Context) ([]types.Genres, error) {
			return nil, errors.New("genres lookup failed")
		}

		handler := wireupHandler(t, repository, InitMockPasswordService(), InitMockJWTService())
		resp := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/api/v1/users/genres", nil)
		r := setupRouter(t, handler, InitMockRefreshMiddleware(func(c *gin.Context) {}))
		r.ServeHTTP(resp, req)

		assertErrorResponse(t, resp.Result(), http.StatusInternalServerError, "SERVER_ERROR")
	})
}

func TestLoginUserBusinessLogic(t *testing.T) {
	t.Run("returns tokens on success", func(t *testing.T) {
		repository := InitMockUserRepository()
		repository.FindUserByEmailFunc = func(ctx context.Context, email string) (types.User, error) {
			return types.User{ID: bson.NewObjectID(), TokenVersion: 2, Password: "stored-hash"}, nil
		}

		passwordService := InitMockPasswordService()
		passwordService.CheckPasswordHashFunc = func(password, hash string) bool {
			return password == "password123" && hash == "stored-hash"
		}

		jwtService := InitMockJWTService()
		jwtService.GenerateTokenFunc = func(payload auth.JWTGeneratePayload) (string, error) {
			return "generated-token", nil
		}

		handler := wireupHandler(t, repository, passwordService, jwtService)
		resp := loginUserRequest(t, setupRouter(t, handler, InitMockRefreshMiddleware(func(c *gin.Context) {})), getLoginUserSchema(nil))

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected status code %d, got %d", http.StatusOK, resp.StatusCode)
		}

		successResponse := decodeResponse[response.SuccessResponse](t, resp)
		if successResponse.State != "success" {
			t.Fatalf("expected success state, got %s", successResponse.State)
		}

		data, ok := successResponse.Data.(map[string]any)
		if !ok {
			t.Fatalf("expected response data to be a map, got %T", successResponse.Data)
		}
		if data["access_token"] != "generated-token" {
			t.Fatalf("expected access token to be %q, got %#v", "generated-token", data["access_token"])
		}

		if len(resp.Cookies()) == 0 {
			t.Fatal("expected refresh token cookie to be set")
		}
	})

	t.Run("returns bad request for invalid credentials when user is not found", func(t *testing.T) {
		repository := InitMockUserRepository()
		repository.FindUserByEmailFunc = func(ctx context.Context, email string) (types.User, error) {
			return types.User{}, shared.ErrUserNotFound
		}

		handler := wireupHandler(t, repository, InitMockPasswordService(), InitMockJWTService())
		resp := loginUserRequest(t, setupRouter(t, handler, InitMockRefreshMiddleware(func(c *gin.Context) {})), getLoginUserSchema())

		assertErrorResponse(t, resp, http.StatusBadRequest, "INVALID_CREDENTIALS")
	})

	t.Run("returns bad request for invalid credentials when password does not match", func(t *testing.T) {
		repository := InitMockUserRepository()
		repository.FindUserByEmailFunc = func(ctx context.Context, email string) (types.User, error) {
			return types.User{ID: bson.NewObjectID(), Password: "stored-hash"}, nil
		}

		passwordService := InitMockPasswordService()
		passwordService.CheckPasswordHashFunc = func(password, hash string) bool {
			return false
		}

		handler := wireupHandler(t, repository, passwordService, InitMockJWTService())
		resp := loginUserRequest(t, setupRouter(t, handler, InitMockRefreshMiddleware(func(c *gin.Context) {})), getLoginUserSchema())

		assertErrorResponse(t, resp, http.StatusBadRequest, "INVALID_CREDENTIALS")
	})

	t.Run("returns internal server error when repository lookup fails unexpectedly", func(t *testing.T) {
		repository := InitMockUserRepository()
		repository.FindUserByEmailFunc = func(ctx context.Context, email string) (types.User, error) {
			return types.User{}, errors.New("lookup failed")
		}

		handler := wireupHandler(t, repository, InitMockPasswordService(), InitMockJWTService())
		resp := loginUserRequest(t, setupRouter(t, handler, InitMockRefreshMiddleware(func(c *gin.Context) {})), getLoginUserSchema())

		assertErrorResponse(t, resp, http.StatusInternalServerError, "SERVER_ERROR")
	})

	t.Run("returns internal server error when jwt generation fails", func(t *testing.T) {
		repository := InitMockUserRepository()
		repository.FindUserByEmailFunc = func(ctx context.Context, email string) (types.User, error) {
			return types.User{ID: bson.NewObjectID(), TokenVersion: 1, Password: "stored-hash"}, nil
		}

		passwordService := InitMockPasswordService()
		passwordService.CheckPasswordHashFunc = func(password, hash string) bool {
			return true
		}

		jwtService := InitMockJWTService()
		jwtService.GenerateTokenFunc = func(payload auth.JWTGeneratePayload) (string, error) {
			return "", errors.New("token failed")
		}

		handler := wireupHandler(t, repository, passwordService, jwtService)
		resp := loginUserRequest(t, setupRouter(t, handler, InitMockRefreshMiddleware(func(c *gin.Context) {})), getLoginUserSchema())

		assertErrorResponse(t, resp, http.StatusInternalServerError, "SERVER_ERROR")
	})
}

func assertValidationError(t *testing.T, resp *http.Response, expectedField, expectedTag string) {
	t.Helper()

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected status code %d, got %d", http.StatusBadRequest, resp.StatusCode)
	}

	errorResponse := decodeResponse[response.ErrorResponse](t, resp)
	if errorResponse.Error.Code != "VALIDATION_ERROR" {
		t.Fatalf("expected error code %s, got %s", "VALIDATION_ERROR", errorResponse.Error.Code)
	}

	message := errorResponse.Error.Message
	if !strings.Contains(message, expectedField) {
		t.Fatalf("expected error message to contain field %q, got %q", expectedField, message)
	}
	if !strings.Contains(message, expectedTag) {
		t.Fatalf("expected error message to contain tag %q, got %q", expectedTag, message)
	}
}

func assertErrorResponse(t *testing.T, resp *http.Response, expectedStatus int, expectedCode string) {
	t.Helper()

	if resp.StatusCode != expectedStatus {
		t.Fatalf("expected status code %d, got %d", expectedStatus, resp.StatusCode)
	}

	errorResponse := decodeResponse[response.ErrorResponse](t, resp)
	if errorResponse.Error.Code != expectedCode {
		t.Fatalf("expected error code %s, got %s", expectedCode, errorResponse.Error.Code)
	}
}

func InitMockRefreshMiddleware(middleware gin.HandlerFunc) gin.HandlerFunc {
	return middleware
}

type MockPasswordService struct{
	HashPasswordFunc func(password string) (string, error)
	CheckPasswordHashFunc func(password, hash string) bool
}

func InitMockPasswordService() *MockPasswordService {
	return &MockPasswordService{
		HashPasswordFunc: func(password string) (string, error) {
			return "hashedPassword", nil
		},
		CheckPasswordHashFunc: func(password, hash string) bool {
			return true
		},
	}
}

func (m *MockPasswordService) HashPassword(password string) (string, error) {
	return m.HashPasswordFunc(password)
}

func (m *MockPasswordService) CheckPasswordHash(password, hash string) bool {
	return m.CheckPasswordHashFunc(password, hash)
}

type MockUserRepository struct {
	CreateUserFunc      func(ctx context.Context, user types.CreateUserSchema) (types.User, error)
	FindUserByEmailFunc func(ctx context.Context, email string) (types.User, error)
	FindUserByIdFunc    func(ctx context.Context, id string) (types.User, error)
	UpdateUserFunc      func(ctx context.Context, user types.UpdateUserSchema) error
	FindAllGenresFunc   func(ctx context.Context) ([]types.Genres, error)
}

func InitMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		CreateUserFunc: func(ctx context.Context, user types.CreateUserSchema) (types.User, error) {
			return types.User{}, nil
		},
		FindUserByEmailFunc: func(ctx context.Context, email string) (types.User, error) {
			return types.User{}, shared.ErrUserNotFound
		},
		FindUserByIdFunc: func(ctx context.Context, id string) (types.User, error) {
			return types.User{}, shared.ErrUserNotFound
		},
		UpdateUserFunc: func(ctx context.Context, user types.UpdateUserSchema) error {
			return nil
		},
		FindAllGenresFunc: func(ctx context.Context) ([]types.Genres, error) {
			return []types.Genres{}, nil
		},
	}
}

func (m *MockUserRepository) CreateUser(ctx context.Context, user types.CreateUserSchema) (types.User, error) {
	return m.CreateUserFunc(ctx, user)
}

func (m *MockUserRepository) FindUserByEmail(ctx context.Context, email string) (types.User, error) {
	return m.FindUserByEmailFunc(ctx, email)
}

func (m *MockUserRepository) FindUserById(ctx context.Context, id string) (types.User, error) {
	return m.FindUserByIdFunc(ctx, id)
}

func (m *MockUserRepository) UpdateUser(ctx context.Context, user types.UpdateUserSchema) error {
	return m.UpdateUserFunc(ctx, user)
}

func (m *MockUserRepository) FindAllGenres(ctx context.Context) ([]types.Genres, error) {
	return m.FindAllGenresFunc(ctx)
}

type MockJWTService struct {
	GenerateTokenFunc func(payload auth.JWTGeneratePayload) (string, error)
	ValidateTokenFunc func(payload auth.JWTValidatePayload) (auth.CustomClaims, error)
}

func InitMockJWTService() *MockJWTService {
	return &MockJWTService{
		GenerateTokenFunc: func(payload auth.JWTGeneratePayload) (string, error) {
			return "mock-token", nil
		},
		ValidateTokenFunc: func(payload auth.JWTValidatePayload) (auth.CustomClaims, error) {
			return auth.CustomClaims{}, nil
		},
	}
}

func (m *MockJWTService) GenerateToken(payload auth.JWTGeneratePayload) (string, error) {
	return m.GenerateTokenFunc(payload)
}

func (m *MockJWTService) ValidateToken(payload auth.JWTValidatePayload) (auth.CustomClaims, error) {
	return m.ValidateTokenFunc(payload)
}