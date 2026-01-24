package ordertoken

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"wappi/internal/domain"
	apperrors "wappi/internal/platform/errors"
	"wappi/internal/platform/errors/mappings"
)

// GetByToken retrieves an order token by its token value
func (r *repository) GetByToken(ctx context.Context, token string) (*domain.OrderToken, apperrors.ApplicationError) {
	query := `
		SELECT id, order_id, token, phone_number, claimed_at, claimed_by_user_id, expires_at, created_at
		FROM order_tokens
		WHERE token = $1
	`

	var orderToken domain.OrderToken
	err := r.db.QueryRowContext(ctx, query, token).Scan(
		&orderToken.ID,
		&orderToken.OrderID,
		&orderToken.Token,
		&orderToken.PhoneNumber,
		&orderToken.ClaimedAt,
		&orderToken.ClaimedByUserID,
		&orderToken.ExpiresAt,
		&orderToken.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.NewApplicationError(mappings.OrderTokenNotFoundError, err)
		}
		return nil, apperrors.NewApplicationError(mappings.InternalServerError, err)
	}

	return &orderToken, nil
}

// GetByOrderID retrieves an order token by order ID
func (r *repository) GetByOrderID(ctx context.Context, orderID string) (*domain.OrderToken, apperrors.ApplicationError) {
	query := `
		SELECT id, order_id, token, phone_number, claimed_at, claimed_by_user_id, expires_at, created_at
		FROM order_tokens
		WHERE order_id = $1
	`

	var orderToken domain.OrderToken
	err := r.db.QueryRowContext(ctx, query, orderID).Scan(
		&orderToken.ID,
		&orderToken.OrderID,
		&orderToken.Token,
		&orderToken.PhoneNumber,
		&orderToken.ClaimedAt,
		&orderToken.ClaimedByUserID,
		&orderToken.ExpiresAt,
		&orderToken.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.NewApplicationError(mappings.OrderTokenNotFoundError, err)
		}
		return nil, apperrors.NewApplicationError(mappings.InternalServerError, err)
	}

	return &orderToken, nil
}

// MarkAsClaimed marks an order token as claimed by a user
func (r *repository) MarkAsClaimed(ctx context.Context, token string, userID string) apperrors.ApplicationError {
	now := time.Now()
	query := `
		UPDATE order_tokens
		SET claimed_at = $1, claimed_by_user_id = $2
		WHERE token = $3
	`

	result, err := r.db.ExecContext(ctx, query, now, userID, token)
	if err != nil {
		return apperrors.NewApplicationError(mappings.InternalServerError, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return apperrors.NewApplicationError(mappings.InternalServerError, err)
	}

	if rowsAffected == 0 {
		return apperrors.NewApplicationError(mappings.OrderTokenNotFoundError, errors.New("order token not found"))
	}

	return nil
}
