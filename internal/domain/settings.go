package domain

import "time"

// Settings represents the application configuration
type Settings struct {
	ID                  string    `json:"id"`
	BusinessName        string    `json:"business_name"`
	BusinessLatitude    float64   `json:"business_latitude"`
	BusinessLongitude   float64   `json:"business_longitude"`
	DefaultMapLatitude  float64   `json:"default_map_latitude"`
	DefaultMapLongitude float64   `json:"default_map_longitude"`
	DefaultMapZoom      int       `json:"default_map_zoom"`
	DefaultItemWeight   int       `json:"default_item_weight"` // in grams
	DeliveryBasePrice   float64   `json:"delivery_base_price"`
	DeliveryPricePerKm  float64   `json:"delivery_price_per_km"`
	DeliveryPricePerKg  float64   `json:"delivery_price_per_kg"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}
