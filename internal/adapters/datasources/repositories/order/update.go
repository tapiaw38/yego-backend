package order

import (
	"context"
	"time"

	"wappi/internal/domain"
	apperrors "wappi/internal/platform/errors"
	"wappi/internal/platform/errors/mappings"
)

// UpdateStatus updates the status of an order
func (r *repository) UpdateStatus(ctx context.Context, id string, status domain.OrderStatus) (*domain.Order, apperrors.ApplicationError) {
	query := `
		UPDATE orders
		SET status = $1, updated_at = $2
		WHERE id = $3
	`

	result, err := r.db.ExecContext(ctx, query, status, time.Now(), id)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.OrderUpdateError, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.OrderUpdateError, err)
	}

	if rowsAffected == 0 {
		return nil, apperrors.NewApplicationError(mappings.OrderNotFoundError, nil)
	}

	return r.GetByID(ctx, id)
}

// Update updates an order (status, eta, etc.)
func (r *repository) Update(ctx context.Context, order *domain.Order) (*domain.Order, apperrors.ApplicationError) {
	order.UpdatedAt = time.Now()

	query := `
		UPDATE orders
		SET status = $1, eta = $2, updated_at = $3
		WHERE id = $4
	`

	result, err := r.db.ExecContext(ctx, query, order.Status, order.ETA, order.UpdatedAt, order.ID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.OrderUpdateError, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.OrderUpdateError, err)
	}

	if rowsAffected == 0 {
		return nil, apperrors.NewApplicationError(mappings.OrderNotFoundError, nil)
	}

	return r.GetByID(ctx, order.ID)
}
