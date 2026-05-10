package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"ai-e-commerce-app-go/backend/internal/config"
	"ai-e-commerce-app-go/backend/internal/handlers"

	"github.com/gin-gonic/gin"
)

func TestHealthReturnsOK(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	handler := handlers.NewHealthHandler(config.Config{AppName: "ai-e-commerce-api"}, nil)
	router.GET("/health", handler.Health)

	request := httptest.NewRequest(http.MethodGet, "/health", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, response.Code)
	}

	var body map[string]string
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatalf("expected valid JSON response: %v", err)
	}

	if body["status"] != "ok" {
		t.Fatalf("expected status ok, got %q", body["status"])
	}

	if body["service"] != "ai-e-commerce-api" {
		t.Fatalf("expected service ai-e-commerce-api, got %q", body["service"])
	}
}
