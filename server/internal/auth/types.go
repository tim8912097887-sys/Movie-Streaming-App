package auth

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

type CustomClaims struct {
	TokenVersion int           `json:"token_version"`
	jwt.StandardClaims
}

type JWTGeneratePayload struct {
	Subject      string        `json:"subject"`
	Secret   string        `json:"secret"`
	Duration time.Duration `json:"duration"`
	TokenVersion int           `json:"token_version"`
}

type JWTValidatePayload struct {
	Token string `json:"token"`
	Secret string `json:"secret"`
}