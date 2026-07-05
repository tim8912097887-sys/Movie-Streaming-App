package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tim8912097887-sys/server/internal/auth"
	"github.com/tim8912097887-sys/server/internal/shared/response"
)

type JWTService interface {
	GenerateToken(payload auth.JWTGeneratePayload) (string, error)
	ValidateToken(payload auth.JWTValidatePayload) (auth.CustomClaims, error)
}

func RefreshTokenMiddleware(jWTService JWTService,secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		refreshToken, err := c.Request.Cookie("refresh_token")
		
		if err != nil || refreshToken == nil || refreshToken.Value == "" {
			c.JSON(http.StatusUnauthorized, response.NewErrorResponse("UNAUTHORIZED", "Unauthorized"))
			c.Abort()
			return
		}

		refreshTokenPayload := auth.JWTValidatePayload{
			Token:  refreshToken.Value,
			Secret: secret,
		}

		var customClaims auth.CustomClaims
		customClaims, err = jWTService.ValidateToken(refreshTokenPayload)

		if err != nil {
			c.JSON(http.StatusUnauthorized, response.NewErrorResponse("UNAUTHORIZED", "Unauthorized"))
			c.Abort()
			return
		}

        c.Set("user_id", customClaims.Subject)
		c.Set("token_version", customClaims.TokenVersion)

		c.Next()
	}
}