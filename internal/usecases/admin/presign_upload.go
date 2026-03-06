package admin

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	s3service "yego/internal/services/s3"
)

type (
	PresignUploadUsecase interface {
		Execute(ctx context.Context, input PresignUploadInput) (*PresignUploadOutput, error)
	}

	presignUploadUsecase struct {
		s3 *s3service.Client
	}

	PresignUploadInput struct {
		Filename    string
		ContentType string
		Folder      string // e.g. "imports/images"
	}

	PresignUploadOutput struct {
		UploadURL string `json:"upload_url"`
		PublicURL string `json:"public_url"`
		Key       string `json:"key"`
	}
)

func NewPresignUploadUsecase(s3Client *s3service.Client) PresignUploadUsecase {
	return &presignUploadUsecase{s3: s3Client}
}

func (u *presignUploadUsecase) Execute(_ context.Context, input PresignUploadInput) (*PresignUploadOutput, error) {
	if !u.s3.IsConfigured() {
		return nil, fmt.Errorf("S3 not configured")
	}

	ext := filepath.Ext(input.Filename)
	if ext == "" {
		ext = extensionFromContentType(input.ContentType)
	}

	folder := input.Folder
	if folder == "" {
		folder = "imports/images"
	}

	key := fmt.Sprintf("%s/%s%s", folder, uuid.New().String(), ext)

	uploadURL, publicURL, err := u.s3.PresignPut(key, 15*time.Minute)
	if err != nil {
		return nil, fmt.Errorf("presign error: %w", err)
	}

	return &PresignUploadOutput{
		UploadURL: uploadURL,
		PublicURL: publicURL,
		Key:       key,
	}, nil
}

func extensionFromContentType(ct string) string {
	ct = strings.ToLower(strings.Split(ct, ";")[0])
	switch ct {
	case "image/jpeg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "image/webp":
		return ".webp"
	case "image/gif":
		return ".gif"
	default:
		return ""
	}
}
