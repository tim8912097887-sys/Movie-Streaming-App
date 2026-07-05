package movies

type Genres struct {
	GenreID int    `json:"genre_id" bson:"genre_id" validate:"required,numeric"`
	Name    string `json:"name" bson:"name" validate:"required"`
}