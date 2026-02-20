package order

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"wappi/internal/adapters/web/middlewares"
	apperrors "wappi/internal/platform/errors"
	"wappi/internal/platform/errors/mappings"
	orderUsecase "wappi/internal/usecases/order"
)

// NewCreatePaymentLinkHandler creates a handler for generating a MercadoPago Checkout Pro payment link
func NewCreatePaymentLinkHandler(usecase orderUsecase.CreatePaymentLinkUsecase, frontendURL string) gin.HandlerFunc {
	return func(c *gin.Context) {
		orderID := c.Param("id")
		if orderID == "" {
			appErr := apperrors.NewApplicationError(mappings.OrderNotFoundError, nil)
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		userID, exists := middlewares.GetUserIDFromContext(c)
		if !exists {
			appErr := apperrors.NewApplicationError(mappings.UnauthorizedError, nil)
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		authHeader := c.GetHeader("Authorization")
		var authToken string
		if authHeader != "" {
			authToken = strings.TrimPrefix(authHeader, "Bearer ")
		}

		output, appErr := usecase.Execute(c, orderUsecase.CreatePaymentLinkInput{
			OrderID:     orderID,
			UserID:      userID,
			AuthToken:   authToken,
			FrontendURL: frontendURL,
		})
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusOK, output)
	}
}
