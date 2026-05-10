package routes

import (
	"context"
	"net/http"
	"strings"

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
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "authorization header required",
			})
			return
		}

		tokenString, ok := strings.CutPrefix(header, "Bearer ")
		if !ok || strings.TrimSpace(tokenString) == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "bearer token required",
			})
			return
		}

		claims, err := tokenParser.ParseToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid or expired token",
			})
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_role", claims.Role)
		c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), userIDContextKey{}, claims.UserID))

		c.Next()
	}
}

type userIDContextKey struct{}
