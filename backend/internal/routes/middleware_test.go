package routes

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestRequestIDGeneratesAndReturnsHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(requestID())
	router.GET("/health", func(c *gin.Context) {
		if c.GetString("request_id") == "" {
			t.Fatal("expected request id in context")
		}
		c.Status(http.StatusNoContent)
	})

	request := httptest.NewRequest(http.MethodGet, "/health", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Header().Get("X-Request-ID") == "" {
		t.Fatal("expected response X-Request-ID header")
	}
}

func TestRequestIDHonorsIncomingHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(requestID())
	router.GET("/health", func(c *gin.Context) {
		if c.GetString("request_id") != "client-request-id" {
			t.Fatalf("expected client request id, got %q", c.GetString("request_id"))
		}
		c.Status(http.StatusNoContent)
	})

	request := httptest.NewRequest(http.MethodGet, "/health", nil)
	request.Header.Set("X-Request-ID", "client-request-id")
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Header().Get("X-Request-ID") != "client-request-id" {
		t.Fatalf("expected client request id header, got %q", response.Header().Get("X-Request-ID"))
	}
}
