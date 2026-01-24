package profile

import (
	"context"

	profileRepo "wappi/internal/adapters/datasources/repositories/profile"
	apperrors "wappi/internal/platform/errors"
)

// CheckCompletedOutput represents the output for checking profile completion
type CheckCompletedOutput struct {
	IsCompleted bool    `json:"is_completed"`
	ProfileID   *string `json:"profile_id,omitempty"`
	Message     string  `json:"message"`
}

// CheckCompletedUsecase defines the interface for checking profile completion
type CheckCompletedUsecase interface {
	Execute(ctx context.Context, userID string) (*CheckCompletedOutput, apperrors.ApplicationError)
}

type checkCompletedUsecase struct {
	repo profileRepo.Repository
}

// NewCheckCompletedUsecase creates a new instance of CheckCompletedUsecase
func NewCheckCompletedUsecase(repo profileRepo.Repository) CheckCompletedUsecase {
	return &checkCompletedUsecase{repo: repo}
}

// Execute checks if the user's profile is completed
func (u *checkCompletedUsecase) Execute(ctx context.Context, userID string) (*CheckCompletedOutput, apperrors.ApplicationError) {
	// Try to get profile by user ID
	profile, err := u.repo.GetByUserID(ctx, userID)
	if err != nil {
		// Profile doesn't exist - not completed
		return &CheckCompletedOutput{
			IsCompleted: false,
			ProfileID:   nil,
			Message:     "El perfil no existe. Por favor complete sus datos.",
		}, nil
	}

	// Check if profile is completed using domain method
	isCompleted := profile.IsCompleted()

	var message string
	if isCompleted {
		message = "Perfil completo"
	} else {
		if profile.PhoneNumber == "" {
			message = "Falta completar el número de teléfono"
		} else if profile.LocationID == nil || *profile.LocationID == "" {
			message = "Falta completar la dirección de entrega"
		} else {
			message = "Perfil incompleto"
		}
	}

	return &CheckCompletedOutput{
		IsCompleted: isCompleted,
		ProfileID:   &profile.ID,
		Message:     message,
	}, nil
}
