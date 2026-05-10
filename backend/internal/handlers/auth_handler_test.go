package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"ai-e-commerce-app-go/backend/internal/handlers"
	"ai-e-commerce-app-go/backend/internal/models"
	"ai-e-commerce-app-go/backend/internal/repositories"
	"ai-e-commerce-app-go/backend/internal/services"

	"github.com/gin-gonic/gin"
)

type stubAuthService struct {
	register func(ctx context.Context, input models.RegisterUserInput) (models.AuthResponse, error)
	login    func(ctx context.Context, input models.LoginUserInput) (models.AuthResponse, error)
	getUser  func(ctx context.Context, id string) (models.User, error)
}

func (s stubAuthService) Register(ctx context.Context, input models.RegisterUserInput) (models.AuthResponse, error) {
	return s.register(ctx, input)
}

func (s stubAuthService) Login(ctx context.Context, input models.LoginUserInput) (models.AuthResponse, error) {
	return s.login(ctx, input)
}

func (s stubAuthService) GetUser(ctx context.Context, id string) (models.User, error) {
	return s.getUser(ctx, id)
}

func TestAuthRegisterReturnsCreated(t *testing.T) {
	gin.SetMode(gin.TestMode)

	service := stubAuthService{
		register: func(ctx context.Context, input models.RegisterUserInput) (models.AuthResponse, error) {
			return models.AuthResponse{
				User: models.User{
					ID:    "f5e5a9b0-4055-421f-bd6f-7e755a16d1af",
					Name:  input.Name,
					Email: input.Email,
					Role:  models.RoleCustomer,
				},
				Token: "signed-token",
			}, nil
		},
	}

	router := gin.New()
	handler := handlers.NewAuthHandler(service)
	router.POST("/api/v1/auth/register", handler.Register)

	body := bytes.NewBufferString(`{"name":"Learner","email":"learner@example.com","password":"password123"}`)
	request := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", body)
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, response.Code)
	}

	var payload struct {
		Data models.AuthResponse `json:"data"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &payload); err != nil {
		t.Fatalf("expected valid JSON response: %v", err)
	}

	if payload.Data.Token != "signed-token" {
		t.Fatalf("expected token, got %q", payload.Data.Token)
	}
}

func TestAuthRegisterReturnsConflict(t *testing.T) {
	gin.SetMode(gin.TestMode)

	service := stubAuthService{
		register: func(ctx context.Context, input models.RegisterUserInput) (models.AuthResponse, error) {
			return models.AuthResponse{}, repositories.ErrConflict
		},
	}

	router := gin.New()
	handler := handlers.NewAuthHandler(service)
	router.POST("/api/v1/auth/register", handler.Register)

	body := bytes.NewBufferString(`{"name":"Learner","email":"learner@example.com","password":"password123"}`)
	request := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", body)
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusConflict {
		t.Fatalf("expected status %d, got %d", http.StatusConflict, response.Code)
	}
}

func TestAuthLoginReturnsUnauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)

	service := stubAuthService{
		login: func(ctx context.Context, input models.LoginUserInput) (models.AuthResponse, error) {
			return models.AuthResponse{}, services.ErrInvalidCredentials
		},
	}

	router := gin.New()
	handler := handlers.NewAuthHandler(service)
	router.POST("/api/v1/auth/login", handler.Login)

	body := bytes.NewBufferString(`{"email":"learner@example.com","password":"wrong"}`)
	request := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", body)
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, response.Code)
	}
}

func TestAuthMeReturnsCurrentUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	service := stubAuthService{
		getUser: func(ctx context.Context, id string) (models.User, error) {
			return models.User{
				ID:    id,
				Name:  "Learner",
				Email: "learner@example.com",
				Role:  models.RoleCustomer,
			}, nil
		},
	}

	router := gin.New()
	handler := handlers.NewAuthHandler(service)
	router.GET("/api/v1/me", func(c *gin.Context) {
		c.Set("user_id", "f5e5a9b0-4055-421f-bd6f-7e755a16d1af")
		handler.Me(c)
	})

	request := httptest.NewRequest(http.MethodGet, "/api/v1/me", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, response.Code)
	}
}

func TestAuthMeReturnsUnauthorizedWithoutUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	service := stubAuthService{
		getUser: func(ctx context.Context, id string) (models.User, error) {
			return models.User{}, errors.New("should not be called")
		},
	}

	router := gin.New()
	handler := handlers.NewAuthHandler(service)
	router.GET("/api/v1/me", handler.Me)

	request := httptest.NewRequest(http.MethodGet, "/api/v1/me", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, response.Code)
	}
}
