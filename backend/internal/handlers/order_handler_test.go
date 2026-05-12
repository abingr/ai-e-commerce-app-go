package handlers_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"ai-e-commerce-app-go/backend/internal/handlers"
	"ai-e-commerce-app-go/backend/internal/models"
	"ai-e-commerce-app-go/backend/internal/repositories"

	"github.com/gin-gonic/gin"
)

type stubOrderService struct {
	createFromCart  func(ctx context.Context, userID string) (models.Order, error)
	listByUser      func(ctx context.Context, userID string) ([]models.Order, error)
	findByIDForUser func(ctx context.Context, userID string, orderID string) (models.Order, error)
}

func (s stubOrderService) CreateFromCart(ctx context.Context, userID string) (models.Order, error) {
	return s.createFromCart(ctx, userID)
}

func (s stubOrderService) ListByUser(ctx context.Context, userID string) ([]models.Order, error) {
	return s.listByUser(ctx, userID)
}

func (s stubOrderService) FindByIDForUser(ctx context.Context, userID string, orderID string) (models.Order, error) {
	return s.findByIDForUser(ctx, userID, orderID)
}

func TestOrderCreateFromCartReturnsCreatedOrder(t *testing.T) {
	gin.SetMode(gin.TestMode)

	service := stubOrderService{
		createFromCart: func(ctx context.Context, userID string) (models.Order, error) {
			if userID != "user-1" {
				t.Fatalf("expected user-1, got %q", userID)
			}

			return models.Order{
				ID:         "8d933930-8022-4df2-91cc-a7a28c0c53ef",
				UserID:     userID,
				Status:     models.OrderStatusConfirmed,
				TotalCents: 25800,
				Items: []models.OrderItem{
					{
						ProductID:      "f5e5a9b0-4055-421f-bd6f-7e755a16d1af",
						ProductName:    "Keyboard",
						Quantity:       2,
						UnitPriceCents: 12900,
						LineTotalCents: 25800,
					},
				},
			}, nil
		},
	}

	router := gin.New()
	handler := handlers.NewOrderHandler(service)
	router.POST("/api/v1/orders", withUser("user-1"), handler.CreateFromCart)

	request := httptest.NewRequest(http.MethodPost, "/api/v1/orders", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, response.Code)
	}

	var payload struct {
		Data models.Order `json:"data"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &payload); err != nil {
		t.Fatalf("expected valid JSON response: %v", err)
	}

	if payload.Data.TotalCents != 25800 {
		t.Fatalf("expected total 25800, got %d", payload.Data.TotalCents)
	}
}

func TestOrderCreateFromCartRejectsEmptyCart(t *testing.T) {
	gin.SetMode(gin.TestMode)

	service := stubOrderService{
		createFromCart: func(ctx context.Context, userID string) (models.Order, error) {
			return models.Order{}, repositories.ErrEmptyCart
		},
	}

	router := gin.New()
	handler := handlers.NewOrderHandler(service)
	router.POST("/api/v1/orders", withUser("user-1"), handler.CreateFromCart)

	request := httptest.NewRequest(http.MethodPost, "/api/v1/orders", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, response.Code)
	}
}

func TestOrderListReturnsOrdersForUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	service := stubOrderService{
		listByUser: func(ctx context.Context, userID string) ([]models.Order, error) {
			return []models.Order{
				{
					ID:         "8d933930-8022-4df2-91cc-a7a28c0c53ef",
					UserID:     userID,
					Status:     models.OrderStatusConfirmed,
					TotalCents: 25800,
				},
			}, nil
		},
	}

	router := gin.New()
	handler := handlers.NewOrderHandler(service)
	router.GET("/api/v1/orders", withUser("user-1"), handler.List)

	request := httptest.NewRequest(http.MethodGet, "/api/v1/orders", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, response.Code)
	}

	var payload struct {
		Data []models.Order `json:"data"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &payload); err != nil {
		t.Fatalf("expected valid JSON response: %v", err)
	}

	if len(payload.Data) != 1 {
		t.Fatalf("expected 1 order, got %d", len(payload.Data))
	}
}

func TestOrderGetByIDRejectsInvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	service := stubOrderService{
		findByIDForUser: func(ctx context.Context, userID string, orderID string) (models.Order, error) {
			t.Fatal("service should not be called for invalid UUID")
			return models.Order{}, nil
		},
	}

	router := gin.New()
	handler := handlers.NewOrderHandler(service)
	router.GET("/api/v1/orders/:id", withUser("user-1"), handler.GetByID)

	request := httptest.NewRequest(http.MethodGet, "/api/v1/orders/not-a-uuid", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, response.Code)
	}
}

func TestOrderGetByIDReturnsNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	service := stubOrderService{
		findByIDForUser: func(ctx context.Context, userID string, orderID string) (models.Order, error) {
			return models.Order{}, repositories.ErrNotFound
		},
	}

	router := gin.New()
	handler := handlers.NewOrderHandler(service)
	router.GET("/api/v1/orders/:id", withUser("user-1"), handler.GetByID)

	request := httptest.NewRequest(http.MethodGet, "/api/v1/orders/8d933930-8022-4df2-91cc-a7a28c0c53ef", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, response.Code)
	}
}
