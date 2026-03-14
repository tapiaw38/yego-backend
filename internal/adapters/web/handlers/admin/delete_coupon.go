package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	adminUsecase "yego/internal/usecases/admin"
)

func NewDeleteCouponHandler(usecase adminUsecase.DeleteCouponUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		appErr := usecase.Execute(c, id)
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}
		c.JSON(http.StatusNoContent, nil)
	}
}
