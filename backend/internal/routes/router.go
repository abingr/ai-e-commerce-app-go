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
	userRepository := repositories.NewUserRepository(deps.DB)
	authService := services.NewAuthService(userRepository, deps.Config.JWTSecret, deps.Config.JWTIssuer)
	authHandler := handlers.NewAuthHandler(authService)
	cartRepository := repositories.NewCartRepository(deps.DB)
	cartService := services.NewCartService(cartRepository)
	cartHandler := handlers.NewCartHandler(cartService)

	router.GET("/health", healthHandler.Health)
	router.GET("/ready", healthHandler.Ready)

	api := router.Group("/api/v1")
	api.POST("/auth/register", authHandler.Register)
	api.POST("/auth/login", authHandler.Login)
	api.GET("/products", productHandler.List)
	api.GET("/products/:id", productHandler.GetByID)
	api.GET("/me", requireAuth(authService), authHandler.Me)

	cart := api.Group("/cart", requireAuth(authService))
	cart.GET("", cartHandler.Get)
	cart.POST("/items", cartHandler.AddItem)
	cart.PATCH("/items/:product_id", cartHandler.UpdateItem)
	cart.DELETE("/items/:product_id", cartHandler.RemoveItem)
	cart.DELETE("/items", cartHandler.Clear)

	admin := api.Group("/admin", requireAuth(authService), requireRole("admin"))
	admin.POST("/products", productHandler.Create)
	admin.PUT("/products/:id", productHandler.Update)
	admin.DELETE("/products/:id", productHandler.Delete)

	return router
}
