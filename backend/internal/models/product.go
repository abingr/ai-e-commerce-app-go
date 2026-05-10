package models

import "time"

type Product struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	Brand         string    `json:"brand"`
	Category      string    `json:"category"`
	PriceCents    int       `json:"price_cents"`
	StockQuantity int       `json:"stock_quantity"`
	ImageURL      string    `json:"image_url"`
	IsActive      bool      `json:"is_active"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type ProductFilters struct {
	Category string
	Search   string
}

type ProductInput struct {
	Name          string `json:"name" binding:"required"`
	Description   string `json:"description" binding:"required"`
	Brand         string `json:"brand" binding:"required"`
	Category      string `json:"category" binding:"required"`
	PriceCents    int    `json:"price_cents" binding:"required,min=0"`
	StockQuantity int    `json:"stock_quantity" binding:"min=0"`
	ImageURL      string `json:"image_url"`
}
