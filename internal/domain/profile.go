package domain

import "time"

// Profile represents a user profile in the system
type Profile struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	PhoneNumber string    `json:"phone_number"`
	LocationID  *string   `json:"location_id,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// IsCompleted checks if the profile has all required fields filled
// A profile is considered complete when it has a phone number and a location
func (p *Profile) IsCompleted() bool {
	hasPhoneNumber := p.PhoneNumber != ""
	hasLocation := p.LocationID != nil && *p.LocationID != ""
	return hasPhoneNumber && hasLocation
}

// ProfileLocation represents a geographic location for a profile
type ProfileLocation struct {
	ID        string  `json:"id"`
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
	Address   string  `json:"address"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ProfileToken represents a token for profile completion link
type ProfileToken struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Token     string    `json:"token"`
	Used      bool      `json:"used"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}
