package services_test

import (
	"context"
	"testing"

	"ai-e-commerce-app-go/backend/internal/models"
	"ai-e-commerce-app-go/backend/internal/services"
)

type stubCartRepository struct {
	cart       models.Cart
	addCalled  bool
	updateSeen int
}

func (r *stubCartRepository) Get(ctx context.Context, userID string) (models.Cart, error) {
	return r.cart, nil
}

func (r *stubCartRepository) AddItem(ctx context.Context, userID string, input models.AddCartItemInput) error {
	r.addCalled = true
	return nil
}

func (r *stubCartRepository) UpdateItem(ctx context.Context, userID string, productID string, quantity int) error {
	r.updateSeen = quantity
	return nil
}

func (r *stubCartRepository) RemoveItem(ctx context.Context, userID string, productID string) error {
	return nil
}

func (r *stubCartRepository) Clear(ctx context.Context, userID string) error {
	return nil
}

func TestCartServiceAddItemReturnsFreshCart(t *testing.T) {
	repository := &stubCartRepository{
		cart: models.Cart{TotalCents: 31800},
	}
	service := services.NewCartService(repository)

	cart, err := service.AddItem(context.Background(), "user-1", models.AddCartItemInput{
		ProductID: "f5e5a9b0-4055-421f-bd6f-7e755a16d1af",
		Quantity:  2,
	})
	if err != nil {
		t.Fatalf("add item: %v", err)
	}

	if !repository.addCalled {
		t.Fatal("expected repository add to be called")
	}

	if cart.TotalCents != 31800 {
		t.Fatalf("expected total 31800, got %d", cart.TotalCents)
	}
}

func TestCartServiceUpdateItemPassesExactQuantity(t *testing.T) {
	repository := &stubCartRepository{}
	service := services.NewCartService(repository)

	_, err := service.UpdateItem(context.Background(), "user-1", "f5e5a9b0-4055-421f-bd6f-7e755a16d1af", models.UpdateCartItemInput{
		Quantity: 3,
	})
	if err != nil {
		t.Fatalf("update item: %v", err)
	}

	if repository.updateSeen != 3 {
		t.Fatalf("expected exact quantity 3, got %d", repository.updateSeen)
	}
}
