package notification

// OrderClaimedPayload contains data about a claimed order
type OrderClaimedPayload struct {
	OrderID   string `json:"order_id"`
	UserID    string `json:"user_id"`
	ProfileID string `json:"profile_id,omitempty"`
	Status    string `json:"status"`
	ETA       string `json:"eta"`
	ClaimedAt string `json:"claimed_at"`
}

// OrderCreatedPayload contains data about a newly created order
type OrderCreatedPayload struct {
	OrderID   string `json:"order_id"`
	ProfileID string `json:"profile_id,omitempty"`
	Status    string `json:"status"`
	ETA       string `json:"eta"`
	CreatedAt string `json:"created_at"`
}

// OrderUpdatedPayload contains data about an order whose status changed
type OrderUpdatedPayload struct {
	OrderID string `json:"order_id"`
	Status  string `json:"status"`
	ETA     string `json:"eta"`
}

// Service defines the interface for sending notifications to clients
// This is a driven port (output port) in hexagonal architecture
type Service interface {
	// NotifyOrderClaimed sends a notification when an order is claimed by a user
	NotifyOrderClaimed(payload OrderClaimedPayload) error
	// NotifyOrderCreated sends a notification when a new order is created
	NotifyOrderCreated(payload OrderCreatedPayload) error
	// NotifyOrderUpdated sends a notification when an order's status changes
	NotifyOrderUpdated(payload OrderUpdatedPayload) error
}
