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

type ProductService interface {
	ListProducts(ctx context.Context, filters models.ProductFilters) ([]models.Product, error)
	GetProduct(ctx context.Context, id string) (models.Product, error)
	CreateProduct(ctx context.Context, input models.ProductInput) (models.Product, error)
	UpdateProduct(ctx context.Context, id string, input models.ProductInput) (models.Product, error)
	DeleteProduct(ctx context.Context, id string) error
}

type ProductHandler struct {
	service ProductService
}

func NewProductHandler(service ProductService) ProductHandler {
	return ProductHandler{service: service}
}

func (h ProductHandler) List(c *gin.Context) {
	filters := models.ProductFilters{
		Category: c.Query("category"),
		Search:   c.Query("search"),
	}

	products, err := h.service.ListProducts(c.Request.Context(), filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to list products",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": products,
	})
}

func (h ProductHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	if _, err := uuid.Parse(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid product id",
		})
		return
	}

	product, err := h.service.GetProduct(c.Request.Context(), id)
	if errors.Is(err, repositories.ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "product not found",
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get product",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": product,
	})
}

func (h ProductHandler) Create(c *gin.Context) {
	var input models.ProductInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid product payload",
		})
		return
	}

	product, err := h.service.CreateProduct(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create product",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": product,
	})
}

func (h ProductHandler) Update(c *gin.Context) {
	id := c.Param("id")
	if _, err := uuid.Parse(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid product id",
		})
		return
	}

	var input models.ProductInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid product payload",
		})
		return
	}

	product, err := h.service.UpdateProduct(c.Request.Context(), id, input)
	if errors.Is(err, repositories.ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "product not found",
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to update product",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": product,
	})
}

func (h ProductHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if _, err := uuid.Parse(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid product id",
		})
		return
	}

	err := h.service.DeleteProduct(c.Request.Context(), id)
	if errors.Is(err, repositories.ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "product not found",
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to delete product",
		})
		return
	}

	c.Status(http.StatusNoContent)
}
