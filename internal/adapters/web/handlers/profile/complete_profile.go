package profile

import (
	"net/http"

	"github.com/gin-gonic/gin"
	profileUsecase "wappi/internal/usecases/profile"
	apperrors "wappi/internal/platform/errors"
	"wappi/internal/platform/errors/mappings"
)

// NewCompleteProfileHandler creates a handler for completing profiles
func NewCompleteProfileHandler(usecase profileUsecase.CompleteProfileUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input profileUsecase.CompleteProfileInput
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
