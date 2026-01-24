package order

import (
	"net/http"

	apperrors "wappi/internal/platform/errors"
	"wappi/internal/platform/errors/mappings"
	orderUsecase "wappi/internal/usecases/order"

	"github.com/gin-gonic/gin"
)

// NewCreateWithLinkHandler creates a handler for creating orders with claim links
// This endpoint is intended to be called by the WhatsApp IA
func NewCreateWithLinkHandler(usecase orderUsecase.CreateWithLinkUsecase, frontendURL string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input orderUsecase.CreateWithLinkInput
		if err := c.ShouldBindJSON(&input); err != nil {
			appErr := apperrors.NewApplicationError(mappings.RequestBodyParsingError, err)
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		output, appErr := usecase.Execute(c.Request.Context(), input, frontendURL)
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusCreated, output)
	}
}
