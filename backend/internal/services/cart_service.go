package services

import (
	"context"

	"ai-e-commerce-app-go/backend/internal/models"
)

type CartRepository interface {
	Get(ctx context.Context, userID string) (models.Cart, error)
	AddItem(ctx context.Context, userID string, input models.AddCartItemInput) error
	UpdateItem(ctx context.Context, userID string, productID string, quantity int) error
	RemoveItem(ctx context.Context, userID string, productID string) error
	Clear(ctx context.Context, userID string) error
}

type CartService struct {
	repository CartRepository
}

func NewCartService(repository CartRepository) CartService {
	return CartService{repository: repository}
}

func (s CartService) GetCart(ctx context.Context, userID string) (models.Cart, error) {
	return s.repository.Get(ctx, userID)
}

func (s CartService) AddItem(ctx context.Context, userID string, input models.AddCartItemInput) (models.Cart, error) {
	if err := s.repository.AddItem(ctx, userID, input); err != nil {
		return models.Cart{}, err
	}

	return s.repository.Get(ctx, userID)
}

func (s CartService) UpdateItem(ctx context.Context, userID string, productID string, input models.UpdateCartItemInput) (models.Cart, error) {
	if err := s.repository.UpdateItem(ctx, userID, productID, input.Quantity); err != nil {
		return models.Cart{}, err
	}

	return s.repository.Get(ctx, userID)
}

func (s CartService) RemoveItem(ctx context.Context, userID string, productID string) (models.Cart, error) {
	if err := s.repository.RemoveItem(ctx, userID, productID); err != nil {
		return models.Cart{}, err
	}

	return s.repository.Get(ctx, userID)
}

func (s CartService) Clear(ctx context.Context, userID string) error {
	return s.repository.Clear(ctx, userID)
}
