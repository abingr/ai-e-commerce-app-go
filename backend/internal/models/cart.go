package models

type CartItem struct {
	ProductID      string `json:"product_id"`
	Name           string `json:"name"`
	Brand          string `json:"brand"`
	Category       string `json:"category"`
	ImageURL       string `json:"image_url"`
	Quantity       int    `json:"quantity"`
	UnitPriceCents int    `json:"unit_price_cents"`
	LineTotalCents int    `json:"line_total_cents"`
}

type Cart struct {
	Items      []CartItem `json:"items"`
	TotalCents int        `json:"total_cents"`
}

type AddCartItemInput struct {
	ProductID string `json:"product_id" binding:"required,uuid"`
	Quantity  int    `json:"quantity" binding:"required,min=1"`
}

type UpdateCartItemInput struct {
	Quantity int `json:"quantity" binding:"required,min=1"`
}
