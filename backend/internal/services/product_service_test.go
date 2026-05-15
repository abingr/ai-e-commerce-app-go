package services_test

import (
	"context"
	"testing"

	"ai-e-commerce-app-go/backend/internal/models"
	"ai-e-commerce-app-go/backend/internal/services"
)

type stubProductRepositoryForService struct {
	filtersSeen models.ProductFilters
	createdSeen models.ProductInput
	deletedID   string
}

func (r *stubProductRepositoryForService) List(ctx context.Context, filters models.ProductFilters) ([]models.Product, error) {
	r.filtersSeen = filters
	return []models.Product{{ID: "f5e5a9b0-4055-421f-bd6f-7e755a16d1af", Name: "Keyboard"}}, nil
}

func (r *stubProductRepositoryForService) FindByID(ctx context.Context, id string) (models.Product, error) {
	return models.Product{ID: id, Name: "Keyboard"}, nil
}

func (r *stubProductRepositoryForService) Create(ctx context.Context, input models.ProductInput) (models.Product, error) {
	r.createdSeen = input
	return models.Product{Name: input.Name, PriceCents: input.PriceCents}, nil
}

func (r *stubProductRepositoryForService) Update(ctx context.Context, id string, input models.ProductInput) (models.Product, error) {
	return models.Product{ID: id, Name: input.Name}, nil
}

func (r *stubProductRepositoryForService) Delete(ctx context.Context, id string) error {
	r.deletedID = id
	return nil
}

func TestProductServicePassesFiltersToRepository(t *testing.T) {
	repository := &stubProductRepositoryForService{}
	service := services.NewProductService(repository)

	_, err := service.ListProducts(context.Background(), models.ProductFilters{
		Category: "Accessories",
		Search:   "keyboard",
	})
	if err != nil {
		t.Fatalf("list products: %v", err)
	}

	if repository.filtersSeen.Category != "Accessories" {
		t.Fatalf("expected category filter, got %q", repository.filtersSeen.Category)
	}

	if repository.filtersSeen.Search != "keyboard" {
		t.Fatalf("expected search filter, got %q", repository.filtersSeen.Search)
	}
}

func TestProductServiceCreateDelegatesInput(t *testing.T) {
	repository := &stubProductRepositoryForService{}
	service := services.NewProductService(repository)

	_, err := service.CreateProduct(context.Background(), models.ProductInput{
		Name:       "USB-C Dock",
		PriceCents: 15900,
	})
	if err != nil {
		t.Fatalf("create product: %v", err)
	}

	if repository.createdSeen.Name != "USB-C Dock" {
		t.Fatalf("expected product name to be delegated, got %q", repository.createdSeen.Name)
	}
}

func TestProductServiceDeleteDelegatesID(t *testing.T) {
	repository := &stubProductRepositoryForService{}
	service := services.NewProductService(repository)

	err := service.DeleteProduct(context.Background(), "f5e5a9b0-4055-421f-bd6f-7e755a16d1af")
	if err != nil {
		t.Fatalf("delete product: %v", err)
	}

	if repository.deletedID != "f5e5a9b0-4055-421f-bd6f-7e755a16d1af" {
		t.Fatalf("expected deleted id to be delegated, got %q", repository.deletedID)
	}
}
