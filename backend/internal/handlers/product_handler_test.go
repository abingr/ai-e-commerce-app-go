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

	"github.com/gin-gonic/gin"
)

type stubProductService struct {
	listProducts  func(ctx context.Context, filters models.ProductFilters) ([]models.Product, error)
	getProduct    func(ctx context.Context, id string) (models.Product, error)
	createProduct func(ctx context.Context, input models.ProductInput) (models.Product, error)
	updateProduct func(ctx context.Context, id string, input models.ProductInput) (models.Product, error)
	deleteProduct func(ctx context.Context, id string) error
}

func (s stubProductService) ListProducts(ctx context.Context, filters models.ProductFilters) ([]models.Product, error) {
	return s.listProducts(ctx, filters)
}

func (s stubProductService) GetProduct(ctx context.Context, id string) (models.Product, error) {
	return s.getProduct(ctx, id)
}

func (s stubProductService) CreateProduct(ctx context.Context, input models.ProductInput) (models.Product, error) {
	return s.createProduct(ctx, input)
}

func (s stubProductService) UpdateProduct(ctx context.Context, id string, input models.ProductInput) (models.Product, error) {
	return s.updateProduct(ctx, id, input)
}

func (s stubProductService) DeleteProduct(ctx context.Context, id string) error {
	return s.deleteProduct(ctx, id)
}

func TestProductListReturnsProducts(t *testing.T) {
	gin.SetMode(gin.TestMode)

	service := stubProductService{
		listProducts: func(ctx context.Context, filters models.ProductFilters) ([]models.Product, error) {
			if filters.Category != "Accessories" {
				t.Fatalf("expected category filter Accessories, got %q", filters.Category)
			}

			return []models.Product{
				{
					ID:            "f5e5a9b0-4055-421f-bd6f-7e755a16d1af",
					Name:          "KeyForge Mechanical Keyboard",
					Brand:         "KeyForge",
					Category:      "Accessories",
					PriceCents:    12900,
					StockQuantity: 40,
					IsActive:      true,
				},
			}, nil
		},
	}

	router := gin.New()
	handler := handlers.NewProductHandler(service)
	router.GET("/api/v1/products", handler.List)

	request := httptest.NewRequest(http.MethodGet, "/api/v1/products?category=Accessories", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, response.Code)
	}

	var body struct {
		Data []models.Product `json:"data"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatalf("expected valid JSON response: %v", err)
	}

	if len(body.Data) != 1 {
		t.Fatalf("expected 1 product, got %d", len(body.Data))
	}

	if body.Data[0].Name != "KeyForge Mechanical Keyboard" {
		t.Fatalf("expected keyboard product, got %q", body.Data[0].Name)
	}
}

func TestProductGetByIDReturnsProduct(t *testing.T) {
	gin.SetMode(gin.TestMode)

	service := stubProductService{
		getProduct: func(ctx context.Context, id string) (models.Product, error) {
			return models.Product{
				ID:         id,
				Name:       "NovaPhone X",
				Brand:      "NovaMobile",
				Category:   "Phones",
				PriceCents: 89900,
				IsActive:   true,
			}, nil
		},
	}

	router := gin.New()
	handler := handlers.NewProductHandler(service)
	router.GET("/api/v1/products/:id", handler.GetByID)

	request := httptest.NewRequest(http.MethodGet, "/api/v1/products/f5e5a9b0-4055-421f-bd6f-7e755a16d1af", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, response.Code)
	}

	var body struct {
		Data models.Product `json:"data"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatalf("expected valid JSON response: %v", err)
	}

	if body.Data.Name != "NovaPhone X" {
		t.Fatalf("expected NovaPhone X, got %q", body.Data.Name)
	}
}

func TestProductGetByIDRejectsInvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	service := stubProductService{
		getProduct: func(ctx context.Context, id string) (models.Product, error) {
			t.Fatal("service should not be called for invalid UUID")
			return models.Product{}, nil
		},
	}

	router := gin.New()
	handler := handlers.NewProductHandler(service)
	router.GET("/api/v1/products/:id", handler.GetByID)

	request := httptest.NewRequest(http.MethodGet, "/api/v1/products/not-a-uuid", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, response.Code)
	}
}

func TestProductGetByIDReturnsNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	service := stubProductService{
		getProduct: func(ctx context.Context, id string) (models.Product, error) {
			return models.Product{}, repositories.ErrNotFound
		},
	}

	router := gin.New()
	handler := handlers.NewProductHandler(service)
	router.GET("/api/v1/products/:id", handler.GetByID)

	request := httptest.NewRequest(http.MethodGet, "/api/v1/products/f5e5a9b0-4055-421f-bd6f-7e755a16d1af", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, response.Code)
	}
}

func TestProductListReturnsInternalServerError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	service := stubProductService{
		listProducts: func(ctx context.Context, filters models.ProductFilters) ([]models.Product, error) {
			return nil, errors.New("database unavailable")
		},
	}

	router := gin.New()
	handler := handlers.NewProductHandler(service)
	router.GET("/api/v1/products", handler.List)

	request := httptest.NewRequest(http.MethodGet, "/api/v1/products", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, response.Code)
	}
}

func TestProductCreateReturnsCreated(t *testing.T) {
	gin.SetMode(gin.TestMode)

	service := stubProductService{
		createProduct: func(ctx context.Context, input models.ProductInput) (models.Product, error) {
			return models.Product{
				ID:            "f5e5a9b0-4055-421f-bd6f-7e755a16d1af",
				Name:          input.Name,
				Brand:         input.Brand,
				Category:      input.Category,
				PriceCents:    input.PriceCents,
				StockQuantity: input.StockQuantity,
				IsActive:      true,
			}, nil
		},
	}

	router := gin.New()
	handler := handlers.NewProductHandler(service)
	router.POST("/api/v1/admin/products", handler.Create)

	body := bytes.NewBufferString(`{"name":"USB-C Docking Station","description":"Multi-port dock.","brand":"DockPro","category":"Accessories","price_cents":15900,"stock_quantity":15,"image_url":"https://example.com/dock.jpg"}`)
	request := httptest.NewRequest(http.MethodPost, "/api/v1/admin/products", body)
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, response.Code)
	}
}

func TestProductUpdateReturnsNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	service := stubProductService{
		updateProduct: func(ctx context.Context, id string, input models.ProductInput) (models.Product, error) {
			return models.Product{}, repositories.ErrNotFound
		},
	}

	router := gin.New()
	handler := handlers.NewProductHandler(service)
	router.PUT("/api/v1/admin/products/:id", handler.Update)

	body := bytes.NewBufferString(`{"name":"USB-C Docking Station","description":"Multi-port dock.","brand":"DockPro","category":"Accessories","price_cents":15900,"stock_quantity":15,"image_url":"https://example.com/dock.jpg"}`)
	request := httptest.NewRequest(http.MethodPut, "/api/v1/admin/products/f5e5a9b0-4055-421f-bd6f-7e755a16d1af", body)
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, response.Code)
	}
}

func TestProductDeleteReturnsNoContent(t *testing.T) {
	gin.SetMode(gin.TestMode)

	service := stubProductService{
		deleteProduct: func(ctx context.Context, id string) error {
			return nil
		},
	}

	router := gin.New()
	handler := handlers.NewProductHandler(service)
	router.DELETE("/api/v1/admin/products/:id", handler.Delete)

	request := httptest.NewRequest(http.MethodDelete, "/api/v1/admin/products/f5e5a9b0-4055-421f-bd6f-7e755a16d1af", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusNoContent {
		t.Fatalf("expected status %d, got %d", http.StatusNoContent, response.Code)
	}
}
