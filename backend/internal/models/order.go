package models

import "time"

const OrderStatusConfirmed = "confirmed"

type Order struct {
	ID         string      `json:"id"`
	UserID     string      `json:"user_id"`
	Status     string      `json:"status"`
	TotalCents int         `json:"total_cents"`
	Items      []OrderItem `json:"items"`
	CreatedAt  time.Time   `json:"created_at"`
	UpdatedAt  time.Time   `json:"updated_at"`
}

type OrderItem struct {
	ID             string    `json:"id"`
	OrderID        string    `json:"order_id"`
	ProductID      string    `json:"product_id"`
	ProductName    string    `json:"product_name"`
	UnitPriceCents int       `json:"unit_price_cents"`
	Quantity       int       `json:"quantity"`
	LineTotalCents int       `json:"line_total_cents"`
	CreatedAt      time.Time `json:"created_at"`
}
