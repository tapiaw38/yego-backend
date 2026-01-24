package ordertoken

import (
	"context"
	"database/sql"

	"wappi/internal/domain"
	apperrors "wappi/internal/platform/errors"
)

// Repository defines the interface for order token data operations
type Repository interface {
	Create(ctx context.Context, token *domain.OrderToken) (*domain.OrderToken, apperrors.ApplicationError)
	GetByToken(ctx context.Context, token string) (*domain.OrderToken, apperrors.ApplicationError)
	GetByOrderID(ctx context.Context, orderID string) (*domain.OrderToken, apperrors.ApplicationError)
	MarkAsClaimed(ctx context.Context, token string, userID string) apperrors.ApplicationError
}

type repository struct {
	db *sql.DB
}

// NewRepository creates a new order token repository
func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}
