package profile

import (
	"context"

	profileRepo "wappi/internal/adapters/datasources/repositories/profile"
	apperrors "wappi/internal/platform/errors"
	"wappi/internal/platform/errors/mappings"

	"github.com/google/uuid"
)

// GetProfileOutput represents the output for getting a profile
type GetProfileOutput struct {
	ID          string          `json:"id"`
	UserID      string          `json:"user_id"`
	PhoneNumber string          `json:"phone_number"`
	Location    *LocationOutput `json:"location,omitempty"`
	CreatedAt   string          `json:"created_at"`
	UpdatedAt   string          `json:"updated_at"`
}

// GetProfileUsecase defines the interface for getting profiles
type GetProfileUsecase interface {
	Execute(ctx context.Context, id string) (*GetProfileOutput, apperrors.ApplicationError)
}

type getProfileUsecase struct {
	repo profileRepo.Repository
}

// NewGetProfileUsecase creates a new instance of GetProfileUsecase
func NewGetProfileUsecase(repo profileRepo.Repository) GetProfileUsecase {
	return &getProfileUsecase{repo: repo}
}

// Execute retrieves a profile by ID
func (u *getProfileUsecase) Execute(ctx context.Context, id string) (*GetProfileOutput, apperrors.ApplicationError) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, apperrors.NewApplicationError(mappings.InvalidUserIDError, err)
	}

	profile, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	output := &GetProfileOutput{
		ID:          profile.ID,
		UserID:      profile.UserID,
		PhoneNumber: profile.PhoneNumber,
		CreatedAt:   profile.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:   profile.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}

	// Get location if exists
	if profile.LocationID != nil {
		location, locErr := u.repo.GetLocationByID(ctx, *profile.LocationID)
		if locErr == nil && location != nil {
			output.Location = &LocationOutput{
				ID:        location.ID,
				Longitude: location.Longitude,
				Latitude:  location.Latitude,
				Address:   location.Address,
			}
		}
	}

	return output, nil
}
