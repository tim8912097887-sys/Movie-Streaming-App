package movies_test

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/tim8912097887-sys/server/internal/auth"
	"github.com/tim8912097887-sys/server/internal/movies"
	"github.com/tim8912097887-sys/server/internal/shared"
	"github.com/tim8912097887-sys/server/internal/shared/middlewares"
	"github.com/tim8912097887-sys/server/internal/shared/response"
	"github.com/tim8912097887-sys/server/internal/shared/types"
)

func setupRouter(t *testing.T,h *movies.Handler,middlware gin.HandlerFunc) *gin.Engine {
	t.Helper()
	r := gin.Default()
	moviesGroup := r.Group("/api/v1/movies")
	h.RegisterRoutes(moviesGroup, middlware)
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

func wireupHandler(t *testing.T, repository *MockMovieRepository, userRepository *MockUserRepository, jwtService *MockJWTService) *movies.Handler {
	t.Helper()
	handlerOpts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}
	logger := slog.New(slog.NewJSONHandler(os.Stderr, handlerOpts))
	slog.SetDefault(logger)
	movieService := movies.NewMovieService(repository, userRepository)
	moviesHandler := movies.NewMovieHandler(movies.MovieHandlerConfig{MovieService: movieService, Logger: logger})
	return moviesHandler
}

func getMovieRequest(t *testing.T, r *gin.Engine,route string) *http.Response {
	t.Helper()
	
	// Construct request
	req := httptest.NewRequest(http.MethodGet,route,nil)

	w := httptest.NewRecorder()
	// Make request
	r.ServeHTTP(w,req)
	return w.Result()
	
}

func getMovieRequestWithHeaders(t *testing.T, r *gin.Engine, route string, headers map[string]string) *http.Response {
	t.Helper()

	req := httptest.NewRequest(http.MethodGet, route, nil)
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w.Result()
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
	if expectedTag == "invalid syntax" {
		if !strings.Contains(message, expectedTag) {
			t.Fatalf("expected error message to contain %q, got %q", expectedTag, message)
		}
		return
	}

	if !strings.Contains(message, expectedField) {
		t.Fatalf("expected error message to contain field %q, got %q", expectedField, message)
	}
	if !strings.Contains(message, expectedTag) {
		t.Fatalf("expected error message to contain tag %q, got %q", expectedTag, message)
	}
}

func assertErrorResponse(t *testing.T, resp *http.Response, expectedStatus int, expectedCode string) response.ErrorResponse {
	t.Helper()

	if resp.StatusCode != expectedStatus {
		t.Fatalf("expected status code %d, got %d", expectedStatus, resp.StatusCode)
	}

	errorResponse := decodeResponse[response.ErrorResponse](t, resp)
	if errorResponse.Error.Code != expectedCode {
		t.Fatalf("expected error code %q, got %q", expectedCode, errorResponse.Error.Code)
	}
	return errorResponse
}

func TestGetMovieValidation(t *testing.T) {
	movieRepository := InitMockMovieRepository()
	userRepository := InitMockUserRepository()
	jwtService := InitMockJWTService()
	movieHandler := wireupHandler(t, movieRepository, userRepository, jwtService)

    tests := []struct {
		name string
		route string
		expectedField string
		expectedTag  string
	} {
		{
			name: "limit should be greater than 0",
			route: "/api/v1/movies?limit=0",
			expectedField: "Limit",
			expectedTag: "gt",
		},
		{
			name: "limit should be less than 100",
			route: "/api/v1/movies?limit=101",
			expectedField: "Limit",
			expectedTag: "lte",
		},
		{
			name: "page should be greater than or equal to 0",
			route: "/api/v1/movies?page=-1",
			expectedField: "Page",
			expectedTag: "min",
		},
		{
			name: "limit should be numeric",
			route: "/api/v1/movies?limit=abc",
			expectedField: "Limit",
			expectedTag: "invalid syntax",
		},
		{
			name: "page should be numeric",
			route: "/api/v1/movies?page=abc",
			expectedField: "Page",
			expectedTag: "invalid syntax",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := setupRouter(t, movieHandler, InitMockAccessMiddleware(func(c *gin.Context) {}))
			resp := getMovieRequest(t, r, tt.route)
			assertValidationError(t, resp, tt.expectedField, tt.expectedTag)
		})
	}
}

func TestGetMoviesSuccess(t *testing.T) {
	movieRepository := InitMockMovieRepository()
	userRepository := InitMockUserRepository()
	jwtService := InitMockJWTService()
	movieHandler := wireupHandler(t, movieRepository, userRepository, jwtService)

	expectedMovie := types.Movie{
		Title:       "Inception",
		Description: "A mind-bending thriller",
		YoutubeID:   "abcd1234",
		Rating:      8.8,
		PosterURL:   "http://example.com/poster.jpg",
	}

	movieRepository.GetMoviesFunc = func(ctx context.Context, paginationParams types.PaginationParams) ([]types.Movie, int, error) {
		return []types.Movie{expectedMovie}, 1, nil
	}

	r := setupRouter(t, movieHandler, InitMockAccessMiddleware(func(c *gin.Context) {}))
	resp := getMovieRequest(t, r, "/api/v1/movies")

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	successResponse := decodeResponse[response.SuccessResponse](t, resp)
	if successResponse.State != "success" {
		t.Fatalf("expected state %q, got %q", "success", successResponse.State)
	}
	if successResponse.Error != nil {
		t.Fatalf("expected no error, got %v", successResponse.Error)
	}

	paginationData, ok := successResponse.Data.(map[string]any)
	if !ok {
		t.Fatalf("expected response data to be an object, got %T", successResponse.Data)
	}

	moviesData, ok := paginationData["movies"].([]any)
	if !ok {
		t.Fatalf("expected movies field to be a slice, got %T", paginationData["movies"])
	}
	if len(moviesData) != 1 {
		t.Fatalf("expected 1 movie in response data, got %d", len(moviesData))
	}

	movieData, ok := moviesData[0].(map[string]any)
	if !ok {
		t.Fatalf("expected movie item to be an object, got %T", moviesData[0])
	}
	if movieData["title"] != expectedMovie.Title {
		t.Fatalf("expected title %q, got %q", expectedMovie.Title, movieData["title"])
	}
	if movieData["description"] != expectedMovie.Description {
		t.Fatalf("expected description %q, got %q", expectedMovie.Description, movieData["description"])
	}
	if movieData["youtube_id"] != expectedMovie.YoutubeID {
		t.Fatalf("expected youtube_id %q, got %q", expectedMovie.YoutubeID, movieData["youtube_id"])
	}
	if movieData["poster_url"] != expectedMovie.PosterURL {
		t.Fatalf("expected poster_url %q, got %q", expectedMovie.PosterURL, movieData["poster_url"])
	}
	if movieData["rating"] != expectedMovie.Rating {
		t.Fatalf("expected rating %v, got %v", expectedMovie.Rating, movieData["rating"])
	}
}

func TestGetMoviesInternalServerError(t *testing.T) {
	movieRepository := InitMockMovieRepository()
	userRepository := InitMockUserRepository()
	jwtService := InitMockJWTService()
	movieHandler := wireupHandler(t, movieRepository, userRepository, jwtService)

	movieRepository.GetMoviesFunc = func(ctx context.Context, paginationParams types.PaginationParams) ([]types.Movie, int, error) {
		return nil, 0, errors.New("database unavailable")
	}

	r := setupRouter(t, movieHandler, InitMockAccessMiddleware(func(c *gin.Context) {}))
	resp := getMovieRequest(t, r, "/api/v1/movies")

	if resp.StatusCode != http.StatusInternalServerError {
		t.Fatalf("expected status code %d, got %d", http.StatusInternalServerError, resp.StatusCode)
	}

	errorResponse := decodeResponse[response.ErrorResponse](t, resp)
	if errorResponse.Error.Code != "SERVER_ERROR" {
		t.Fatalf("expected error code %q, got %q", "SERVER_ERROR", errorResponse.Error.Code)
	}
	if !strings.Contains(errorResponse.Error.Message, "Internal Server Error") {
		t.Fatalf("expected error message to include %q, got %q", "Internal Server Error", errorResponse.Error.Message)
	}
}

func TestGetUserMovieAuthError(t *testing.T) {
	movieRepository := InitMockMovieRepository()
	userRepository := InitMockUserRepository()
	jwtService := InitMockJWTService()
	movieHandler := wireupHandler(t, movieRepository, userRepository, jwtService)

	secret := "test-secret"
	validToken, err := auth.NewJWTService().GenerateToken(auth.JWTGeneratePayload{
		Subject:      "user-123",
		Secret:       secret,
		Duration:     time.Hour,
		TokenVersion: 1,
	})
	if err != nil {
		t.Fatal(err)
	}

	jwtService.ValidateTokenFunc = func(payload auth.JWTValidatePayload) (auth.CustomClaims, error) {
		if payload.Token == validToken && payload.Secret == secret {
			return auth.CustomClaims{ 
				TokenVersion: 1 ,
			    StandardClaims: jwt.StandardClaims{
					Subject: "user-123",
				},	
			}, nil
		}
		return auth.CustomClaims{}, errors.New("invalid token")
	}

	tests := []struct {
		name           string
		header         string
		expectedStatus int
		expectedCode   string
		expectedMsg    string
	}{
		{
			name:           "missing authorization header",
			header:         "",
			expectedStatus: http.StatusUnauthorized,
			expectedCode:   "UNAUTHORIZED",
			expectedMsg:    "Unauthorized",
		},
		{
			name:           "missing bearer prefix",
			header:         "Token " + validToken,
			expectedStatus: http.StatusUnauthorized,
			expectedCode:   "UNAUTHORIZED",
			expectedMsg:    "Unauthorized",
		},
		{
			name:           "empty bearer token",
			header:         "Bearer ",
			expectedStatus: http.StatusUnauthorized,
			expectedCode:   "UNAUTHORIZED",
			expectedMsg:    "Unauthorized",
		},
		{
			name:           "invalid bearer token",
			header:         "Bearer invalid.token.value",
			expectedStatus: http.StatusUnauthorized,
			expectedCode:   "UNAUTHORIZED",
			expectedMsg:    "Unauthorized",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := setupRouter(t, movieHandler, middlewares.AccessTokenMiddleware(jwtService, secret))
			resp := getMovieRequestWithHeaders(t, r, "/api/v1/movies/user", map[string]string{
				"Authorization": tt.header,
			})
			errorCtx := assertErrorResponse(t, resp, tt.expectedStatus, tt.expectedCode)
        
			if !strings.Contains(errorCtx.Error.Message, tt.expectedMsg) {
				t.Fatalf("expected error message to contain %q, got %q", tt.expectedMsg, errorCtx.Error.Message)
			}
		})
	}
}

func TestGetUserMovieSuccess(t *testing.T) {
	movieRepository := InitMockMovieRepository()
	userRepository := InitMockUserRepository()
	jwtService := InitMockJWTService()
	movieHandler := wireupHandler(t, movieRepository, userRepository, jwtService)

	secret := "test-secret"
	validToken, err := auth.NewJWTService().GenerateToken(auth.JWTGeneratePayload{
		Subject:      "user-123",
		Secret:       secret,
		Duration:     time.Hour,
		TokenVersion: 1,
	})
	if err != nil {
		t.Fatal(err)
	}

	movieRepository.GetMoviesByGenresFunc = func(ctx context.Context, genres []types.Genres) ([]types.Movie, error) {
		return []types.Movie{{
			Title:       "Inception",
			Description: "A mind-bending thriller",
			YoutubeID:   "abcd1234",
			Rating:      8.8,
			PosterURL:   "http://example.com/poster.jpg",
		}}, nil
	}

	userRepository.FindUserByIdFunc = func(ctx context.Context, id string) (types.User, error) {
		return types.User{
			ID:           types.User{}.ID,
			TokenVersion: 1,
			FavoriteGenres: []types.Genres{{GenreID: 1, Name: "Sci-Fi"}},
		}, nil
	}

	jwtService.ValidateTokenFunc = func(payload auth.JWTValidatePayload) (auth.CustomClaims, error) {
		if payload.Token == validToken && payload.Secret == secret {
			return auth.CustomClaims{
				TokenVersion: 1,
				StandardClaims: jwt.StandardClaims{
					Subject: "user-123",
				},
			}, nil
		}
		return auth.CustomClaims{}, errors.New("invalid token")
	}

	r := setupRouter(t, movieHandler, middlewares.AccessTokenMiddleware(jwtService, secret))
	resp := getMovieRequestWithHeaders(t, r, "/api/v1/movies/user", map[string]string{
		"Authorization": "Bearer " + validToken,
	})

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	successResponse := decodeResponse[response.SuccessResponse](t, resp)
	if successResponse.State != "success" {
		t.Fatalf("expected state %q, got %q", "success", successResponse.State)
	}
	if successResponse.Error != nil {
		t.Fatalf("expected no error, got %v", successResponse.Error)
	}

	paginationData, ok := successResponse.Data.(map[string]any)
	if !ok {
		t.Fatalf("expected response data to be an object, got %T", successResponse.Data)
	}

	moviesData, ok := paginationData["movies"].([]any)
	if !ok {
		t.Fatalf("expected movies field to be a slice, got %T", paginationData["movies"])
	}
	if len(moviesData) != 1 {
		t.Fatalf("expected 1 movie in response data, got %d", len(moviesData))
	}

	movieData, ok := moviesData[0].(map[string]any)
	if !ok {
		t.Fatalf("expected movie item to be an object, got %T", moviesData[0])
	}
	if movieData["title"] != "Inception" {
		t.Fatalf("expected title %q, got %q", "Inception", movieData["title"])
	}
}

func TestGetUserMovieInternalServerError(t *testing.T) {
	movieRepository := InitMockMovieRepository()
	userRepository := InitMockUserRepository()
	jwtService := InitMockJWTService()
	movieHandler := wireupHandler(t, movieRepository, userRepository, jwtService)

	secret := "test-secret"
	validToken, err := auth.NewJWTService().GenerateToken(auth.JWTGeneratePayload{
		Subject:      "user-123",
		Secret:       secret,
		Duration:     time.Hour,
		TokenVersion: 1,
	})
	if err != nil {
		t.Fatal(err)
	}

	userRepository.FindUserByIdFunc = func(ctx context.Context, id string) (types.User, error) {
		return types.User{
			ID:           types.User{}.ID,
			TokenVersion: 1,
			FavoriteGenres: []types.Genres{{GenreID: 1, Name: "Sci-Fi"}},
		}, nil
	}

	movieRepository.GetMoviesByGenresFunc = func(ctx context.Context, genres []types.Genres) ([]types.Movie, error) {
		return nil, errors.New("service failure")
	}

	jwtService.ValidateTokenFunc = func(payload auth.JWTValidatePayload) (auth.CustomClaims, error) {
		if payload.Token == validToken && payload.Secret == secret {
			return auth.CustomClaims{
				TokenVersion: 1,
				StandardClaims: jwt.StandardClaims{
					Subject: "user-123",
				},
			}, nil
		}
		return auth.CustomClaims{}, errors.New("invalid token")
	}

	r := setupRouter(t, movieHandler, middlewares.AccessTokenMiddleware(jwtService, secret))
	resp := getMovieRequestWithHeaders(t, r, "/api/v1/movies/user", map[string]string{
		"Authorization": "Bearer " + validToken,
	})

	errorCtx := assertErrorResponse(t, resp, http.StatusInternalServerError, "SERVER_ERROR")
	if !strings.Contains(errorCtx.Error.Message, "Internal Server Error") {
		t.Fatalf("expected error message to contain %q, got %q", "Internal Server Error", errorCtx.Error.Message)
	}
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

type MockMovieRepository struct {
	GetMoviesFunc         func(ctx context.Context, paginationParams types.PaginationParams) ([]types.Movie, int, error)
	GetMoviesByGenresFunc func(ctx context.Context, genres []types.Genres) ([]types.Movie, error)
	GetMovieByIdFunc      func(ctx context.Context, id string) (types.Movie, error)
}

func InitMockMovieRepository() *MockMovieRepository {
	return &MockMovieRepository{
		GetMoviesFunc: func(ctx context.Context, paginationParams types.PaginationParams) ([]types.Movie, int, error) {
			return []types.Movie{}, 1, nil
		},
		GetMoviesByGenresFunc: func(ctx context.Context, genres []types.Genres) ([]types.Movie, error) {
			return []types.Movie{}, nil
		},
		GetMovieByIdFunc: func(ctx context.Context, id string) (types.Movie, error) {
			return types.Movie{}, shared.ErrMovieNotFound
		},
	}
}

func (m *MockMovieRepository) GetMovies(ctx context.Context, paginationParams types.PaginationParams) ([]types.Movie, int, error) {
	return m.GetMoviesFunc(ctx, paginationParams)
}

func (m *MockMovieRepository) GetMoviesByGenres(ctx context.Context, genres []types.Genres) ([]types.Movie, error) {
	return m.GetMoviesByGenresFunc(ctx, genres)
}

func (m *MockMovieRepository) GetMovieById(ctx context.Context, id string) (types.Movie, error) {
	return m.GetMovieByIdFunc(ctx, id)
}

type MockUserRepository struct {
	CreateUserFunc     func(ctx context.Context, user types.CreateUserSchema) (types.User, error)
	FindUserByEmailFunc func(ctx context.Context, email string) (types.User, error)
	FindUserByIdFunc   func(ctx context.Context, id string) (types.User, error)
	UpdateUserFunc     func(ctx context.Context, user types.UpdateUserSchema) error
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

func InitMockAccessMiddleware(middleware gin.HandlerFunc) gin.HandlerFunc {
	return middleware
}