package profile

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"wappi/internal/adapters/web/middlewares"
	apperrors "wappi/internal/platform/errors"
	"wappi/internal/platform/errors/mappings"
	profileUsecase "wappi/internal/usecases/profile"
)

// NewGenerateLinkHandler creates a handler for generating profile links
// User ID is extracted from JWT token (set by AuthMiddleware)
func NewGenerateLinkHandler(usecase profileUsecase.GenerateLinkUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user ID from JWT context (set by AuthMiddleware)
		userID, exists := middlewares.GetUserIDFromContext(c)
		if !exists || userID == "" {
			appErr := apperrors.NewApplicationError(mappings.UnauthorizedError, nil)
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		input := profileUsecase.GenerateLinkInput{
			UserID: userID,
		}

		output, appErr := usecase.Execute(c.Request.Context(), input)
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusCreated, output)
	}
}
