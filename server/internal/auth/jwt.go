package auth

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

type jWTService struct {
}

func NewJWTService() *jWTService {
	return &jWTService{}
}

func (j *jWTService) GenerateToken(payload JWTGeneratePayload) (string, error) {
	claims := &CustomClaims{
		TokenVersion: payload.TokenVersion,
		StandardClaims: jwt.StandardClaims{
			Subject: payload.Subject,
			ExpiresAt: time.Now().Add(payload.Duration).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(payload.Secret))
	
	if err != nil {
		return "", err
	}
	
	return tokenString,nil
}

func (j *jWTService) ValidateToken(payload JWTValidatePayload) (CustomClaims, error) {
	claims := &CustomClaims{}

	token, err := jwt.ParseWithClaims(payload.Token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(payload.Secret), nil
	})

	if err != nil {
		return CustomClaims{}, err
	}

	if !token.Valid {
		return CustomClaims{}, err
	}

	return *claims, nil
}