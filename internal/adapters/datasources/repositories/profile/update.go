package profile

import (
	"context"
	"time"

	"wappi/internal/domain"
	apperrors "wappi/internal/platform/errors"
	"wappi/internal/platform/errors/mappings"
)

// Update updates a profile in the database
func (r *repository) Update(ctx context.Context, profile *domain.Profile) (*domain.Profile, apperrors.ApplicationError) {
	profile.UpdatedAt = time.Now()

	query := `
		UPDATE profiles
		SET phone_number = $1, location_id = $2, updated_at = $3
		WHERE id = $4
	`

	result, err := r.db.ExecContext(ctx, query,
		profile.PhoneNumber,
		profile.LocationID,
		profile.UpdatedAt,
		profile.ID,
	)

	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.ProfileUpdateError, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.ProfileUpdateError, err)
	}

	if rowsAffected == 0 {
		return nil, apperrors.NewApplicationError(mappings.ProfileNotFoundError, nil)
	}

	return r.GetByID(ctx, profile.ID)
}

// MarkTokenUsed marks a token as used
func (r *repository) MarkTokenUsed(ctx context.Context, token string) apperrors.ApplicationError {
	query := `
		UPDATE profile_tokens
		SET used = TRUE
		WHERE token = $1
	`

	_, err := r.db.ExecContext(ctx, query, token)
	if err != nil {
		return apperrors.NewApplicationError(mappings.ProfileUpdateError, err)
	}

	return nil
}
