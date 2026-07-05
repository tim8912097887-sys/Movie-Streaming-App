package users

import (
	"github.com/tim8912097887-sys/server/internal/movies"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type CreateUserSchema struct {
	Name            string     `json:"name" bson:"name" validate:"required,min=3,max=50"`
	Email           string     `json:"email" bson:"email" validate:"required,email,max=60"`
	Password        string     `json:"password" bson:"password" validate:"required,min=8"`
	FavoriteGenres []movies.Genres `json:"favorite_genres" bson:"favorite_genres" validate:"required,dive"`
}

type LoginUserSchema struct {
	Email    string `json:"email" validate:"required,email,max=60"`
	Password string `json:"password" validate:"required,min=8"`
}

type User struct {
	ID              bson.ObjectID        `json:"id" bson:"_id,omitempty"`
	Name            string     `json:"name" bson:"name"`
	Email           string     `json:"email" bson:"email"`
	Password        string     `json:"password" bson:"password"`
	TokenVersion    int        `json:"token_version" bson:"token_version"`
	FavoriteGenres []movies.Genres `json:"favorite_genres" bson:"favorite_genres"`
}

type UserDTO struct {
	Name            string     `json:"name"`
	Email           string     `json:"email"`
}

type UpdateUserSchema struct {
	Id              string  
	TokenVersion    int   
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type AccessTokenResponse struct {
	AccessToken  string `json:"access_token"`
}