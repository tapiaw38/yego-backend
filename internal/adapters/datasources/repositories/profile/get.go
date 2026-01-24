package profile

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"wappi/internal/domain"
	apperrors "wappi/internal/platform/errors"
	"wappi/internal/platform/errors/mappings"
)

// GetByUserID retrieves a profile by user ID
func (r *repository) GetByUserID(ctx context.Context, userID string) (*domain.Profile, apperrors.ApplicationError) {
	query := `
		SELECT id, user_id, phone_number, location_id, created_at, updated_at
		FROM profiles
		WHERE user_id = $1
	`

	var profile domain.Profile
	var locationID sql.NullString
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&profile.ID,
		&profile.UserID,
		&profile.PhoneNumber,
		&locationID,
		&profile.CreatedAt,
		&profile.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.NewApplicationError(mappings.ProfileNotFoundError, err)
		}
		return nil, apperrors.NewApplicationError(mappings.InternalServerError, err)
	}

	if locationID.Valid {
		profile.LocationID = &locationID.String
	}

	return &profile, nil
}

// GetByID retrieves a profile by ID
func (r *repository) GetByID(ctx context.Context, id string) (*domain.Profile, apperrors.ApplicationError) {
	query := `
		SELECT id, user_id, phone_number, location_id, created_at, updated_at
		FROM profiles
		WHERE id = $1
	`

	var profile domain.Profile
	var locationID sql.NullString
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&profile.ID,
		&profile.UserID,
		&profile.PhoneNumber,
		&locationID,
		&profile.CreatedAt,
		&profile.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.NewApplicationError(mappings.ProfileNotFoundError, err)
		}
		return nil, apperrors.NewApplicationError(mappings.InternalServerError, err)
	}

	if locationID.Valid {
		profile.LocationID = &locationID.String
	}

	return &profile, nil
}

// GetToken retrieves a profile token
func (r *repository) GetToken(ctx context.Context, token string) (*domain.ProfileToken, apperrors.ApplicationError) {
	query := `
		SELECT id, user_id, token, used, expires_at, created_at
		FROM profile_tokens
		WHERE token = $1
	`

	var profileToken domain.ProfileToken
	err := r.db.QueryRowContext(ctx, query, token).Scan(
		&profileToken.ID,
		&profileToken.UserID,
		&profileToken.Token,
		&profileToken.Used,
		&profileToken.ExpiresAt,
		&profileToken.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.NewApplicationError(mappings.ProfileTokenNotFoundError, err)
		}
		return nil, apperrors.NewApplicationError(mappings.InternalServerError, err)
	}

	// Check if token is expired
	if time.Now().After(profileToken.ExpiresAt) {
		return nil, apperrors.NewApplicationError(mappings.ProfileTokenExpiredError, nil)
	}

	// Check if token is already used
	if profileToken.Used {
		return nil, apperrors.NewApplicationError(mappings.ProfileTokenUsedError, nil)
	}

	return &profileToken, nil
}

// GetLocationByID retrieves a location by ID
func (r *repository) GetLocationByID(ctx context.Context, id string) (*domain.ProfileLocation, apperrors.ApplicationError) {
	query := `
		SELECT id, longitude, latitude, address, created_at, updated_at
		FROM profile_locations
		WHERE id = $1
	`

	var location domain.ProfileLocation
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&location.ID,
		&location.Longitude,
		&location.Latitude,
		&location.Address,
		&location.CreatedAt,
		&location.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.NewApplicationError(mappings.ProfileNotFoundError, err)
		}
		return nil, apperrors.NewApplicationError(mappings.InternalServerError, err)
	}

	return &location, nil
}

// GetAll retrieves all profiles
func (r *repository) GetAll(ctx context.Context) ([]*domain.Profile, apperrors.ApplicationError) {
	query := `
		SELECT id, user_id, phone_number, location_id, created_at, updated_at
		FROM profiles
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.InternalServerError, err)
	}
	defer rows.Close()

	var profiles []*domain.Profile
	for rows.Next() {
		var profile domain.Profile
		var locationID sql.NullString
		err := rows.Scan(
			&profile.ID,
			&profile.UserID,
			&profile.PhoneNumber,
			&locationID,
			&profile.CreatedAt,
			&profile.UpdatedAt,
		)
		if err != nil {
			return nil, apperrors.NewApplicationError(mappings.InternalServerError, err)
		}

		if locationID.Valid {
			profile.LocationID = &locationID.String
		}

		profiles = append(profiles, &profile)
	}

	if err = rows.Err(); err != nil {
		return nil, apperrors.NewApplicationError(mappings.InternalServerError, err)
	}

	return profiles, nil
}
