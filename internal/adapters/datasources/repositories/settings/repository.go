package settings

import (
	"context"
	"database/sql"
	"time"

	"wappi/internal/domain"
	apperrors "wappi/internal/platform/errors"
	"wappi/internal/platform/errors/mappings"

	"github.com/google/uuid"
)

// Repository defines the interface for settings operations
type Repository interface {
	Get(ctx context.Context) (*domain.Settings, apperrors.ApplicationError)
	Upsert(ctx context.Context, settings *domain.Settings) (*domain.Settings, apperrors.ApplicationError)
}

type repository struct {
	db *sql.DB
}

// NewRepository creates a new settings repository
func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

// Get retrieves the settings (there should only be one row)
func (r *repository) Get(ctx context.Context) (*domain.Settings, apperrors.ApplicationError) {
	query := `
		SELECT id, business_name, business_latitude, business_longitude,
			   default_map_latitude, default_map_longitude, default_map_zoom,
			   default_item_weight, delivery_base_price, delivery_price_per_km,
			   delivery_price_per_kg, created_at, updated_at
		FROM settings
		LIMIT 1
	`

	var s domain.Settings
	err := r.db.QueryRowContext(ctx, query).Scan(
		&s.ID, &s.BusinessName, &s.BusinessLatitude, &s.BusinessLongitude,
		&s.DefaultMapLatitude, &s.DefaultMapLongitude, &s.DefaultMapZoom,
		&s.DefaultItemWeight, &s.DeliveryBasePrice, &s.DeliveryPricePerKm,
		&s.DeliveryPricePerKg, &s.CreatedAt, &s.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		// Return default settings if none exist
		return &domain.Settings{
			ID:                  "",
			BusinessName:        "",
			BusinessLatitude:    -34.6037, // Buenos Aires default
			BusinessLongitude:   -58.3816,
			DefaultMapLatitude:  -34.6037,
			DefaultMapLongitude: -58.3816,
			DefaultMapZoom:      13,
			DefaultItemWeight:   500, // 500g default
			DeliveryBasePrice:   500,
			DeliveryPricePerKm:  200,
			DeliveryPricePerKg:  100,
		}, nil
	}

	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.SettingsGetError, err)
	}

	return &s, nil
}

// Upsert creates or updates the settings
func (r *repository) Upsert(ctx context.Context, settings *domain.Settings) (*domain.Settings, apperrors.ApplicationError) {
	now := time.Now()

	// Check if settings exist
	existing, _ := r.Get(ctx)

	if existing.ID == "" {
		// Create new settings
		settings.ID = uuid.New().String()
		settings.CreatedAt = now
		settings.UpdatedAt = now

		query := `
			INSERT INTO settings (
				id, business_name, business_latitude, business_longitude,
				default_map_latitude, default_map_longitude, default_map_zoom,
				default_item_weight, delivery_base_price, delivery_price_per_km,
				delivery_price_per_kg, created_at, updated_at
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		`

		_, err := r.db.ExecContext(ctx, query,
			settings.ID, settings.BusinessName, settings.BusinessLatitude, settings.BusinessLongitude,
			settings.DefaultMapLatitude, settings.DefaultMapLongitude, settings.DefaultMapZoom,
			settings.DefaultItemWeight, settings.DeliveryBasePrice, settings.DeliveryPricePerKm,
			settings.DeliveryPricePerKg, settings.CreatedAt, settings.UpdatedAt,
		)

		if err != nil {
			return nil, apperrors.NewApplicationError(mappings.SettingsCreateError, err)
		}
	} else {
		// Update existing settings
		settings.ID = existing.ID
		settings.CreatedAt = existing.CreatedAt
		settings.UpdatedAt = now

		query := `
			UPDATE settings SET
				business_name = $1, business_latitude = $2, business_longitude = $3,
				default_map_latitude = $4, default_map_longitude = $5, default_map_zoom = $6,
				default_item_weight = $7, delivery_base_price = $8, delivery_price_per_km = $9,
				delivery_price_per_kg = $10, updated_at = $11
			WHERE id = $12
		`

		_, err := r.db.ExecContext(ctx, query,
			settings.BusinessName, settings.BusinessLatitude, settings.BusinessLongitude,
			settings.DefaultMapLatitude, settings.DefaultMapLongitude, settings.DefaultMapZoom,
			settings.DefaultItemWeight, settings.DeliveryBasePrice, settings.DeliveryPricePerKm,
			settings.DeliveryPricePerKg, settings.UpdatedAt, settings.ID,
		)

		if err != nil {
			return nil, apperrors.NewApplicationError(mappings.SettingsUpdateError, err)
		}
	}

	return settings, nil
}
