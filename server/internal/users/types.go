package users

import "github.com/tim8912097887-sys/server/internal/movies"

type CreateUserSchema struct {
	Name            string     `json:"name" validate:"required,min=3,max=50"`
	Email           string     `json:"email" validate:"required,email,max=60"`
	Password        string     `json:"password" validate:"required,min=8"`
	Favorite_Genres []movies.Genres `json:"favorite_genres" validate:"required,dive"`
}

type LoginUserSchema struct {
	Email    string `json:"email" validate:"required,email,max=60"`
	Password string `json:"password" validate:"required,min=8"`
}

type User struct {
	ID              string        `json:"id"`
	Name            string     `json:"name"`
	Email           string     `json:"email"`
	Password        string     `json:"password"`
	TokenVersion    int        `json:"token_version"`
	Favorite_Genres []movies.Genres `json:"favorite_genres"`
}

type UserDTO struct {
	Name            string     `json:"name"`
	Email           string     `json:"email"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type AccessTokenResponse struct {
	AccessToken  string `json:"access_token"`
}