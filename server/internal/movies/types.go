package movies

import (
	"time"

	"github.com/tim8912097887-sys/server/internal/shared/types"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type PaginationParams struct {
    // If not provided, defaults to 0. Must be 0 or greater.
    Offset int `form:"offset,default=0" binding:"min=0"`
    
    // If not provided, defaults to 10. Must be between 1 and 100.
    Limit  int `form:"limit,default=10" binding:"gt=0,lte=100"`
}

type Movie struct {
	ID          bson.ObjectID `json:"id" bson:"_id,omitempty"`
	Title       string        `json:"title" bson:"title"`
	Description string        `json:"description" bson:"description"`
	YoutubeID   string        `json:"youtube_id" bson:"youtube_id"`
	Rating      float64       `json:"rating" bson:"rating"`
	PosterURL   string        `json:"poster_url" bson:"poster_url"`
	Genres      []types.Genres      `json:"genres" bson:"genres"`
	CreatedAt   time.Time         `json:"created_at" bson:"created_at"`
}

type MovieDTO struct {
	ID          bson.ObjectID `json:"id" bson:"_id,omitempty"`
	Title       string        `json:"title" bson:"title"`
	Description string        `json:"description" bson:"description"`
	YoutubeID   string        `json:"youtube_id" bson:"youtube_id"`
	Rating      float64       `json:"rating" bson:"rating"`
	PosterURL   string        `json:"poster_url" bson:"poster_url"`
}
