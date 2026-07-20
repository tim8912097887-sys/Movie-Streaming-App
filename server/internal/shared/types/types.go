package types

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type User struct {
	ID             bson.ObjectID `json:"id" bson:"_id,omitempty"`
	Name           string        `json:"name" bson:"name"`
	Email          string        `json:"email" bson:"email"`
	Password       string        `json:"password" bson:"password"`
	TokenVersion   int           `json:"token_version" bson:"token_version"`
	Role           string        `json:"role" bson:"role"`
	FavoriteGenres []Genres      `json:"favorite_genres" bson:"favorite_genres"`
}

type Genres struct {
	GenreID int    `json:"genre_id" bson:"genre_id" validate:"required,numeric"`
	Name    string `json:"name" bson:"name" validate:"required"`
}

type CreateUserSchema struct {
	Name            string     `json:"name" bson:"name" validate:"required,min=3,max=50"`
	Email           string     `json:"email" bson:"email" validate:"required,email,max=60"`
	Password        string     `json:"password" bson:"password" validate:"required,min=8"`
	FavoriteGenres []Genres `json:"favorite_genres" bson:"favorite_genres" validate:"required,dive"`
}

type UpdateUserSchema struct {
	Id              string  
	TokenVersion    int   
}

type PaginationParams struct {
    // If not provided, defaults to 0. Must be 0 or greater.
	Page int `form:"page,default=0" binding:"numeric,min=0"`
    
    // If not provided, defaults to 10. Must be between 1 and 100.
    Limit  int `form:"limit,default=10" binding:"numeric,gt=0,lte=100"`
}

type Movie struct {
	ID          bson.ObjectID `json:"id" bson:"_id,omitempty"`
	Title       string        `json:"title" bson:"title"`
	Description string        `json:"description" bson:"description"`
	YoutubeID   string        `json:"youtube_id" bson:"youtube_id"`
	Rating      float64       `json:"rating" bson:"rating"`
	PosterURL   string        `json:"poster_url" bson:"poster_url"`
	Genres      []Genres      `json:"genres" bson:"genres"`
	CreatedAt   time.Time         `json:"created_at" bson:"created_at"`
}

type UpdateMovieSchema struct {
	ID          bson.ObjectID `json:"id" bson:"_id,omitempty"`
	Rating      float64 `json:"rating" bson:"rating"`
}