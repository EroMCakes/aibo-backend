package middlewares

import (
	"net/http"

	"aibo/internal/database"
	"aibo/internal/types"

	"github.com/gin-gonic/gin"
)

// PremiumMiddleware is a middleware that checks if a user is a premium user.
// If the user is not a premium user, it returns a 403 status with a JSON response containing the error message "This feature requires a premium subscription".
// If the user is not found, it returns a 404 status with a JSON response containing the error message "User not found".
// If the user is a premium user, it calls the next handler in the chain.
func PremiumMiddleware(repo *database.AiboRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString("aibo_id")

		user, err := repo.GetAiboByID(userID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, types.ErrorResponse{Error: "User not found"})
			return
		}

		if !user.IsPremium {
			c.AbortWithStatusJSON(http.StatusForbidden, types.ErrorResponse{Error: "This feature requires a premium subscription"})
			return
		}

		c.Next()
	}
}
