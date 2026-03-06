package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"yego/internal/usecases/admin"
)

type presignUploadRequest struct {
	Filename    string `json:"filename" binding:"required"`
	ContentType string `json:"content_type"`
	Folder      string `json:"folder"`
}

func NewPresignUploadHandler(usecase admin.PresignUploadUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req presignUploadRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		output, err := usecase.Execute(c.Request.Context(), admin.PresignUploadInput{
			Filename:    req.Filename,
			ContentType: req.ContentType,
			Folder:      req.Folder,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, output)
	}
}
