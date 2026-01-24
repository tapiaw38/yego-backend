package order

import (
	"net/http"

	apperrors "wappi/internal/platform/errors"
	"wappi/internal/platform/errors/mappings"
	orderUsecase "wappi/internal/usecases/order"

	"github.com/gin-gonic/gin"
)

// NewCreateHandler creates a handler for creating orders
func NewCreateHandler(usecase orderUsecase.CreateUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input orderUsecase.CreateInput
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

		c.JSON(http.StatusCreated, output)
	}
}
