package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	apperrors "yego/internal/platform/errors"
	"yego/internal/platform/errors/mappings"
	adminUsecase "yego/internal/usecases/admin"
)

type AssignDeliveryInput struct {
	DeliveryUserID string `json:"delivery_user_id"`
}

func NewAssignDeliveryHandler(usecase adminUsecase.AssignDeliveryUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		orderID := c.Param("id")

		var input AssignDeliveryInput
		if err := c.ShouldBindJSON(&input); err != nil {
			appErr := apperrors.NewApplicationError(mappings.RequestBodyParsingError, err)
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		output, appErr := usecase.Execute(c, orderID, adminUsecase.AssignDeliveryInput{
			DeliveryUserID: input.DeliveryUserID,
		})
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusOK, output)
	}
}
