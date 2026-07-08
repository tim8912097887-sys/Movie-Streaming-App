package admin

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/tim8912097887-sys/server/internal/shared"
	"github.com/tim8912097887-sys/server/internal/shared/types"
	"google.golang.org/genai"
)

type AdminServiceConfig struct {
	UserRepository UserRepository
	MovieRepository MovieRepository
	GenaiClient *genai.Client
}

type UserRepository interface {
	CreateUser(ctx context.Context, user types.CreateUserSchema) (types.User, error)
	FindUserByEmail(ctx context.Context, email string) (types.User, error)
	FindUserById(ctx context.Context, id string) (types.User, error)
	UpdateUser(ctx context.Context, user types.UpdateUserSchema) error
}

type MovieRepository interface {
	GetMovies(ctx context.Context, paginationParams types.PaginationParams) ([]types.Movie,error)
    GetMoviesByGenres(ctx context.Context, genres []types.Genres) ([]types.Movie,error)
	GetMovieById(ctx context.Context, id string) (types.Movie, error)
	UpdateMovie(ctx context.Context, movie types.UpdateMovieSchema) error
}

type service struct{
	userRepository UserRepository
	movieRepository MovieRepository
	genaiClient *genai.Client
}


func NewAdminService(adminServiceConfig AdminServiceConfig) *service {
	return &service{
		userRepository: adminServiceConfig.UserRepository,
		movieRepository: adminServiceConfig.MovieRepository,
		genaiClient: adminServiceConfig.GenaiClient,
	}
}

func (s *service) CreateRating(ctx context.Context, userId string,tokenVersion int, movieId string, adminReview string) error {

	// check if user exists and token version matches
	existingUser, err := s.userRepository.FindUserById(ctx, userId)
	if err != nil {
		if err == shared.ErrUserNotFound {
			return shared.ErrUserNotFound
		}
		return err
	}

	if existingUser.TokenVersion != tokenVersion {
		return shared.ErrTokenVersionMismatch
	}
	// Check admin role
	if existingUser.Role != "admin" {
		return shared.ErrForbidden
	}

	// check if movie exists
	existingMovie, err := s.movieRepository.GetMovieById(ctx, movieId)
	if err != nil {
		if err == shared.ErrMovieNotFound {
			return shared.ErrMovieNotFound
		}
		return err
	}

	// Get Rating from gemini

	config := &genai.GenerateContentConfig{
    ResponseMIMEType: "application/json",
    ResponseSchema: &genai.Schema{
			Type: genai.TypeObject,
			Properties: map[string]*genai.Schema{
				"rating": {
					Type: genai.TypeNumber,
				},
			},
			Required: []string{"rating"},
		},
	}

	prompt := fmt.Sprintf(`
		Based on the following review, rate the movie.

		Title: %s

		Review:
		%s

		Return ONLY valid JSON:

		{
		"rating": 8.5
		}
		`, existingMovie.Title, adminReview)

	resp, err := s.genaiClient.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash",
		genai.Text(prompt),
		config,
	)

	if err != nil {
		return err
	}

	var result struct {
		Rating float64 `json:"rating"`
	}

	err = json.Unmarshal([]byte(resp.Text()), &result)
	
	if err != nil {
		return err
	}

	// update movie rating

	updateMovieSchema := types.UpdateMovieSchema{
		ID: existingMovie.ID,
		Rating: result.Rating,
	}
	err = s.movieRepository.UpdateMovie(ctx, updateMovieSchema)
	if err != nil {
		return err
	}

	return nil
}