package users

import "errors"

var (
	ErrUserAlreadyExists = errors.New("user already exists")
    ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotFound = errors.New("user not found")
	ErrTokenVersionMismatch = errors.New("token version mismatch")
)