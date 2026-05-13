package routes

import (
	"context"
	"strings"

	"ai-e-commerce-app-go/backend/internal/handlers"
	"ai-e-commerce-app-go/backend/internal/services"

	"github.com/gin-gonic/gin"
)

type TokenParser interface {
	ParseToken(tokenString string) (services.AuthClaims, error)
}

func requireAuth(tokenParser TokenParser) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			handlers.JSONError(c, 401, "UNAUTHORIZED", "authorization header required")
			c.Abort()
			return
		}

		tokenString, ok := strings.CutPrefix(header, "Bearer ")
		if !ok || strings.TrimSpace(tokenString) == "" {
			handlers.JSONError(c, 401, "UNAUTHORIZED", "bearer token required")
			c.Abort()
			return
		}

		claims, err := tokenParser.ParseToken(tokenString)
		if err != nil {
			handlers.JSONError(c, 401, "UNAUTHORIZED", "invalid or expired token")
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_role", claims.Role)
		c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), userIDContextKey{}, claims.UserID))

		c.Next()
	}
}

func requireRole(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("user_role")
		if !exists || userRole != role {
			handlers.JSONError(c, 403, "FORBIDDEN", "insufficient permissions")
			c.Abort()
			return
		}

		c.Next()
	}
}

type userIDContextKey struct{}
