package web

import (
	"github.com/gin-gonic/gin"
	adminHandler "wappi/internal/adapters/web/handlers/admin"
	orderHandler "wappi/internal/adapters/web/handlers/order"
	profileHandler "wappi/internal/adapters/web/handlers/profile"
	settingsHandler "wappi/internal/adapters/web/handlers/settings"
	"wappi/internal/adapters/web/middlewares"
	adminUsecase "wappi/internal/usecases/admin"
	orderUsecase "wappi/internal/usecases/order"
	profileUsecase "wappi/internal/usecases/profile"
	settingsUsecase "wappi/internal/usecases/settings"
)

// RegisterRoutes registers all application routes
func RegisterRoutes(app *gin.Engine, orderUsecases *orderUsecase.Usecases, profileUsecases *profileUsecase.Usecases, adminUsecases *adminUsecase.Usecases, settingsUsecases *settingsUsecase.Usecases, frontendURL string) {
	api := app.Group("/api")

	// Public order routes (tracking by UUID - no auth needed)
	orders := api.Group("/orders")
	{
		orders.GET("/:id", orderHandler.NewGetHandler(orderUsecases.Get))
		// Create order with link (admin functionality, no auth for now)
		orders.POST("/create-with-link", orderHandler.NewCreateWithLinkHandler(orderUsecases.CreateWithLink, frontendURL))
	}

	// Protected order routes (require auth)
	ordersAuth := api.Group("/orders")
	ordersAuth.Use(middlewares.AuthMiddleware())
	{
		ordersAuth.POST("", orderHandler.NewCreateHandler(orderUsecases.Create))
		ordersAuth.PATCH("/:id/status", orderHandler.NewUpdateStatusHandler(orderUsecases.UpdateStatus))
		// Claim order with token
		ordersAuth.POST("/claim/:token", orderHandler.NewClaimHandler(orderUsecases.Claim))
		// List current user's orders
		ordersAuth.GET("/my", orderHandler.NewListMyHandler(orderUsecases.ListMyOrders))
	}

	// Public profile routes (token-based access)
	profiles := api.Group("/profiles")
	{
		profiles.GET("/validate/:token", profileHandler.NewValidateTokenHandler(profileUsecases.ValidateToken))
		profiles.POST("/complete", profileHandler.NewCompleteProfileHandler(profileUsecases.CompleteProfile))
		profiles.GET("/:id", profileHandler.NewGetHandler(profileUsecases.Get))
		profiles.PUT("/:id", profileHandler.NewUpdateHandler(profileUsecases.Update))
	}

	// Protected profile routes (require auth)
	profilesAuth := api.Group("/profiles")
	profilesAuth.Use(middlewares.AuthMiddleware())
	{
		profilesAuth.POST("/generate-link", profileHandler.NewGenerateLinkHandler(profileUsecases.GenerateLink))
		profilesAuth.GET("/check-completed", profileHandler.NewCheckCompletedHandler(profileUsecases.CheckCompleted))
	}

	// Settings routes (public for reading, could add auth for updating)
	settings := api.Group("/settings")
	{
		settings.GET("", settingsHandler.NewGetHandler(settingsUsecases.Get))
		settings.PUT("", settingsHandler.NewUpdateHandler(settingsUsecases.Update))
		settings.POST("/calculate-delivery", settingsHandler.NewCalculateDeliveryFeeHandler(settingsUsecases.CalculateDeliveryFee))
	}

	// Admin routes (require auth)
	admin := api.Group("/admin")
	admin.Use(middlewares.AuthMiddleware())
	{
		admin.GET("/profiles", adminHandler.NewListProfilesHandler(adminUsecases.ListProfiles))
		admin.GET("/orders", adminHandler.NewListOrdersHandler(adminUsecases.ListOrders))
		admin.PUT("/orders/:id", adminHandler.NewUpdateOrderHandler(adminUsecases.UpdateOrder))
	}
}
