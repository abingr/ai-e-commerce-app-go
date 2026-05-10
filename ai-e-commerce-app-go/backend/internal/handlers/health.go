package handlers

import (
	"context"
	"net/http"
	"time"

	"ai-e-commerce-app-go/backend/internal/config"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type HealthHandler struct {
	config config.Config
	db     *pgxpool.Pool
}

func NewHealthHandler(config config.Config, db *pgxpool.Pool) HealthHandler {
	return HealthHandler{
		config: config,
		db:     db,
	}
}

func (h HealthHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"service": h.config.AppName,
	})
}

func (h HealthHandler) Ready(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()

	if err := h.db.Ping(ctx); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status":   "not_ready",
			"database": "disconnected",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":   "ready",
		"database": "connected",
	})
}
