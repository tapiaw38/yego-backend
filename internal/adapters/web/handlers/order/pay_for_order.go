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

type PayForOrderRequestBody struct {
	SecurityCode string `json:"security_code" binding:"required"`
}

// NewPayForOrderHandler creates a handler for processing payment for an order
func NewPayForOrderHandler(usecase orderUsecase.PayForOrderUsecase) gin.HandlerFunc {
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

		var body PayForOrderRequestBody
		if err := c.ShouldBindJSON(&body); err != nil {
			appErr := apperrors.NewApplicationError(mappings.OrderInvalidStatusError, err)
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		authHeader := c.GetHeader("Authorization")
		var authToken string
		if authHeader != "" {
			authToken = strings.TrimPrefix(authHeader, "Bearer ")
		}

		output, appErr := usecase.Execute(c, orderUsecase.PayForOrderInput{
			OrderID:      orderID,
			UserID:       userID,
			AuthToken:    authToken,
			SecurityCode: body.SecurityCode,
		})
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusOK, output)
	}
}
