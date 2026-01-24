package order

import (
	"net/http"

	"wappi/internal/adapters/web/middlewares"
	apperrors "wappi/internal/platform/errors"
	"wappi/internal/platform/errors/mappings"
	orderUsecase "wappi/internal/usecases/order"

	"github.com/gin-gonic/gin"
)

// NewClaimHandler creates a handler for claiming orders via token
// This endpoint requires authentication - user_id comes from JWT context
func NewClaimHandler(usecase orderUsecase.ClaimUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Param("token")
		if token == "" {
			appErr := apperrors.NewApplicationError(mappings.OrderTokenNotFoundError, nil)
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		// Get user_id from JWT context (set by auth middleware)
		userID, exists := middlewares.GetUserIDFromContext(c)
		if !exists {
			appErr := apperrors.NewApplicationError(mappings.UnauthorizedError, nil)
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		input := orderUsecase.ClaimInput{
			Token:  token,
			UserID: userID,
		}

		output, appErr := usecase.Execute(c.Request.Context(), input)
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusOK, output)
	}
}
