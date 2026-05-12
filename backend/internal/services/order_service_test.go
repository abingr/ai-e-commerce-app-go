package services_test

import (
	"context"
	"testing"

	"ai-e-commerce-app-go/backend/internal/models"
	"ai-e-commerce-app-go/backend/internal/services"
)

type stubOrderRepository struct {
	createCalled bool
	userSeen     string
}

func (r *stubOrderRepository) CreateFromCart(ctx context.Context, userID string) (models.Order, error) {
	r.createCalled = true
	r.userSeen = userID
	return models.Order{ID: "8d933930-8022-4df2-91cc-a7a28c0c53ef", UserID: userID}, nil
}

func (r *stubOrderRepository) ListByUser(ctx context.Context, userID string) ([]models.Order, error) {
	r.userSeen = userID
	return []models.Order{{ID: "8d933930-8022-4df2-91cc-a7a28c0c53ef", UserID: userID}}, nil
}

func (r *stubOrderRepository) FindByIDForUser(ctx context.Context, userID string, orderID string) (models.Order, error) {
	r.userSeen = userID
	return models.Order{ID: orderID, UserID: userID}, nil
}

func TestOrderServiceCreateFromCartDelegatesToRepository(t *testing.T) {
	repository := &stubOrderRepository{}
	service := services.NewOrderService(repository)

	order, err := service.CreateFromCart(context.Background(), "user-1")
	if err != nil {
		t.Fatalf("create order: %v", err)
	}

	if !repository.createCalled {
		t.Fatal("expected repository create to be called")
	}

	if order.UserID != "user-1" {
		t.Fatalf("expected user-1, got %q", order.UserID)
	}
}
