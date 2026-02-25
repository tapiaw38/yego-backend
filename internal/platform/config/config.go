package config

import (
	"os"

	"github.com/joho/godotenv"
)

// ConfigurationService holds all application configuration
type ConfigurationService struct {
	ServerPort              string
	DatabaseURL             string
	GinMode                 string
	FrontendURL             string
	BackendURL              string
	JWTSecret               string
	PaymentServiceURL       string
	PaymentAPIKey           string
	AuthAPIURL              string
	MPAccessToken string
}

var instance *ConfigurationService

// GetInstance returns the singleton configuration instance
func GetInstance() *ConfigurationService {
	if instance == nil {
		_ = godotenv.Load()

		instance = &ConfigurationService{
			ServerPort:       getEnvOrDefault("SERVER_PORT", "8080"),
			DatabaseURL:      getEnvOrDefault("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/yego?sslmode=disable"),
			GinMode:          getEnvOrDefault("GIN_MODE", "debug"),
			FrontendURL:         getEnvOrDefault("FRONTEND_URL", "http://localhost:5173"),
			BackendURL:          getEnvOrDefault("BACKEND_URL", ""),
			JWTSecret:           getEnvOrDefault("JWT_SECRET", ""),
			PaymentServiceURL:   getEnvOrDefault("PAYMENT_SERVICE_URL", "http://localhost:8008"),
			PaymentAPIKey:       getEnvOrDefault("PAYMENT_API_KEY", ""),
			AuthAPIURL:          getEnvOrDefault("AUTH_API_URL", "http://localhost:8082"),
			MPAccessToken: getEnvOrDefault("MP_ACCESS_TOKEN", ""),
		}
	}
	return instance
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
