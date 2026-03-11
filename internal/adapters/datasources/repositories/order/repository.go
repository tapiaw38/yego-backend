package order

import (
	"context"
	"database/sql"

	"yego/internal/domain"
	apperrors "yego/internal/platform/errors"
)

// Repository defines the interface for order data operations
type Repository interface {
	Create(ctx context.Context, order *domain.Order) (*domain.Order, apperrors.ApplicationError)
	GetByID(ctx context.Context, id string) (*domain.Order, apperrors.ApplicationError)
	GetAll(ctx context.Context) ([]*domain.Order, apperrors.ApplicationError)
	GetByUserID(ctx context.Context, userID string) ([]*domain.Order, apperrors.ApplicationError)
	GetByDeliveryUserID(ctx context.Context, deliveryUserID string) ([]*DeliveryOrderRow, apperrors.ApplicationError)
	UpdateStatus(ctx context.Context, id string, status domain.OrderStatus) (*domain.Order, apperrors.ApplicationError)
	Update(ctx context.Context, order *domain.Order) (*domain.Order, apperrors.ApplicationError)
	AssignUser(ctx context.Context, orderID string, userID string) apperrors.ApplicationError
	AssignProfile(ctx context.Context, orderID string, profileID string) apperrors.ApplicationError
	AssignDelivery(ctx context.Context, orderID string, deliveryUserID string) (*domain.Order, apperrors.ApplicationError)
	AcceptDelivery(ctx context.Context, orderID string, deliveryUserID string) (*domain.Order, apperrors.ApplicationError)
}

type repository struct {
	db *sql.DB
}

// NewRepository creates a new order repository
func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}
