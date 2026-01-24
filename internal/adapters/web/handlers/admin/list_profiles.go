package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	adminUsecase "wappi/internal/usecases/admin"
)

// NewListProfilesHandler creates a handler for listing all profiles
func NewListProfilesHandler(usecase adminUsecase.ListProfilesUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		output, appErr := usecase.Execute(c.Request.Context())
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusOK, output)
	}
}
