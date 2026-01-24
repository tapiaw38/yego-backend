package profile

import (
	"context"
	"database/sql"

	"wappi/internal/domain"
	apperrors "wappi/internal/platform/errors"
)

// Repository defines the interface for profile data operations
type Repository interface {
	Create(ctx context.Context, profile *domain.Profile) (*domain.Profile, apperrors.ApplicationError)
	GetByUserID(ctx context.Context, userID string) (*domain.Profile, apperrors.ApplicationError)
	GetByID(ctx context.Context, id string) (*domain.Profile, apperrors.ApplicationError)
	GetAll(ctx context.Context) ([]*domain.Profile, apperrors.ApplicationError)
	Update(ctx context.Context, profile *domain.Profile) (*domain.Profile, apperrors.ApplicationError)
	CreateToken(ctx context.Context, token *domain.ProfileToken) (*domain.ProfileToken, apperrors.ApplicationError)
	GetToken(ctx context.Context, token string) (*domain.ProfileToken, apperrors.ApplicationError)
	MarkTokenUsed(ctx context.Context, token string) apperrors.ApplicationError
	CreateLocation(ctx context.Context, location *domain.ProfileLocation) (*domain.ProfileLocation, apperrors.ApplicationError)
	GetLocationByID(ctx context.Context, id string) (*domain.ProfileLocation, apperrors.ApplicationError)
}

type repository struct {
	db *sql.DB
}

// NewRepository creates a new profile repository
func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}
