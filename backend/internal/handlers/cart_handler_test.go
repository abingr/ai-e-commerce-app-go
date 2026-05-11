package handlers_test

import (
	"bytes"
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

type stubCartService struct {
	getCart    func(ctx context.Context, userID string) (models.Cart, error)
	addItem    func(ctx context.Context, userID string, input models.AddCartItemInput) (models.Cart, error)
	updateItem func(ctx context.Context, userID string, productID string, input models.UpdateCartItemInput) (models.Cart, error)
	removeItem func(ctx context.Context, userID string, productID string) (models.Cart, error)
	clear      func(ctx context.Context, userID string) error
}

func (s stubCartService) GetCart(ctx context.Context, userID string) (models.Cart, error) {
	return s.getCart(ctx, userID)
}

func (s stubCartService) AddItem(ctx context.Context, userID string, input models.AddCartItemInput) (models.Cart, error) {
	return s.addItem(ctx, userID, input)
}

func (s stubCartService) UpdateItem(ctx context.Context, userID string, productID string, input models.UpdateCartItemInput) (models.Cart, error) {
	return s.updateItem(ctx, userID, productID, input)
}

func (s stubCartService) RemoveItem(ctx context.Context, userID string, productID string) (models.Cart, error) {
	return s.removeItem(ctx, userID, productID)
}

func (s stubCartService) Clear(ctx context.Context, userID string) error {
	return s.clear(ctx, userID)
}

func TestCartGetReturnsCart(t *testing.T) {
	gin.SetMode(gin.TestMode)

	service := stubCartService{
		getCart: func(ctx context.Context, userID string) (models.Cart, error) {
			if userID != "user-1" {
				t.Fatalf("expected user-1, got %q", userID)
			}

			return models.Cart{
				Items: []models.CartItem{
					{
						ProductID:      "f5e5a9b0-4055-421f-bd6f-7e755a16d1af",
						Name:           "Keyboard",
						Quantity:       2,
						UnitPriceCents: 12900,
						LineTotalCents: 25800,
					},
				},
				TotalCents: 25800,
			}, nil
		},
	}

	router := gin.New()
	handler := handlers.NewCartHandler(service)
	router.GET("/api/v1/cart", withUser("user-1"), handler.Get)

	request := httptest.NewRequest(http.MethodGet, "/api/v1/cart", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, response.Code)
	}

	var payload struct {
		Data models.Cart `json:"data"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &payload); err != nil {
		t.Fatalf("expected valid JSON response: %v", err)
	}

	if payload.Data.TotalCents != 25800 {
		t.Fatalf("expected total 25800, got %d", payload.Data.TotalCents)
	}
}

func TestCartAddItemReturnsUpdatedCart(t *testing.T) {
	gin.SetMode(gin.TestMode)

	service := stubCartService{
		addItem: func(ctx context.Context, userID string, input models.AddCartItemInput) (models.Cart, error) {
			if input.Quantity != 2 {
				t.Fatalf("expected quantity 2, got %d", input.Quantity)
			}

			return models.Cart{TotalCents: 31800}, nil
		},
	}

	router := gin.New()
	handler := handlers.NewCartHandler(service)
	router.POST("/api/v1/cart/items", withUser("user-1"), handler.AddItem)

	body := bytes.NewBufferString(`{"product_id":"f5e5a9b0-4055-421f-bd6f-7e755a16d1af","quantity":2}`)
	request := httptest.NewRequest(http.MethodPost, "/api/v1/cart/items", body)
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, response.Code)
	}
}

func TestCartAddItemReturnsNotFoundForMissingProduct(t *testing.T) {
	gin.SetMode(gin.TestMode)

	service := stubCartService{
		addItem: func(ctx context.Context, userID string, input models.AddCartItemInput) (models.Cart, error) {
			return models.Cart{}, repositories.ErrNotFound
		},
	}

	router := gin.New()
	handler := handlers.NewCartHandler(service)
	router.POST("/api/v1/cart/items", withUser("user-1"), handler.AddItem)

	body := bytes.NewBufferString(`{"product_id":"f5e5a9b0-4055-421f-bd6f-7e755a16d1af","quantity":2}`)
	request := httptest.NewRequest(http.MethodPost, "/api/v1/cart/items", body)
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, response.Code)
	}
}

func TestCartUpdateItemRejectsInvalidProductID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	service := stubCartService{
		updateItem: func(ctx context.Context, userID string, productID string, input models.UpdateCartItemInput) (models.Cart, error) {
			t.Fatal("service should not be called for invalid UUID")
			return models.Cart{}, nil
		},
	}

	router := gin.New()
	handler := handlers.NewCartHandler(service)
	router.PATCH("/api/v1/cart/items/:product_id", withUser("user-1"), handler.UpdateItem)

	body := bytes.NewBufferString(`{"quantity":3}`)
	request := httptest.NewRequest(http.MethodPatch, "/api/v1/cart/items/not-a-uuid", body)
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, response.Code)
	}
}

func TestCartRemoveItemReturnsUpdatedCart(t *testing.T) {
	gin.SetMode(gin.TestMode)

	service := stubCartService{
		removeItem: func(ctx context.Context, userID string, productID string) (models.Cart, error) {
			return models.Cart{Items: []models.CartItem{}, TotalCents: 0}, nil
		},
	}

	router := gin.New()
	handler := handlers.NewCartHandler(service)
	router.DELETE("/api/v1/cart/items/:product_id", withUser("user-1"), handler.RemoveItem)

	request := httptest.NewRequest(http.MethodDelete, "/api/v1/cart/items/f5e5a9b0-4055-421f-bd6f-7e755a16d1af", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, response.Code)
	}
}

func TestCartClearReturnsNoContent(t *testing.T) {
	gin.SetMode(gin.TestMode)

	service := stubCartService{
		clear: func(ctx context.Context, userID string) error {
			return nil
		},
	}

	router := gin.New()
	handler := handlers.NewCartHandler(service)
	router.DELETE("/api/v1/cart/items", withUser("user-1"), handler.Clear)

	request := httptest.NewRequest(http.MethodDelete, "/api/v1/cart/items", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusNoContent {
		t.Fatalf("expected status %d, got %d", http.StatusNoContent, response.Code)
	}
}

func withUser(userID string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("user_id", userID)
	}
}
