package types

import "go.mongodb.org/mongo-driver/v2/bson"

type User struct {
	ID             bson.ObjectID `json:"id" bson:"_id,omitempty"`
	Name           string        `json:"name" bson:"name"`
	Email          string        `json:"email" bson:"email"`
	Password       string        `json:"password" bson:"password"`
	TokenVersion   int           `json:"token_version" bson:"token_version"`
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