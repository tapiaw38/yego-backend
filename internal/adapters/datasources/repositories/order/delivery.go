package order

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"yego/internal/domain"
	apperrors "yego/internal/platform/errors"
	"yego/internal/platform/errors/mappings"
)

// DeliveryOrderRow is an order enriched with profile/location data for delivery views.
type DeliveryOrderRow struct {
	Order       domain.Order
	PhoneNumber string
	Address     sql.NullString
	Latitude    sql.NullFloat64
	Longitude   sql.NullFloat64
}

// GetByDeliveryUserID retrieves all orders assigned to a delivery user,
// enriched with the customer's phone number and location.
func (r *repository) GetByDeliveryUserID(ctx context.Context, deliveryUserID string) ([]*DeliveryOrderRow, apperrors.ApplicationError) {
	query := `
		SELECT
			o.id, o.profile_id, o.user_id, o.status, o.status_message, o.eta, o.data,
			o.delivery_user_id, o.delivery_accepted_at, o.created_at, o.updated_at,
			COALESCE(p.phone_number, '') AS phone_number,
			pl.address,
			pl.latitude,
			pl.longitude
		FROM orders o
		LEFT JOIN profiles p ON p.id = o.profile_id
		LEFT JOIN profile_locations pl ON pl.id = p.location_id
		WHERE o.delivery_user_id = $1
		ORDER BY o.created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, deliveryUserID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.InternalServerError, err)
	}
	defer rows.Close()

	var result []*DeliveryOrderRow
	for rows.Next() {
		var row DeliveryOrderRow
		var dataJSON []byte
		var statusMessage sql.NullString
		var deliveryUID sql.NullString
		var deliveryAcceptedAt sql.NullTime

		if err := rows.Scan(
			&row.Order.ID,
			&row.Order.ProfileID,
			&row.Order.UserID,
			&row.Order.Status,
			&statusMessage,
			&row.Order.ETA,
			&dataJSON,
			&deliveryUID,
			&deliveryAcceptedAt,
			&row.Order.CreatedAt,
			&row.Order.UpdatedAt,
			&row.PhoneNumber,
			&row.Address,
			&row.Latitude,
			&row.Longitude,
		); err != nil {
			return nil, apperrors.NewApplicationError(mappings.InternalServerError, err)
		}

		if dataJSON != nil {
			_ = row.Order.SetDataFromJSON(dataJSON)
		}
		if statusMessage.Valid {
			row.Order.StatusMessage = &statusMessage.String
		}
		if deliveryUID.Valid {
			row.Order.DeliveryUserID = &deliveryUID.String
		}
		if deliveryAcceptedAt.Valid {
			row.Order.DeliveryAcceptedAt = &deliveryAcceptedAt.Time
		}

		result = append(result, &row)
	}

	if err = rows.Err(); err != nil {
		return nil, apperrors.NewApplicationError(mappings.InternalServerError, err)
	}

	return result, nil
}

// AssignDelivery sets the delivery_user_id for an order (called by admin).
func (r *repository) AssignDelivery(ctx context.Context, orderID string, deliveryUserID string) (*domain.Order, apperrors.ApplicationError) {
	query := `
		UPDATE orders
		SET delivery_user_id = $1, updated_at = $2
		WHERE id = $3
	`

	result, err := r.db.ExecContext(ctx, query, deliveryUserID, time.Now(), orderID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.InternalServerError, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.InternalServerError, err)
	}

	if rowsAffected == 0 {
		return nil, apperrors.NewApplicationError(mappings.OrderNotFoundError, errors.New("order not found"))
	}

	return r.GetByID(ctx, orderID)
}

// AcceptDelivery sets delivery_accepted_at for an order (called by delivery user).
// Only the assigned delivery user can accept.
func (r *repository) AcceptDelivery(ctx context.Context, orderID string, deliveryUserID string) (*domain.Order, apperrors.ApplicationError) {
	query := `
		UPDATE orders
		SET delivery_accepted_at = $1, updated_at = $1
		WHERE id = $2 AND delivery_user_id = $3 AND delivery_accepted_at IS NULL
	`

	now := time.Now()
	result, err := r.db.ExecContext(ctx, query, now, orderID, deliveryUserID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.InternalServerError, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.InternalServerError, err)
	}

	if rowsAffected == 0 {
		return nil, apperrors.NewApplicationError(mappings.OrderNotFoundError,
			errors.New("order not found, already accepted, or not assigned to you"))
	}

	return r.GetByID(ctx, orderID)
}
