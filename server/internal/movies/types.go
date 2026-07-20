package movies

import (
	"go.mongodb.org/mongo-driver/v2/bson"
)

type MovieDTO struct {
	ID          bson.ObjectID `json:"id" bson:"_id,omitempty"`
	Title       string        `json:"title" bson:"title"`
	Description string        `json:"description" bson:"description"`
	YoutubeID   string        `json:"youtube_id" bson:"youtube_id"`
	Rating      float64       `json:"rating" bson:"rating"`
	PosterURL   string        `json:"poster_url" bson:"poster_url"`
}

type PaginationResponse struct {
	Movies []MovieDTO `json:"movies" bson:"movies"`
	TotalPage int        `json:"total_page" bson:"total_page"`
	CurrentPage int `json:"current_page" bson:"current_page"`
}
