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
		JSONValidationError(c, "invalid registration payload", err)
		return
	}

	response, err := h.service.Register(c.Request.Context(), input)
	if errors.Is(err, repositories.ErrConflict) {
		JSONError(c, http.StatusConflict, "CONFLICT", "email already registered")
		return
	}

	if err != nil {
		RecordError(c, err)
		JSONError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to register user")
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": response,
	})
}

func (h AuthHandler) Login(c *gin.Context) {
	var input models.LoginUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		JSONValidationError(c, "invalid login payload", err)
		return
	}

	response, err := h.service.Login(c.Request.Context(), input)
	if errors.Is(err, services.ErrInvalidCredentials) {
		JSONError(c, http.StatusUnauthorized, "UNAUTHORIZED", "invalid email or password")
		return
	}

	if err != nil {
		RecordError(c, err)
		JSONError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to login")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": response,
	})
}

func (h AuthHandler) Me(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		JSONError(c, http.StatusUnauthorized, "UNAUTHORIZED", "authentication required")
		return
	}

	user, err := h.service.GetUser(c.Request.Context(), userID.(string))
	if errors.Is(err, repositories.ErrNotFound) {
		JSONError(c, http.StatusUnauthorized, "UNAUTHORIZED", "authentication required")
		return
	}

	if err != nil {
		RecordError(c, err)
		JSONError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to get current user")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": user,
	})
}
