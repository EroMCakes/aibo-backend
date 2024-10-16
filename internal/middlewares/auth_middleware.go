package middlewares

import (
	"net/http"
	"strings"

	"aibo/internal/types"
	"aibo/internal/utilitaries"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, types.ErrorResponse{Error: "Authorization header is required"})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, types.ErrorResponse{Error: "Authorization header format must be Bearer {token}"})
			return
		}

		claims, err := utilitaries.ValidateJWT(parts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, types.ErrorResponse{Error: "Invalid or expired token"})
			return
		}

		c.Set("user_id", claims.AiboID)
		c.Next()
	}
}
