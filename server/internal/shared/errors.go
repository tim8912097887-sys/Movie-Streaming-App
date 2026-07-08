package shared

import "errors"

var (
	ErrUserAlreadyExists = errors.New("user already exists")
    ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotFound = errors.New("user not found")
	ErrTokenVersionMismatch = errors.New("token version mismatch")
    ErrForbidden = errors.New("forbidden")
	ErrMovieNotFound = errors.New("movie not found")
)