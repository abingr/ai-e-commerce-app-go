package routes

import (
	"log/slog"

	"ai-e-commerce-app-go/backend/internal/config"
	"ai-e-commerce-app-go/backend/internal/handlers"
	"ai-e-commerce-app-go/backend/internal/repositories"
	"ai-e-commerce-app-go/backend/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Dependencies struct {
	Config config.Config
	DB     *pgxpool.Pool
	Logger *slog.Logger
}

func New(deps Dependencies) *gin.Engine {
	if deps.Config.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(cors())
	router.Use(requestLogger(deps.Logger))

	healthHandler := handlers.NewHealthHandler(deps.Config, deps.DB)
	productRepository := repositories.NewProductRepository(deps.DB)
	productService := services.NewProductService(productRepository)
	productHandler := handlers.NewProductHandler(productService)

	router.GET("/health", healthHandler.Health)
	router.GET("/ready", healthHandler.Ready)

	api := router.Group("/api/v1")
	api.GET("/products", productHandler.List)
	api.GET("/products/:id", productHandler.GetByID)

	return router
}
