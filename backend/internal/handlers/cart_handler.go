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
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get cart",
		})
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
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid cart item payload",
		})
		return
	}

	cart, err := h.service.AddItem(c.Request.Context(), userID, input)
	if errors.Is(err, repositories.ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "product not found",
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to add item to cart",
		})
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
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid product id",
		})
		return
	}

	var input models.UpdateCartItemInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid cart item payload",
		})
		return
	}

	cart, err := h.service.UpdateItem(c.Request.Context(), userID, productID, input)
	if errors.Is(err, repositories.ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "cart item not found",
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to update cart item",
		})
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
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid product id",
		})
		return
	}

	cart, err := h.service.RemoveItem(c.Request.Context(), userID, productID)
	if errors.Is(err, repositories.ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "cart item not found",
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to remove cart item",
		})
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
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to clear cart",
		})
		return
	}

	c.Status(http.StatusNoContent)
}

func getUserID(c *gin.Context) (string, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "authentication required",
		})
		return "", false
	}

	return userID.(string), true
}
