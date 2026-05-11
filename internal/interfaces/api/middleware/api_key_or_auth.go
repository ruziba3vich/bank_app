package middleware

import (
	"crypto/subtle"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/prodonik/bank_app/internal/infrastructure/auth"
	"github.com/prodonik/bank_app/internal/interfaces/api/dto"
)

// APIKeyOrAuth accepts either a valid X-Api-Key header (matching apiKey) or a
// Bearer JWT. Used for endpoints the browser-extension token-refresher can hit
// without holding a bank-app JWT.
func APIKeyOrAuth(jwtService *auth.JWTService, apiKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if apiKey != "" {
			if provided := c.GetHeader("X-Api-Key"); provided != "" {
				if subtle.ConstantTimeCompare([]byte(provided), []byte(apiKey)) == 1 {
					c.Next()
					return
				}
				c.AbortWithStatusJSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "invalid api key"})
				return
			}
		}

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "authorization header or api key required"})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "invalid authorization header format"})
			return
		}

		claims, err := jwtService.ValidateAccessToken(parts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "invalid or expired token"})
			return
		}

		c.Set("user_id", claims.UserID)
		c.Next()
	}
}
