package admin

import (
	"context"

	profileRepo "wappi/internal/adapters/datasources/repositories/profile"
	apperrors "wappi/internal/platform/errors"
)

// ProfileOutput represents a profile in the admin list
type ProfileOutput struct {
	ID          string          `json:"id"`
	UserID      string          `json:"user_id"`
	PhoneNumber string          `json:"phone_number"`
	Location    *LocationOutput `json:"location,omitempty"`
	CreatedAt   string          `json:"created_at"`
	UpdatedAt   string          `json:"updated_at"`
}

// LocationOutput represents a location in the output
type LocationOutput struct {
	ID        string  `json:"id"`
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
	Address   string  `json:"address"`
}

// ListProfilesOutput represents the output for listing profiles
type ListProfilesOutput struct {
	Profiles []ProfileOutput `json:"profiles"`
	Total    int             `json:"total"`
}

// ListProfilesUsecase defines the interface for listing profiles
type ListProfilesUsecase interface {
	Execute(ctx context.Context) (*ListProfilesOutput, apperrors.ApplicationError)
}

type listProfilesUsecase struct {
	repo profileRepo.Repository
}

// NewListProfilesUsecase creates a new instance of ListProfilesUsecase
func NewListProfilesUsecase(repo profileRepo.Repository) ListProfilesUsecase {
	return &listProfilesUsecase{repo: repo}
}

// Execute lists all profiles
func (u *listProfilesUsecase) Execute(ctx context.Context) (*ListProfilesOutput, apperrors.ApplicationError) {
	profiles, err := u.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	output := &ListProfilesOutput{
		Profiles: make([]ProfileOutput, 0, len(profiles)),
		Total:    len(profiles),
	}

	for _, p := range profiles {
		profileOutput := ProfileOutput{
			ID:          p.ID,
			UserID:      p.UserID,
			PhoneNumber: p.PhoneNumber,
			CreatedAt:   p.CreatedAt.Format("2006-01-02T15:04:05Z"),
			UpdatedAt:   p.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		}

		// Get location if exists
		if p.LocationID != nil {
			location, locErr := u.repo.GetLocationByID(ctx, *p.LocationID)
			if locErr == nil && location != nil {
				profileOutput.Location = &LocationOutput{
					ID:        location.ID,
					Longitude: location.Longitude,
					Latitude:  location.Latitude,
					Address:   location.Address,
				}
			}
		}

		output.Profiles = append(output.Profiles, profileOutput)
	}

	return output, nil
}
