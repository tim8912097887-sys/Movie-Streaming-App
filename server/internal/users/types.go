package users

type LoginUserSchema struct {
	Email    string `json:"email" validate:"required,email,max=60"`
	Password string `json:"password" validate:"required,min=8"`
}

type UserDTO struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type GenreDTO struct {
	GenreID int    `json:"genre_id"`
	Name    string `json:"name"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type AccessTokenResponse struct {
	AccessToken string `json:"access_token"`
}