package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthRequired middleware checks for valid authentication
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header required",
			})
			c.Abort()
			return
		}

		// Check Bearer token format
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid authorization header format",
			})
			c.Abort()
			return
		}

		token := parts[1]
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Token required",
			})
			c.Abort()
			return
		}

		// TODO: Validate JWT token and extract user information
		// For now, we'll use a placeholder user ID
		userID := "user-123" // This should be extracted from the validated JWT

		// Set user ID in context
		c.Set("user_id", userID)
		c.Next()
	}
}

// GetUserIDFromContext extracts user ID from Gin context
func GetUserIDFromContext(c *gin.Context) (string, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return "", false
	}

	if id, ok := userID.(string); ok {
		return id, true
	}

	return "", false
}

// GetUserIDFromStdContext extracts user ID from standard context
func GetUserIDFromStdContext(ctx context.Context) (string, bool) {
	userID := ctx.Value("user_id")
	if userID == nil {
		return "", false
	}

	if id, ok := userID.(string); ok {
		return id, true
	}

	return "", false
}

// SetUserIDInContext sets user ID in standard context
func SetUserIDInContext(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, "user_id", userID)
}