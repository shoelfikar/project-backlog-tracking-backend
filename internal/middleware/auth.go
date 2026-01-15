package middleware

import (
	"github.com/gin-gonic/gin"

	"sprint-backlog/internal/utils"
)

// AuthMiddleware validates JWT token and sets user info in context
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.RespondUnauthorized(c, "Authorization header is required")
			c.Abort()
			return
		}

		// Extract token from header
		tokenString := utils.ExtractTokenFromHeader(authHeader)
		if tokenString == "" {
			utils.RespondUnauthorized(c, "Invalid authorization header format")
			c.Abort()
			return
		}

		// Validate token
		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			if err == utils.ErrExpiredToken {
				utils.RespondUnauthorized(c, "Token has expired")
			} else {
				utils.RespondUnauthorized(c, "Invalid token")
			}
			c.Abort()
			return
		}

		// Set user info in context
		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("google_id", claims.GoogleID)

		c.Next()
	}
}
