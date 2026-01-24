package profile

import (
	"net/http"

	"github.com/gin-gonic/gin"
	profileUsecase "wappi/internal/usecases/profile"
)

// NewUpdateHandler creates a handler for updating a profile
func NewUpdateHandler(usecase profileUsecase.UpdateProfileUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		var input profileUsecase.UpdateProfileInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		output, appErr := usecase.Execute(c.Request.Context(), id, input)
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusOK, output)
	}
}
