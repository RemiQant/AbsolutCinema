package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"absolutcinema-backend/internal/auth"
)

// AuthMiddleware validates JWT access token from cookies
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get access token from cookie
		accessToken, err := c.Cookie("access_token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Access token not found",
			})
			c.Abort()
			return
		}

		// Validate access token
		claims, err := auth.ValidateAccessToken(accessToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid or expired access token",
			})
			c.Abort()
			return
		}

		// Set user information in context for use in handlers
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_username", claims.Username)
		c.Set("user_role", claims.Role)

		c.Next()
	}
}

// OptionalAuthMiddleware validates JWT if present, but doesn't require it
func OptionalAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get access token from cookie
		accessToken, err := c.Cookie("access_token")
		if err != nil {
			// No token present, continue without authentication
			c.Next()
			return
		}

		// Validate access token
		claims, err := auth.ValidateAccessToken(accessToken)
		if err != nil {
			// Invalid token, continue without authentication
			c.Next()
			return
		}

		// Set user information in context
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_username", claims.Username)
		c.Set("user_role", claims.Role)

		c.Next()
	}
}
