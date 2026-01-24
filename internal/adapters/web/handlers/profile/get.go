package profile

import (
	"net/http"

	"github.com/gin-gonic/gin"
	profileUsecase "wappi/internal/usecases/profile"
)

// NewGetHandler creates a handler for getting a profile by ID
func NewGetHandler(usecase profileUsecase.GetProfileUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		output, appErr := usecase.Execute(c.Request.Context(), id)
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusOK, output)
	}
}
