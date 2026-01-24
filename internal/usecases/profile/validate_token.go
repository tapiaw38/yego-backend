package profile

import (
	"context"

	profileRepo "wappi/internal/adapters/datasources/repositories/profile"
	apperrors "wappi/internal/platform/errors"
)

// ValidateTokenOutput represents the output for token validation
type ValidateTokenOutput struct {
	Valid   bool   `json:"valid"`
	UserID  string `json:"user_id"`
	Message string `json:"message,omitempty"`
}

// ValidateTokenUsecase defines the interface for validating profile tokens
type ValidateTokenUsecase interface {
	Execute(ctx context.Context, token string) (*ValidateTokenOutput, apperrors.ApplicationError)
}

type validateTokenUsecase struct {
	repo profileRepo.Repository
}

// NewValidateTokenUsecase creates a new instance of ValidateTokenUsecase
func NewValidateTokenUsecase(repo profileRepo.Repository) ValidateTokenUsecase {
	return &validateTokenUsecase{repo: repo}
}

// Execute validates a profile token
func (u *validateTokenUsecase) Execute(ctx context.Context, token string) (*ValidateTokenOutput, apperrors.ApplicationError) {
	profileToken, err := u.repo.GetToken(ctx, token)
	if err != nil {
		return nil, err
	}

	return &ValidateTokenOutput{
		Valid:  true,
		UserID: profileToken.UserID,
	}, nil
}
