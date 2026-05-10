package services

import (
	"context"

	"ai-e-commerce-app-go/backend/internal/models"
)

type ProductRepository interface {
	List(ctx context.Context, filters models.ProductFilters) ([]models.Product, error)
	FindByID(ctx context.Context, id string) (models.Product, error)
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
