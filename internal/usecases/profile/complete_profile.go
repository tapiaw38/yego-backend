package profile

import (
	"context"

	profileRepo "wappi/internal/adapters/datasources/repositories/profile"
	"wappi/internal/domain"
	apperrors "wappi/internal/platform/errors"
)

// CompleteProfileInput represents the input for completing a profile
type CompleteProfileInput struct {
	Token       string  `json:"token" binding:"required"`
	PhoneNumber string  `json:"phone_number" binding:"required"`
	Longitude   float64 `json:"longitude" binding:"required"`
	Latitude    float64 `json:"latitude" binding:"required"`
	Address     string  `json:"address"`
}

// CompleteProfileOutput represents the output after completing a profile
type CompleteProfileOutput struct {
	ID          string   `json:"id"`
	UserID      string   `json:"user_id"`
	PhoneNumber string   `json:"phone_number"`
	Location    *LocationOutput `json:"location"`
	CreatedAt   string   `json:"created_at"`
	UpdatedAt   string   `json:"updated_at"`
}

// LocationOutput represents location data in the output
type LocationOutput struct {
	ID        string  `json:"id"`
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
	Address   string  `json:"address"`
}

// CompleteProfileUsecase defines the interface for completing profiles
type CompleteProfileUsecase interface {
	Execute(ctx context.Context, input CompleteProfileInput) (*CompleteProfileOutput, apperrors.ApplicationError)
}

type completeProfileUsecase struct {
	repo profileRepo.Repository
}

// NewCompleteProfileUsecase creates a new instance of CompleteProfileUsecase
func NewCompleteProfileUsecase(repo profileRepo.Repository) CompleteProfileUsecase {
	return &completeProfileUsecase{repo: repo}
}

// Execute completes a user profile
func (u *completeProfileUsecase) Execute(ctx context.Context, input CompleteProfileInput) (*CompleteProfileOutput, apperrors.ApplicationError) {
	// Validate token
	profileToken, err := u.repo.GetToken(ctx, input.Token)
	if err != nil {
		return nil, err
	}

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
	existingProfile, _ := u.repo.GetByUserID(ctx, profileToken.UserID)

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
			UserID:      profileToken.UserID,
			PhoneNumber: input.PhoneNumber,
			LocationID:  &createdLocation.ID,
		}
		resultProfile, err = u.repo.Create(ctx, newProfile)
		if err != nil {
			return nil, err
		}
	}

	// Mark token as used
	if markErr := u.repo.MarkTokenUsed(ctx, input.Token); markErr != nil {
		return nil, markErr
	}

	return &CompleteProfileOutput{
		ID:          resultProfile.ID,
		UserID:      resultProfile.UserID,
		PhoneNumber: resultProfile.PhoneNumber,
		Location: &LocationOutput{
			ID:        createdLocation.ID,
			Longitude: createdLocation.Longitude,
			Latitude:  createdLocation.Latitude,
			Address:   createdLocation.Address,
		},
		CreatedAt: resultProfile.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt: resultProfile.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}, nil
}
