package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv      string
	AppName     string
	HTTPPort    string
	DatabaseURL string
	JWTSecret   string
	JWTIssuer   string
}

func Load() Config {
	_ = godotenv.Load()

	return Config{
		AppEnv:      getEnv("ECOMMERCE_APP_ENV", "development"),
		AppName:     getEnv("ECOMMERCE_APP_NAME", "ai-e-commerce-api"),
		HTTPPort:    getEnv("ECOMMERCE_HTTP_PORT", "8080"),
		DatabaseURL: getEnv("ECOMMERCE_DATABASE_URL", "postgres://ecommerce_user:ecommerce_password@127.0.0.1:55432/ecommerce_db?sslmode=disable"),
		JWTSecret:   getEnv("ECOMMERCE_JWT_SECRET", "change-this-development-secret"),
		JWTIssuer:   getEnv("ECOMMERCE_JWT_ISSUER", "ai-e-commerce-api"),
	}
}

func getEnv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}
