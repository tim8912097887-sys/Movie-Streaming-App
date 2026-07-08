package movies

import (
	"context"

	"github.com/tim8912097887-sys/server/internal/shared"
	"github.com/tim8912097887-sys/server/internal/shared/types"
)

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
}

type service struct {
	repository MovieRepository
	userRepository UserRepository
}

func NewMovieService(repository MovieRepository, userRepository UserRepository) *service {
	return &service{
		repository: repository,
		userRepository: userRepository,
	}
}

func (s *service) GetMovies(ctx context.Context, paginationParams types.PaginationParams) ([]MovieDTO,error) {
	movies, err := s.repository.GetMovies(ctx, paginationParams)
	if err != nil {
	   return []MovieDTO{}, nil	
	}

	var moviesDTO = make([]MovieDTO,len(movies))

	for i, movie := range movies {
		moviesDTO[i] = s.DataToDTO(movie)
	}

	return moviesDTO, err
}

func (s *service) GetUserMovie(ctx context.Context, userId string,tokenVersion int) ([]MovieDTO,error) {

	// check if user exists and token version matches
	existingUser, err := s.userRepository.FindUserById(ctx, userId)
	if err != nil {
		if err == shared.ErrUserNotFound {
			return []MovieDTO{}, shared.ErrUserNotFound
		}
		return []MovieDTO{}, err
	}

	if existingUser.TokenVersion != tokenVersion {
		return []MovieDTO{}, shared.ErrTokenVersionMismatch
	}

	movies, err := s.repository.GetMoviesByGenres(ctx,existingUser.FavoriteGenres)

	if err != nil {
		return []MovieDTO{}, err
	}

	var moviesDTO = make([]MovieDTO,len(movies))
	for i, movie := range movies {
		moviesDTO[i] = s.DataToDTO(movie)
	}
	return moviesDTO, nil
}

func (s *service) DataToDTO(data types.Movie) MovieDTO {
	return MovieDTO{
		ID:          data.ID,
		Title:       data.Title,
		Description: data.Description,
		YoutubeID:   data.YoutubeID,
		Rating:      data.Rating,
		PosterURL:   data.PosterURL,
	}
}