package admin

import (
	"yego/internal/domain"
)

// CouponOutput represents a coupon in the admin API response
type CouponOutput struct {
	ID                string   `json:"id"`
	Code              string   `json:"code"`
	Description       *string  `json:"description,omitempty"`
	DiscountType      string   `json:"discount_type"`
	DiscountValue     float64  `json:"discount_value"`
	MaxUses           *int     `json:"max_uses,omitempty"`
	CurrentUses       int      `json:"current_uses"`
	UsageLimitPerUser int      `json:"usage_limit_per_user"`
	MinOrderAmount    *float64 `json:"min_order_amount,omitempty"`
	ValidFrom         *string  `json:"valid_from,omitempty"`
	ValidUntil        *string  `json:"valid_until,omitempty"`
	Active            bool     `json:"active"`
	IconURL           *string  `json:"icon_url,omitempty"`
	CoverURL          *string  `json:"cover_url,omitempty"`
	CreatedAt         string   `json:"created_at"`
	UpdatedAt         string   `json:"updated_at"`
}

func toCouponOutput(c *domain.Coupon) *CouponOutput {
	out := &CouponOutput{
		ID:                c.ID,
		Code:              c.Code,
		Description:       c.Description,
		DiscountType:      string(c.DiscountType),
		DiscountValue:     c.DiscountValue,
		MaxUses:           c.MaxUses,
		CurrentUses:       c.CurrentUses,
		UsageLimitPerUser: c.UsageLimitPerUser,
		MinOrderAmount:    c.MinOrderAmount,
		Active:            c.Active,
		IconURL:           c.IconURL,
		CoverURL:          c.CoverURL,
		CreatedAt:         c.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:         c.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
	if c.ValidFrom != nil {
		s := c.ValidFrom.Format("2006-01-02T15:04:05Z")
		out.ValidFrom = &s
	}
	if c.ValidUntil != nil {
		s := c.ValidUntil.Format("2006-01-02T15:04:05Z")
		out.ValidUntil = &s
	}
	return out
}

// OrderOutput represents an order in the admin list
type OrderOutput struct {
	ID            string            `json:"id"`
	ProfileID     *string           `json:"profile_id,omitempty"`
	UserID        *string           `json:"user_id,omitempty"`
	Status        string            `json:"status"`
	StatusMessage *string           `json:"status_message,omitempty"`
	StatusIndex   int               `json:"status_index"`
	ETA           string            `json:"eta"`
	Data          *domain.OrderData `json:"data,omitempty"`
	CreatedAt     string            `json:"created_at"`
	UpdatedAt     string            `json:"updated_at"`
	AllStatuses   []string          `json:"all_statuses"`
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

	return OrderOutput{
		ID:            order.ID,
		ProfileID:     order.ProfileID,
		UserID:        order.UserID,
		Status:        string(order.Status),
		StatusMessage: order.StatusMessage,
		StatusIndex:   order.StatusIndex(),
		ETA:           order.ETA,
		Data:          order.Data,
		CreatedAt:     order.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:     order.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		AllStatuses:   allStatuses,
	}
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
