package admin

import (
	"yego/internal/domain"
)

// OrderOutput represents an order in the admin list
type OrderOutput struct {
	ID                 string            `json:"id"`
	ProfileID          *string           `json:"profile_id,omitempty"`
	UserID             *string           `json:"user_id,omitempty"`
	Status             string            `json:"status"`
	StatusMessage      *string           `json:"status_message,omitempty"`
	StatusIndex        int               `json:"status_index"`
	ETA                string            `json:"eta"`
	Data               *domain.OrderData `json:"data,omitempty"`
	DeliveryUserID     *string           `json:"delivery_user_id,omitempty"`
	DeliveryAcceptedAt *string           `json:"delivery_accepted_at,omitempty"`
	CreatedAt          string            `json:"created_at"`
	UpdatedAt          string            `json:"updated_at"`
	AllStatuses        []string          `json:"all_statuses"`
}

// ProfileOutput represents a profile in the admin list
type ProfileOutput struct {
	ID          string          `json:"id"`
	UserID      string          `json:"user_id"`
	PhoneNumber string          `json:"phone_number"`
	Location    *LocationOutput `json:"location,omitempty"`
	CreatedAt   string          `json:"created_at"`
	UpdatedAt   string          `json:"updated_at"`
}

// LocationOutput represents a location in the output
type LocationOutput struct {
	ID        string  `json:"id"`
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
	Address   string  `json:"address"`
}

// TransactionOutput represents a transaction in the admin list
type TransactionOutput struct {
	ID               string  `json:"id"`
	OrderID          string  `json:"order_id"`
	UserID           string  `json:"user_id"`
	ProfileID        *string `json:"profile_id,omitempty"`
	Amount           float64 `json:"amount"`
	Currency         string  `json:"currency"`
	Status           string  `json:"status"`
	PaymentID        *int    `json:"payment_id,omitempty"`
	GatewayPaymentID *string `json:"gateway_payment_id,omitempty"`
	CollectorID      *string `json:"collector_id,omitempty"`
	Description      *string `json:"description,omitempty"`
	CreatedAt        string  `json:"created_at"`
	UpdatedAt        string  `json:"updated_at"`
}

// toOrderOutput converts a domain order to output
func toOrderOutput(order *domain.Order) OrderOutput {
	allStatuses := make([]string, len(domain.ValidStatuses))
	for i, s := range domain.ValidStatuses {
		allStatuses[i] = string(s)
	}

	output := OrderOutput{
		ID:            order.ID,
		ProfileID:     order.ProfileID,
		UserID:        order.UserID,
		Status:        string(order.Status),
		StatusMessage: order.StatusMessage,
		StatusIndex:   order.StatusIndex(),
		ETA:           order.ETA,
		Data:          order.Data,
		DeliveryUserID: order.DeliveryUserID,
		CreatedAt:     order.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:     order.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		AllStatuses:   allStatuses,
	}

	if order.DeliveryAcceptedAt != nil {
		formatted := order.DeliveryAcceptedAt.Format("2006-01-02T15:04:05Z")
		output.DeliveryAcceptedAt = &formatted
	}

	return output
}

// toTransactionOutput converts a domain transaction to output
func toTransactionOutput(transaction *domain.Transaction) TransactionOutput {
	return TransactionOutput{
		ID:               transaction.ID,
		OrderID:          transaction.OrderID,
		UserID:           transaction.UserID,
		ProfileID:        transaction.ProfileID,
		Amount:           transaction.Amount,
		Currency:         transaction.Currency,
		Status:           transaction.Status,
		PaymentID:        transaction.PaymentID,
		GatewayPaymentID: transaction.GatewayPaymentID,
		CollectorID:      transaction.CollectorID,
		Description:      transaction.Description,
		CreatedAt:        transaction.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:        transaction.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}
