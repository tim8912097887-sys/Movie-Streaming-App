package auth

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

type CustomClaims struct {
	Email  string `json:"email"`
	jwt.StandardClaims
}

type JWTGeneratePayload struct {
	Subject      string        `json:"subject"`
	Secret   string        `json:"secret"`
	Email  string        `json:"email"`
	Duration time.Duration `json:"duration"`
}