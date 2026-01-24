package domain

import "time"

// OrderStatus represents the possible states of an order
type OrderStatus string

const (
	StatusCreated   OrderStatus = "CREATED"
	StatusConfirmed OrderStatus = "CONFIRMED"
	StatusPreparing OrderStatus = "PREPARING"
	StatusOnTheWay  OrderStatus = "ON_THE_WAY"
	StatusDelivered OrderStatus = "DELIVERED"
	StatusCancelled OrderStatus = "CANCELLED"
)

// ValidStatuses contains all valid order statuses in order
var ValidStatuses = []OrderStatus{
	StatusCreated,
	StatusConfirmed,
	StatusPreparing,
	StatusOnTheWay,
	StatusDelivered,
	StatusCancelled,
}

// IsValidStatus checks if a status string is valid
func IsValidStatus(s string) bool {
	for _, status := range ValidStatuses {
		if string(status) == s {
			return true
		}
	}
	return false
}

// Order represents a customer order in the system
type Order struct {
	ID        string      `json:"id"`
	ProfileID *string     `json:"profile_id,omitempty"`
	UserID    *string     `json:"user_id,omitempty"`
	Status    OrderStatus `json:"status"`
	ETA       string      `json:"eta"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

// OrderToken represents a token for claiming an order via link
type OrderToken struct {
	ID              string     `json:"id"`
	OrderID         string     `json:"order_id"`
	Token           string     `json:"token"`
	PhoneNumber     *string    `json:"phone_number,omitempty"`
	ClaimedAt       *time.Time `json:"claimed_at,omitempty"`
	ClaimedByUserID *string    `json:"claimed_by_user_id,omitempty"`
	ExpiresAt       time.Time  `json:"expires_at"`
	CreatedAt       time.Time  `json:"created_at"`
}

// StatusIndex returns the position of the current status in the workflow
func (o *Order) StatusIndex() int {
	for i, s := range ValidStatuses {
		if s == o.Status {
			return i
		}
	}
	return -1
}
