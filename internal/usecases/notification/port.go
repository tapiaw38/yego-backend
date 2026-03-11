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

// OrderAssignedToDeliveryPayload contains data when an order is assigned to delivery users
type OrderAssignedToDeliveryPayload struct {
	OrderID    string `json:"order_id"`
	Status     string `json:"status"`
	ETA        string `json:"eta"`
	AssignedAt string `json:"assigned_at"`
}

// DeliveryAcceptedPayload contains data when a delivery user accepts an order
type DeliveryAcceptedPayload struct {
	OrderID            string `json:"order_id"`
	DeliveryUserID     string `json:"delivery_user_id"`
	DeliveryAcceptedAt string `json:"delivery_accepted_at"`
}

// OrderStatusUpdatedPayload contains data when an order status changes
type OrderStatusUpdatedPayload struct {
	OrderID        string `json:"order_id"`
	Status         string `json:"status"`
	StatusMessage  string `json:"status_message,omitempty"`
	ETA            string `json:"eta"`
}

// Service defines the interface for sending notifications to clients
type Service interface {
	NotifyOrderClaimed(payload OrderClaimedPayload) error
	NotifyDeliveryUsers(payload OrderAssignedToDeliveryPayload) error
	NotifyManagersDeliveryAccepted(payload DeliveryAcceptedPayload) error
	NotifyDeliveryUserOrderUpdated(deliveryUserID string, payload OrderStatusUpdatedPayload) error
}
