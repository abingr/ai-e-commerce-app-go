package services

import (
	"context"

	"ai-e-commerce-app-go/backend/internal/models"
)

type OrderRepository interface {
	CreateFromCart(ctx context.Context, userID string) (models.Order, error)
	ListByUser(ctx context.Context, userID string) ([]models.Order, error)
	FindByIDForUser(ctx context.Context, userID string, orderID string) (models.Order, error)
}

type OrderService struct {
	repository OrderRepository
}

func NewOrderService(repository OrderRepository) OrderService {
	return OrderService{repository: repository}
}

func (s OrderService) CreateFromCart(ctx context.Context, userID string) (models.Order, error) {
	return s.repository.CreateFromCart(ctx, userID)
}

func (s OrderService) ListByUser(ctx context.Context, userID string) ([]models.Order, error) {
	return s.repository.ListByUser(ctx, userID)
}

func (s OrderService) FindByIDForUser(ctx context.Context, userID string, orderID string) (models.Order, error) {
	return s.repository.FindByIDForUser(ctx, userID, orderID)
}
