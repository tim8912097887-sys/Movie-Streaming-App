package movies

type Genres struct {
	GenreID int    `json:"genre_id" validate:"required,numeric"`
	Name    string `json:"name" validate:"required"`
}