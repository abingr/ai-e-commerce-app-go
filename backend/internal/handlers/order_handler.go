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

type OrderService interface {
	CreateFromCart(ctx context.Context, userID string) (models.Order, error)
	ListByUser(ctx context.Context, userID string) ([]models.Order, error)
	FindByIDForUser(ctx context.Context, userID string, orderID string) (models.Order, error)
}

type OrderHandler struct {
	service OrderService
}

func NewOrderHandler(service OrderService) OrderHandler {
	return OrderHandler{service: service}
}

func (h OrderHandler) CreateFromCart(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	order, err := h.service.CreateFromCart(c.Request.Context(), userID)
	if errors.Is(err, repositories.ErrEmptyCart) {
		JSONError(c, http.StatusBadRequest, "BAD_REQUEST", "cart is empty")
		return
	}

	if err != nil {
		RecordError(c, err)
		JSONError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to create order")
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": order,
	})
}

func (h OrderHandler) List(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	orders, err := h.service.ListByUser(c.Request.Context(), userID)
	if err != nil {
		RecordError(c, err)
		JSONError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to list orders")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": orders,
	})
}

func (h OrderHandler) GetByID(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	orderID := c.Param("id")
	if _, err := uuid.Parse(orderID); err != nil {
		JSONError(c, http.StatusBadRequest, "BAD_REQUEST", "invalid order id")
		return
	}

	order, err := h.service.FindByIDForUser(c.Request.Context(), userID, orderID)
	if errors.Is(err, repositories.ErrNotFound) {
		JSONError(c, http.StatusNotFound, "NOT_FOUND", "order not found")
		return
	}

	if err != nil {
		RecordError(c, err)
		JSONError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to get order")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": order,
	})
}
