package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tim8912097887-sys/server/internal/auth"
	"github.com/tim8912097887-sys/server/internal/shared/response"
)

func AccessTokenMiddleware(jwtService JWTService,secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		bearToken := c.Request.Header.Get("Authorization")
		
		const prefix = "Bearer "
		if bearToken == "" || !strings.HasPrefix(bearToken, prefix) {
			c.JSON(http.StatusUnauthorized, response.NewErrorResponse("UNAUTHORIZED", "Unauthorized"))
			c.Abort()
			return
		}

		accessToken := bearToken[len("Bearer "):]

		if accessToken == "" {
			c.JSON(http.StatusUnauthorized, response.NewErrorResponse("UNAUTHORIZED", "Unauthorized"))
			c.Abort()
			return
		}

		accessTokenPayload := auth.JWTValidatePayload{
			Token:  accessToken,
			Secret: secret,
		}

		var customClaims auth.CustomClaims
		customClaims, err := jwtService.ValidateToken(accessTokenPayload)

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