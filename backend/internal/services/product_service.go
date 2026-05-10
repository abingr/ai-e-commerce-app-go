package services

import (
	"context"

	"ai-e-commerce-app-go/backend/internal/models"
)

type ProductRepository interface {
	List(ctx context.Context, filters models.ProductFilters) ([]models.Product, error)
	FindByID(ctx context.Context, id string) (models.Product, error)
	Create(ctx context.Context, input models.ProductInput) (models.Product, error)
	Update(ctx context.Context, id string, input models.ProductInput) (models.Product, error)
	Delete(ctx context.Context, id string) error
}

type ProductService struct {
	repository ProductRepository
}

func NewProductService(repository ProductRepository) ProductService {
	return ProductService{repository: repository}
}

func (s ProductService) ListProducts(ctx context.Context, filters models.ProductFilters) ([]models.Product, error) {
	return s.repository.List(ctx, filters)
}

func (s ProductService) GetProduct(ctx context.Context, id string) (models.Product, error) {
	return s.repository.FindByID(ctx, id)
}

func (s ProductService) CreateProduct(ctx context.Context, input models.ProductInput) (models.Product, error) {
	return s.repository.Create(ctx, input)
}

func (s ProductService) UpdateProduct(ctx context.Context, id string, input models.ProductInput) (models.Product, error) {
	return s.repository.Update(ctx, id, input)
}

func (s ProductService) DeleteProduct(ctx context.Context, id string) error {
	return s.repository.Delete(ctx, id)
}
