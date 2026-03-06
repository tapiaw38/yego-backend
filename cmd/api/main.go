package main

import (
	"log"

	"yego/internal/adapters/datasources"
	"yego/internal/adapters/web"
	paymentHandler "yego/internal/adapters/web/handlers/payment"
	websocketHandler "yego/internal/adapters/web/handlers/websocket"
	"yego/internal/adapters/web/integrations"
	"yego/internal/adapters/web/middlewares"
	"yego/internal/platform/appcontext"
	"yego/internal/platform/config"
	"yego/internal/platform/database"
	s3service "yego/internal/services/s3"
	"yego/internal/usecases"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.GetInstance()
	log.Printf("Starting Order Tracking API on port %s", cfg.ServerPort)

	db := database.GetInstance()

	if err := database.RunMigrations(); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	ds := datasources.CreateDatasources(db)
	integrations := integrations.CreateIntegration(cfg)
	contextFactory := appcontext.NewFactory(ds, integrations, cfg)

	s3Client := s3service.NewClient(cfg.S3Region, cfg.S3Bucket, cfg.S3AccessKeyID, cfg.S3SecretAccessKey)
	useCases := usecases.CreateUsecases(contextFactory, s3Client)

	gin.SetMode(cfg.GinMode)
	app := gin.Default()

	app.Use(middlewares.CORSMiddleware())

	hub := integrations.WebSocket.GetHub()
	wsHandler := websocketHandler.NewHandler(hub)
	paymentCheckHandler := paymentHandler.NewHandler(contextFactory)
	web.RegisterRoutes(app, useCases, wsHandler, paymentCheckHandler, cfg)

	app.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	if err := app.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
