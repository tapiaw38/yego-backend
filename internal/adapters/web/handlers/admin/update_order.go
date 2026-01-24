package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	adminUsecase "wappi/internal/usecases/admin"
)

// NewUpdateOrderHandler creates a handler for updating an order
func NewUpdateOrderHandler(usecase adminUsecase.UpdateOrderUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		var input adminUsecase.UpdateOrderInput
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
