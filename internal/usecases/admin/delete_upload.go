package admin

import (
	"context"
	"fmt"
	"strings"

	s3service "yego/internal/services/s3"
)

type (
	DeleteUploadUsecase interface {
		Execute(ctx context.Context, key string) error
	}

	deleteUploadUsecase struct {
		s3 *s3service.Client
	}
)

func NewDeleteUploadUsecase(s3Client *s3service.Client) DeleteUploadUsecase {
	return &deleteUploadUsecase{s3: s3Client}
}

func (u *deleteUploadUsecase) Execute(_ context.Context, key string) error {
	if !u.s3.IsConfigured() {
		return fmt.Errorf("S3 not configured")
	}
	// Accept either a key ("imports/images/abc.jpg") or a full URL
	if strings.HasPrefix(key, "https://") {
		// Extract key from URL: https://{bucket}.s3.{region}.amazonaws.com/{key}
		parts := strings.SplitN(key, ".amazonaws.com/", 2)
		if len(parts) == 2 {
			key = parts[1]
			// Strip query string if any
			if idx := strings.Index(key, "?"); idx != -1 {
				key = key[:idx]
			}
		}
	}
	return u.s3.DeleteObject(key)
}
