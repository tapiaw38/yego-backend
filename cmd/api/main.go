package main

import (
	"log"

	orderRepo "wappi/internal/adapters/datasources/repositories/order"
	ordertokenRepo "wappi/internal/adapters/datasources/repositories/ordertoken"
	profileRepo "wappi/internal/adapters/datasources/repositories/profile"
	settingsRepo "wappi/internal/adapters/datasources/repositories/settings"
	"wappi/internal/adapters/web"
	"wappi/internal/adapters/web/middlewares"
	"wappi/internal/platform/config"
	"wappi/internal/platform/database"
	adminUsecase "wappi/internal/usecases/admin"
	orderUsecase "wappi/internal/usecases/order"
	profileUsecase "wappi/internal/usecases/profile"
	settingsUsecase "wappi/internal/usecases/settings"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg := config.GetInstance()
	log.Printf("Starting Order Tracking API on port %s", cfg.ServerPort)

	// Initialize database
	db := database.GetInstance()
	if err := database.RunMigrations(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize repositories
	orderRepository := orderRepo.NewRepository(db)
	orderTokenRepository := ordertokenRepo.NewRepository(db)
	profileRepository := profileRepo.NewRepository(db)
	settingsRepository := settingsRepo.NewRepository(db)

	// Initialize use cases
	orderUsecases := orderUsecase.NewUsecases(orderRepository, orderTokenRepository, profileRepository)
	profileUsecases := profileUsecase.NewUsecases(profileRepository)
	adminUsecases := adminUsecase.NewUsecases(profileRepository, orderRepository)
	settingsUsecases := settingsUsecase.NewUsecases(settingsRepository)

	// Setup Gin
	gin.SetMode(cfg.GinMode)
	app := gin.Default()

	// Apply CORS middleware
	app.Use(middlewares.CORSMiddleware())

	// Register routes
	web.RegisterRoutes(app, orderUsecases, profileUsecases, adminUsecases, settingsUsecases, cfg.FrontendURL)

	// Health check endpoint
	app.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Start server
	if err := app.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
