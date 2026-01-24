package profile

import (
	"net/http"

	"github.com/gin-gonic/gin"
	profileUsecase "wappi/internal/usecases/profile"
)

// NewValidateTokenHandler creates a handler for validating profile tokens
func NewValidateTokenHandler(usecase profileUsecase.ValidateTokenUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Param("token")

		output, appErr := usecase.Execute(c.Request.Context(), token)
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusOK, output)
	}
}
