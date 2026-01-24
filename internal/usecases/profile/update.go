package profile

import (
	"context"

	"wappi/internal/domain"
	profileRepo "wappi/internal/adapters/datasources/repositories/profile"
	apperrors "wappi/internal/platform/errors"
	"wappi/internal/platform/errors/mappings"

	"github.com/google/uuid"
)

// UpdateProfileInput represents the input for updating a profile
type UpdateProfileInput struct {
	PhoneNumber string  `json:"phone_number"`
	Longitude   float64 `json:"longitude"`
	Latitude    float64 `json:"latitude"`
	Address     string  `json:"address"`
}

// UpdateProfileUsecase defines the interface for updating profiles
type UpdateProfileUsecase interface {
	Execute(ctx context.Context, id string, input UpdateProfileInput) (*GetProfileOutput, apperrors.ApplicationError)
}

type updateProfileUsecase struct {
	repo profileRepo.Repository
}

// NewUpdateProfileUsecase creates a new instance of UpdateProfileUsecase
func NewUpdateProfileUsecase(repo profileRepo.Repository) UpdateProfileUsecase {
	return &updateProfileUsecase{repo: repo}
}

// Execute updates a profile
func (u *updateProfileUsecase) Execute(ctx context.Context, id string, input UpdateProfileInput) (*GetProfileOutput, apperrors.ApplicationError) {
	// Validate profile ID
	if _, err := uuid.Parse(id); err != nil {
		return nil, apperrors.NewApplicationError(mappings.InvalidUserIDError, err)
	}

	// Get existing profile
	existingProfile, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Create or update location
	location := &domain.ProfileLocation{
		ID:        uuid.New().String(),
		Longitude: input.Longitude,
		Latitude:  input.Latitude,
		Address:   input.Address,
	}

	// If profile already has a location, we'll create a new one
	// (In a real app, you might want to update the existing location instead)
	savedLocation, locErr := u.repo.CreateLocation(ctx, location)
	if locErr != nil {
		return nil, locErr
	}

	// Update profile
	existingProfile.PhoneNumber = input.PhoneNumber
	existingProfile.LocationID = &savedLocation.ID

	updatedProfile, err := u.repo.Update(ctx, existingProfile)
	if err != nil {
		return nil, err
	}

	// Build output
	output := &GetProfileOutput{
		ID:          updatedProfile.ID,
		UserID:      updatedProfile.UserID,
		PhoneNumber: updatedProfile.PhoneNumber,
		CreatedAt:   updatedProfile.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:   updatedProfile.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}

	if updatedProfile.LocationID != nil {
		output.Location = &LocationOutput{
			ID:        savedLocation.ID,
			Longitude: savedLocation.Longitude,
			Latitude:  savedLocation.Latitude,
			Address:   savedLocation.Address,
		}
	}

	return output, nil
}
