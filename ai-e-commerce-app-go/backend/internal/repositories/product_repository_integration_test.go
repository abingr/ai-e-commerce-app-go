package repositories_test

import (
	"context"
	"os"
	"testing"
	"time"

	"ai-e-commerce-app-go/backend/internal/config"
	"ai-e-commerce-app-go/backend/internal/database"
	"ai-e-commerce-app-go/backend/internal/models"
	"ai-e-commerce-app-go/backend/internal/repositories"
)

func TestProductRepositoryListIntegration(t *testing.T) {
	if os.Getenv("ECOMMERCE_RUN_INTEGRATION_TESTS") != "true" {
		t.Skip("set ECOMMERCE_RUN_INTEGRATION_TESTS=true after running migrations to enable this database test")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cfg := config.Load()
	db, err := database.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		t.Fatalf("connect database: %v", err)
	}
	defer db.Close()

	repository := repositories.NewProductRepository(db)

	products, err := repository.List(ctx, models.ProductFilters{})
	if err != nil {
		t.Fatalf("list products: %v", err)
	}

	if len(products) == 0 {
		t.Fatal("expected seeded products, got none")
	}
}
