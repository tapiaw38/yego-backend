package profile

import (
	"context"

	profileRepo "wappi/internal/adapters/datasources/repositories/profile"
	"wappi/internal/domain"
	apperrors "wappi/internal/platform/errors"
)

// CreateOrUpdateProfileInput represents the input for creating or updating a profile
type CreateOrUpdateProfileInput struct {
	PhoneNumber string  `json:"phone_number" binding:"required"`
	Longitude   float64 `json:"longitude" binding:"required"`
	Latitude    float64 `json:"latitude" binding:"required"`
	Address     string  `json:"address"`
}

// CreateOrUpdateProfileUsecase defines the interface for creating or updating profiles
type CreateOrUpdateProfileUsecase interface {
	Execute(ctx context.Context, userID string, input CreateOrUpdateProfileInput) (*GetProfileOutput, apperrors.ApplicationError)
}

type createOrUpdateProfileUsecase struct {
	repo profileRepo.Repository
}

// NewCreateOrUpdateProfileUsecase creates a new instance of CreateOrUpdateProfileUsecase
func NewCreateOrUpdateProfileUsecase(repo profileRepo.Repository) CreateOrUpdateProfileUsecase {
	return &createOrUpdateProfileUsecase{repo: repo}
}

// Execute creates or updates a profile for the authenticated user
func (u *createOrUpdateProfileUsecase) Execute(ctx context.Context, userID string, input CreateOrUpdateProfileInput) (*GetProfileOutput, apperrors.ApplicationError) {
	// Create location first
	location := &domain.ProfileLocation{
		Longitude: input.Longitude,
		Latitude:  input.Latitude,
		Address:   input.Address,
	}

	createdLocation, err := u.repo.CreateLocation(ctx, location)
	if err != nil {
		return nil, err
	}

	// Check if profile already exists for this user
	existingProfile, _ := u.repo.GetByUserID(ctx, userID)

	var resultProfile *domain.Profile

	if existingProfile != nil {
		// Update existing profile
		existingProfile.PhoneNumber = input.PhoneNumber
		existingProfile.LocationID = &createdLocation.ID
		resultProfile, err = u.repo.Update(ctx, existingProfile)
		if err != nil {
			return nil, err
		}
	} else {
		// Create new profile
		newProfile := &domain.Profile{
			UserID:      userID,
			PhoneNumber: input.PhoneNumber,
			LocationID:  &createdLocation.ID,
		}
		resultProfile, err = u.repo.Create(ctx, newProfile)
		if err != nil {
			return nil, err
		}
	}

	// Build output
	output := &GetProfileOutput{
		ID:          resultProfile.ID,
		UserID:      resultProfile.UserID,
		PhoneNumber: resultProfile.PhoneNumber,
		CreatedAt:   resultProfile.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:   resultProfile.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}

	if resultProfile.LocationID != nil {
		output.Location = &LocationOutput{
			ID:        createdLocation.ID,
			Longitude: createdLocation.Longitude,
			Latitude:  createdLocation.Latitude,
			Address:   createdLocation.Address,
		}
	}

	return output, nil
}

