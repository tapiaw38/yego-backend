package order

import (
	"context"
	"time"

	"github.com/google/uuid"
	"wappi/internal/domain"
	apperrors "wappi/internal/platform/errors"
	"wappi/internal/platform/errors/mappings"
)

// Create inserts a new order into the database
func (r *repository) Create(ctx context.Context, order *domain.Order) (*domain.Order, apperrors.ApplicationError) {
	order.ID = uuid.New().String()
	order.Status = domain.StatusCreated
	order.CreatedAt = time.Now()
	order.UpdatedAt = time.Now()

	dataJSON, err := order.DataJSON()
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.OrderCreateError, err)
	}

	query := `
		INSERT INTO orders (id, profile_id, user_id, status, eta, data, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err = r.db.ExecContext(ctx, query,
		order.ID,
		order.ProfileID,
		order.UserID,
		order.Status,
		order.ETA,
		dataJSON,
		order.CreatedAt,
		order.UpdatedAt,
	)

	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.OrderCreateError, err)
	}

	return order, nil
}
