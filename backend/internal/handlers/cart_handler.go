package handlers

import (
	"context"
	"errors"
	"net/http"

	"ai-e-commerce-app-go/backend/internal/models"
	"ai-e-commerce-app-go/backend/internal/repositories"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CartService interface {
	GetCart(ctx context.Context, userID string) (models.Cart, error)
	AddItem(ctx context.Context, userID string, input models.AddCartItemInput) (models.Cart, error)
	UpdateItem(ctx context.Context, userID string, productID string, input models.UpdateCartItemInput) (models.Cart, error)
	RemoveItem(ctx context.Context, userID string, productID string) (models.Cart, error)
	Clear(ctx context.Context, userID string) error
}

type CartHandler struct {
	service CartService
}

func NewCartHandler(service CartService) CartHandler {
	return CartHandler{service: service}
}

func (h CartHandler) Get(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	cart, err := h.service.GetCart(c.Request.Context(), userID)
	if err != nil {
		RecordError(c, err)
		JSONError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to get cart")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": cart,
	})
}

func (h CartHandler) AddItem(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	var input models.AddCartItemInput
	if err := c.ShouldBindJSON(&input); err != nil {
		JSONValidationError(c, "invalid cart item payload", err)
		return
	}

	cart, err := h.service.AddItem(c.Request.Context(), userID, input)
	if errors.Is(err, repositories.ErrNotFound) {
		JSONError(c, http.StatusNotFound, "NOT_FOUND", "product not found")
		return
	}

	if err != nil {
		RecordError(c, err)
		JSONError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to add item to cart")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": cart,
	})
}

func (h CartHandler) UpdateItem(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	productID := c.Param("product_id")
	if _, err := uuid.Parse(productID); err != nil {
		JSONError(c, http.StatusBadRequest, "BAD_REQUEST", "invalid product id")
		return
	}

	var input models.UpdateCartItemInput
	if err := c.ShouldBindJSON(&input); err != nil {
		JSONValidationError(c, "invalid cart item payload", err)
		return
	}

	cart, err := h.service.UpdateItem(c.Request.Context(), userID, productID, input)
	if errors.Is(err, repositories.ErrNotFound) {
		JSONError(c, http.StatusNotFound, "NOT_FOUND", "cart item not found")
		return
	}

	if err != nil {
		RecordError(c, err)
		JSONError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to update cart item")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": cart,
	})
}

func (h CartHandler) RemoveItem(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	productID := c.Param("product_id")
	if _, err := uuid.Parse(productID); err != nil {
		JSONError(c, http.StatusBadRequest, "BAD_REQUEST", "invalid product id")
		return
	}

	cart, err := h.service.RemoveItem(c.Request.Context(), userID, productID)
	if errors.Is(err, repositories.ErrNotFound) {
		JSONError(c, http.StatusNotFound, "NOT_FOUND", "cart item not found")
		return
	}

	if err != nil {
		RecordError(c, err)
		JSONError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to remove cart item")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": cart,
	})
}

func (h CartHandler) Clear(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	if err := h.service.Clear(c.Request.Context(), userID); err != nil {
		RecordError(c, err)
		JSONError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to clear cart")
		return
	}

	c.Status(http.StatusNoContent)
}

func getUserID(c *gin.Context) (string, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		JSONError(c, http.StatusUnauthorized, "UNAUTHORIZED", "authentication required")
		return "", false
	}

	return userID.(string), true
}
