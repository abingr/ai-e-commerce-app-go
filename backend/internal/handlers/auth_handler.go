package handlers

import (
	"context"
	"errors"
	"net/http"

	"ai-e-commerce-app-go/backend/internal/models"
	"ai-e-commerce-app-go/backend/internal/repositories"
	"ai-e-commerce-app-go/backend/internal/services"

	"github.com/gin-gonic/gin"
)

type AuthService interface {
	Register(ctx context.Context, input models.RegisterUserInput) (models.AuthResponse, error)
	Login(ctx context.Context, input models.LoginUserInput) (models.AuthResponse, error)
	GetUser(ctx context.Context, id string) (models.User, error)
}

type AuthHandler struct {
	service AuthService
}

func NewAuthHandler(service AuthService) AuthHandler {
	return AuthHandler{service: service}
}

func (h AuthHandler) Register(c *gin.Context) {
	var input models.RegisterUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid registration payload",
		})
		return
	}

	response, err := h.service.Register(c.Request.Context(), input)
	if errors.Is(err, repositories.ErrConflict) {
		c.JSON(http.StatusConflict, gin.H{
			"error": "email already registered",
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to register user",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": response,
	})
}

func (h AuthHandler) Login(c *gin.Context) {
	var input models.LoginUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid login payload",
		})
		return
	}

	response, err := h.service.Login(c.Request.Context(), input)
	if errors.Is(err, services.ErrInvalidCredentials) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid email or password",
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to login",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": response,
	})
}

func (h AuthHandler) Me(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "authentication required",
		})
		return
	}

	user, err := h.service.GetUser(c.Request.Context(), userID.(string))
	if errors.Is(err, repositories.ErrNotFound) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "authentication required",
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get current user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": user,
	})
}
