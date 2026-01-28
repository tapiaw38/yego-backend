package order

import (
	"context"
	"database/sql"
	"errors"

	"wappi/internal/domain"
	apperrors "wappi/internal/platform/errors"
	"wappi/internal/platform/errors/mappings"
)

// GetByID retrieves an order by its ID
func (r *repository) GetByID(ctx context.Context, id string) (*domain.Order, apperrors.ApplicationError) {
	query := `
		SELECT id, profile_id, user_id, status, eta, data, created_at, updated_at
		FROM orders
		WHERE id = $1
	`

	var order domain.Order
	var dataJSON []byte
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&order.ID,
		&order.ProfileID,
		&order.UserID,
		&order.Status,
		&order.ETA,
		&dataJSON,
		&order.CreatedAt,
		&order.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.NewApplicationError(mappings.OrderNotFoundError, err)
		}
		return nil, apperrors.NewApplicationError(mappings.InternalServerError, err)
	}

	if dataJSON != nil {
		if err := order.SetDataFromJSON(dataJSON); err != nil {
			return nil, apperrors.NewApplicationError(mappings.InternalServerError, err)
		}
	}

	return &order, nil
}

// GetAll retrieves all orders
func (r *repository) GetAll(ctx context.Context) ([]*domain.Order, apperrors.ApplicationError) {
	query := `
		SELECT id, profile_id, user_id, status, eta, data, created_at, updated_at
		FROM orders
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.InternalServerError, err)
	}
	defer rows.Close()

	var orders []*domain.Order
	for rows.Next() {
		var order domain.Order
		var dataJSON []byte
		err := rows.Scan(
			&order.ID,
			&order.ProfileID,
			&order.UserID,
			&order.Status,
			&order.ETA,
			&dataJSON,
			&order.CreatedAt,
			&order.UpdatedAt,
		)
		if err != nil {
			return nil, apperrors.NewApplicationError(mappings.InternalServerError, err)
		}
		if dataJSON != nil {
			_ = order.SetDataFromJSON(dataJSON)
		}
		orders = append(orders, &order)
	}

	if err = rows.Err(); err != nil {
		return nil, apperrors.NewApplicationError(mappings.InternalServerError, err)
	}

	return orders, nil
}

// GetByUserID retrieves all orders for a specific user
func (r *repository) GetByUserID(ctx context.Context, userID string) ([]*domain.Order, apperrors.ApplicationError) {
	query := `
		SELECT id, profile_id, user_id, status, eta, data, created_at, updated_at
		FROM orders
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.InternalServerError, err)
	}
	defer rows.Close()

	var orders []*domain.Order
	for rows.Next() {
		var order domain.Order
		var dataJSON []byte
		err := rows.Scan(
			&order.ID,
			&order.ProfileID,
			&order.UserID,
			&order.Status,
			&order.ETA,
			&dataJSON,
			&order.CreatedAt,
			&order.UpdatedAt,
		)
		if err != nil {
			return nil, apperrors.NewApplicationError(mappings.InternalServerError, err)
		}
		if dataJSON != nil {
			_ = order.SetDataFromJSON(dataJSON)
		}
		orders = append(orders, &order)
	}

	if err = rows.Err(); err != nil {
		return nil, apperrors.NewApplicationError(mappings.InternalServerError, err)
	}

	return orders, nil
}

// AssignUser assigns a user to an order
func (r *repository) AssignUser(ctx context.Context, orderID string, userID string) apperrors.ApplicationError {
	query := `
		UPDATE orders
		SET user_id = $1, updated_at = NOW()
		WHERE id = $2
	`

	result, err := r.db.ExecContext(ctx, query, userID, orderID)
	if err != nil {
		return apperrors.NewApplicationError(mappings.InternalServerError, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return apperrors.NewApplicationError(mappings.InternalServerError, err)
	}

	if rowsAffected == 0 {
		return apperrors.NewApplicationError(mappings.OrderNotFoundError, errors.New("order not found"))
	}

	return nil
}

// AssignProfile assigns a profile to an order
func (r *repository) AssignProfile(ctx context.Context, orderID string, profileID string) apperrors.ApplicationError {
	query := `
		UPDATE orders
		SET profile_id = $1, updated_at = NOW()
		WHERE id = $2
	`

	result, err := r.db.ExecContext(ctx, query, profileID, orderID)
	if err != nil {
		return apperrors.NewApplicationError(mappings.InternalServerError, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return apperrors.NewApplicationError(mappings.InternalServerError, err)
	}

	if rowsAffected == 0 {
		return apperrors.NewApplicationError(mappings.OrderNotFoundError, errors.New("order not found"))
	}

	return nil
}
