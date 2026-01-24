package profile

import (
	"net/http"

	"wappi/internal/adapters/web/middlewares"
	apperrors "wappi/internal/platform/errors"
	"wappi/internal/platform/errors/mappings"
	profileUsecase "wappi/internal/usecases/profile"

	"github.com/gin-gonic/gin"
)

func NewCreateOrUpdateHandler(usecase profileUsecase.CreateOrUpdateProfileUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := middlewares.GetUserIDFromContext(c)
		if !exists || userID == "" {
			appErr := apperrors.NewApplicationError(mappings.UnauthorizedError, nil)
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		var input profileUsecase.CreateOrUpdateProfileInput
		if err := c.ShouldBindJSON(&input); err != nil {
			appErr := apperrors.NewApplicationError(mappings.RequestBodyParsingError, err)
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		output, appErr := usecase.Execute(c.Request.Context(), userID, input)
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusOK, output)
	}
}
