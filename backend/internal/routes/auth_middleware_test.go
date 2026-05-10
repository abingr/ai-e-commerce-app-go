package routes

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestRequireRoleAllowsMatchingRole(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.GET("/admin", func(c *gin.Context) {
		c.Set("user_role", "admin")
	}, requireRole("admin"), func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})

	request := httptest.NewRequest(http.MethodGet, "/admin", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusNoContent {
		t.Fatalf("expected status %d, got %d", http.StatusNoContent, response.Code)
	}
}

func TestRequireRoleRejectsWrongRole(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.GET("/admin", func(c *gin.Context) {
		c.Set("user_role", "customer")
	}, requireRole("admin"), func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})

	request := httptest.NewRequest(http.MethodGet, "/admin", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusForbidden {
		t.Fatalf("expected status %d, got %d", http.StatusForbidden, response.Code)
	}
}
