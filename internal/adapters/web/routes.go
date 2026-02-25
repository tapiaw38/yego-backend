package web

import (
	"github.com/gin-gonic/gin"
	adminHandler "yego/internal/adapters/web/handlers/admin"
	orderHandler "yego/internal/adapters/web/handlers/order"
	paymentHandler "yego/internal/adapters/web/handlers/payment"
	profileHandler "yego/internal/adapters/web/handlers/profile"
	settingsHandler "yego/internal/adapters/web/handlers/settings"
	websocketHandler "yego/internal/adapters/web/handlers/websocket"
	"yego/internal/adapters/web/middlewares"
	"yego/internal/platform/config"
	"yego/internal/usecases"
)

// RegisterRoutes registers all application routes
func RegisterRoutes(app *gin.Engine, useCases *usecases.Usecases, wsHandler *websocketHandler.Handler, paymentCheckHandler *paymentHandler.Handler, cfg *config.ConfigurationService) {
	api := app.Group("/api")

	// Public order routes (tracking by UUID - no auth needed)
	orders := api.Group("/orders")
	{
		orders.GET("/:id", orderHandler.NewGetHandler(useCases.Order.GetUsecase))
		orders.POST("/create-with-link", orderHandler.NewCreateWithLinkHandler(useCases.Order.CreateWithLinkUsecase, cfg.FrontendURL))
		orders.GET("/claim/:token/info", orderHandler.NewGetClaimInfoHandler(useCases.Order.GetClaimInfoUsecase))
		// MercadoPago webhook — called by MP servers, no auth
		orders.POST("/webhook/mp", orderHandler.NewPaymentWebhookHandler(useCases.Order.HandlePaymentWebhookUsecase))
	}

	// Protected order routes (require auth)
	ordersAuth := api.Group("/orders")
	ordersAuth.Use(middlewares.AuthMiddleware())
	{
		ordersAuth.POST("", orderHandler.NewCreateHandler(useCases.Order.CreateUsecase))
		ordersAuth.PATCH("/:id/status", orderHandler.NewUpdateStatusHandler(useCases.Order.UpdateStatusUsecase))
		ordersAuth.POST("/claim/:token", orderHandler.NewClaimHandler(useCases.Order.ClaimUsecase))
		ordersAuth.POST("/:id/pay", orderHandler.NewPayForOrderHandler(useCases.Order.PayForOrderUsecase))
		ordersAuth.POST("/:id/payment-link", orderHandler.NewCreatePaymentLinkHandler(useCases.Order.CreatePaymentLinkUsecase, cfg.FrontendURL, cfg.BackendURL))
		ordersAuth.GET("/my", orderHandler.NewListMyHandler(useCases.Order.ListMyOrdersUsecase))
	}

	// Public profile routes (token-based access)
	profiles := api.Group("/profiles")
	{
		profiles.GET("/validate/:token", profileHandler.NewValidateTokenHandler(useCases.Profile.ValidateTokenUsecase))
		profiles.POST("/complete", profileHandler.NewCompleteProfileHandler(useCases.Profile.CompleteProfileUsecase))
		profiles.GET("/:id", profileHandler.NewGetHandler(useCases.Profile.GetUsecase))
		profiles.PUT("/:id", profileHandler.NewUpdateHandler(useCases.Profile.UpdateUsecase))
	}

	// Protected profile routes (require auth)
	profilesAuth := api.Group("/profiles")
	profilesAuth.Use(middlewares.AuthMiddleware())
	{
		profilesAuth.POST("/generate-link", profileHandler.NewGenerateLinkHandler(useCases.Profile.GenerateLinkUsecase))
		profilesAuth.GET("/check-completed", profileHandler.NewCheckCompletedHandler(useCases.Profile.CheckCompletedUsecase))
		profilesAuth.POST("/upsert", profileHandler.NewUpsertHandler(useCases.Profile.UpsertUsecase))
	}

	// Settings routes (public for reading, could add auth for updating)
	settings := api.Group("/settings")
	{
		settings.GET("", settingsHandler.NewGetHandler(useCases.Settings.GetUsecase))
		settings.PUT("", settingsHandler.NewUpdateHandler(useCases.Settings.UpdateUsecase))
		settings.POST("/calculate-delivery", settingsHandler.NewCalculateDeliveryFeeHandler(useCases.Settings.CalculateDeliveryFeeUsecase))
	}

	// Admin routes (require auth)
	admin := api.Group("/admin")
	admin.Use(middlewares.AuthMiddleware())
	{
		admin.GET("/profiles", adminHandler.NewListProfilesHandler(useCases.Admin.ListProfilesUsecase))
		admin.GET("/orders", adminHandler.NewListOrdersHandler(useCases.Admin.ListOrdersUsecase))
		admin.GET("/transactions", adminHandler.NewListTransactionsHandler(useCases.Admin.ListTransactionsUsecase))
		admin.PUT("/orders/:id", adminHandler.NewUpdateOrderHandler(useCases.Admin.UpdateOrderUsecase))
	}

	// Payment routes (require auth)
	payment := api.Group("/payment")
	payment.Use(middlewares.AuthMiddleware())
	{
		payment.GET("/check/:user_id", paymentCheckHandler.CheckPaymentMethod)
	}

	// Payment methods proxy (require auth) — frontend talks to backend, backend adds X-API-Key
	pmHandler := paymentHandler.NewPaymentMethodsHandler(cfg)
	paymentMethods := api.Group("/payment-methods")
	paymentMethods.Use(middlewares.AuthMiddleware())
	{
		paymentMethods.GET("", pmHandler.List)
		paymentMethods.GET("/default", pmHandler.GetDefault)
		paymentMethods.POST("", pmHandler.Create)
		paymentMethods.PUT("/:id", pmHandler.Update)
		paymentMethods.DELETE("/:id", pmHandler.Delete)
	}

	// MercadoPago proxy (require auth)
	mpProxy := api.Group("/mercadopago")
	mpProxy.Use(middlewares.AuthMiddleware())
	{
		mpProxy.POST("/token", pmHandler.CreateToken)
		mpProxy.GET("/payment_method", pmHandler.GetPaymentMethod)
	}

	// WebSocket routes (auth handled via query parameter in handler)
	app.GET("/ws/notifications", wsHandler.HandleWebSocket)
}
