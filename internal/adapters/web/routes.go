package web

import (
	adminHandler "wappi/internal/adapters/web/handlers/admin"
	orderHandler "wappi/internal/adapters/web/handlers/order"
	profileHandler "wappi/internal/adapters/web/handlers/profile"
	"wappi/internal/adapters/web/middlewares"
	adminUsecase "wappi/internal/usecases/admin"
	orderUsecase "wappi/internal/usecases/order"
	profileUsecase "wappi/internal/usecases/profile"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers all application routes
func RegisterRoutes(app *gin.Engine, orderUsecases *orderUsecase.Usecases, profileUsecases *profileUsecase.Usecases, adminUsecases *adminUsecase.Usecases, frontendURL string) {
	api := app.Group("/api")

	// Public order routes (tracking by UUID - no auth needed)
	orders := api.Group("/orders")
	{
		orders.GET("/:id", orderHandler.NewGetHandler(orderUsecases.Get))
		// WhatsApp IA endpoint - creates order with claim link (public for IA access)
		orders.POST("/create-with-link", orderHandler.NewCreateWithLinkHandler(orderUsecases.CreateWithLink, frontendURL))
	}

	// Protected order routes (require auth)
	ordersAuth := api.Group("/orders")
	ordersAuth.Use(middlewares.AuthMiddleware())
	{
		ordersAuth.POST("", orderHandler.NewCreateHandler(orderUsecases.Create))
		ordersAuth.PATCH("/:id/status", orderHandler.NewUpdateStatusHandler(orderUsecases.UpdateStatus))
		// Claim order via token - requires auth to identify user
		ordersAuth.POST("/claim/:token", orderHandler.NewClaimHandler(orderUsecases.Claim))
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
		profilesAuth.POST("", profileHandler.NewCreateOrUpdateHandler(profileUsecases.CreateOrUpdate))
		profilesAuth.PUT("", profileHandler.NewCreateOrUpdateHandler(profileUsecases.CreateOrUpdate))
	}

	// Admin routes (open for now - no authentication)
	admin := api.Group("/admin")
	{
		admin.GET("/profiles", adminHandler.NewListProfilesHandler(adminUsecases.ListProfiles))
		admin.GET("/orders", adminHandler.NewListOrdersHandler(adminUsecases.ListOrders))
		admin.PUT("/orders/:id", adminHandler.NewUpdateOrderHandler(adminUsecases.UpdateOrder))
	}
}
