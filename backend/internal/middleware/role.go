package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RequireRole creates middleware that checks if the user has one of the required roles
func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user role from context (set by AuthMiddleware)
		userRole, exists := c.Get("user_role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized - no role found",
			})
			c.Abort()
			return
		}

		// Check if user has one of the required roles
		role := userRole.(string)
		hasRole := false
		for _, requiredRole := range roles {
			if role == requiredRole {
				hasRole = true
				break
			}
		}

		if !hasRole {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Forbidden - insufficient permissions",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireAdmin is a convenience middleware for admin-only routes
func RequireAdmin() gin.HandlerFunc {
	return RequireRole("admin")
}

// RequireCustomer is a convenience middleware for customer routes
func RequireCustomer() gin.HandlerFunc {
	return RequireRole("customer")
}

// RequireAdminOrCustomer allows both admin and customer roles
func RequireAdminOrCustomer() gin.HandlerFunc {
	return RequireRole("admin", "customer")
}
