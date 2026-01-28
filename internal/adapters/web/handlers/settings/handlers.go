package settings

import (
	"net/http"

	apperrors "wappi/internal/platform/errors"
	"wappi/internal/platform/errors/mappings"
	settingsUsecase "wappi/internal/usecases/settings"

	"github.com/gin-gonic/gin"
)

// NewGetHandler creates a handler for getting settings
func NewGetHandler(usecase settingsUsecase.GetUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		output, appErr := usecase.Execute(c.Request.Context())
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusOK, output.Settings)
	}
}

// NewUpdateHandler creates a handler for updating settings
func NewUpdateHandler(usecase settingsUsecase.UpdateUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input settingsUsecase.UpdateInput
		if err := c.ShouldBindJSON(&input); err != nil {
			appErr := apperrors.NewApplicationError(mappings.RequestBodyParsingError, err)
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		output, appErr := usecase.Execute(c.Request.Context(), input)
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusOK, output.Settings)
	}
}

// NewCalculateDeliveryFeeHandler creates a handler for calculating delivery fee
func NewCalculateDeliveryFeeHandler(usecase settingsUsecase.CalculateDeliveryFeeUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input settingsUsecase.CalculateDeliveryFeeInput
		if err := c.ShouldBindJSON(&input); err != nil {
			appErr := apperrors.NewApplicationError(mappings.RequestBodyParsingError, err)
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		output, appErr := usecase.Execute(c.Request.Context(), input)
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusOK, output)
	}
}
