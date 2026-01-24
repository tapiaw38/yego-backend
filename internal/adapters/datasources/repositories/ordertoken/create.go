package ordertoken

import (
	"context"
	"time"

	"github.com/google/uuid"
	"wappi/internal/domain"
	apperrors "wappi/internal/platform/errors"
	"wappi/internal/platform/errors/mappings"
)

// Create inserts a new order token into the database
func (r *repository) Create(ctx context.Context, token *domain.OrderToken) (*domain.OrderToken, apperrors.ApplicationError) {
	token.ID = uuid.New().String()
	token.Token = uuid.New().String()
	token.CreatedAt = time.Now()

	query := `
		INSERT INTO order_tokens (id, order_id, token, phone_number, expires_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.db.ExecContext(ctx, query,
		token.ID,
		token.OrderID,
		token.Token,
		token.PhoneNumber,
		token.ExpiresAt,
		token.CreatedAt,
	)

	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.OrderTokenCreateError, err)
	}

	return token, nil
}
