package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"yego/internal/usecases/admin"
)

func NewDeleteUploadHandler(usecase admin.DeleteUploadUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.Query("key")
		if key == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "key is required"})
			return
		}
		if err := usecase.Execute(c.Request.Context(), key); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.Status(http.StatusNoContent)
	}
}
