package profile

import (
	"context"
	"time"

	"github.com/google/uuid"
	"wappi/internal/domain"
	apperrors "wappi/internal/platform/errors"
	"wappi/internal/platform/errors/mappings"
)

// Create inserts a new profile into the database
func (r *repository) Create(ctx context.Context, profile *domain.Profile) (*domain.Profile, apperrors.ApplicationError) {
	profile.ID = uuid.New().String()
	profile.CreatedAt = time.Now()
	profile.UpdatedAt = time.Now()

	query := `
		INSERT INTO profiles (id, user_id, phone_number, location_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.db.ExecContext(ctx, query,
		profile.ID,
		profile.UserID,
		profile.PhoneNumber,
		profile.LocationID,
		profile.CreatedAt,
		profile.UpdatedAt,
	)

	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.ProfileCreateError, err)
	}

	return profile, nil
}

// CreateLocation inserts a new location into the database
func (r *repository) CreateLocation(ctx context.Context, location *domain.ProfileLocation) (*domain.ProfileLocation, apperrors.ApplicationError) {
	location.ID = uuid.New().String()
	location.CreatedAt = time.Now()
	location.UpdatedAt = time.Now()

	query := `
		INSERT INTO profile_locations (id, longitude, latitude, address, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.db.ExecContext(ctx, query,
		location.ID,
		location.Longitude,
		location.Latitude,
		location.Address,
		location.CreatedAt,
		location.UpdatedAt,
	)

	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.LocationCreateError, err)
	}

	return location, nil
}

// CreateToken inserts a new profile token into the database
func (r *repository) CreateToken(ctx context.Context, token *domain.ProfileToken) (*domain.ProfileToken, apperrors.ApplicationError) {
	token.ID = uuid.New().String()
	token.Token = uuid.New().String()
	token.CreatedAt = time.Now()

	query := `
		INSERT INTO profile_tokens (id, user_id, token, used, expires_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.db.ExecContext(ctx, query,
		token.ID,
		token.UserID,
		token.Token,
		token.Used,
		token.ExpiresAt,
		token.CreatedAt,
	)

	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.ProfileTokenCreateError, err)
	}

	return token, nil
}
