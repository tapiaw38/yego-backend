package order

import (
	"net/http"

	"github.com/gin-gonic/gin"
	orderUsecase "wappi/internal/usecases/order"
)

// NewGetHandler creates a handler for getting orders
func NewGetHandler(usecase orderUsecase.GetUsecase) gin.HandlerFunc {
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
