package admin

type CreateRatingParams struct {
	MovieId string `validate:"required"`
}

type CreateRatingBody struct {
	AdminReview string `json:"admin_review" validate:"required,min=1,max=500"`
}