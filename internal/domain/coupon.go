package domain

import "time"

type DiscountType string

const (
	DiscountTypePercentage DiscountType = "PERCENTAGE"
	DiscountTypeFixed      DiscountType = "FIXED"
)

type Coupon struct {
	ID                string       `json:"id"`
	Code              string       `json:"code"`
	Description       *string      `json:"description,omitempty"`
	DiscountType      DiscountType `json:"discount_type"`
	DiscountValue     float64      `json:"discount_value"`
	MaxUses           *int         `json:"max_uses,omitempty"`
	CurrentUses       int          `json:"current_uses"`
	UsageLimitPerUser int          `json:"usage_limit_per_user"`
	MinOrderAmount    *float64     `json:"min_order_amount,omitempty"`
	ValidFrom         *time.Time   `json:"valid_from,omitempty"`
	ValidUntil        *time.Time   `json:"valid_until,omitempty"`
	Active            bool         `json:"active"`
	IconURL           *string      `json:"icon_url,omitempty"`
	CoverURL          *string      `json:"cover_url,omitempty"`
	CreatedAt         time.Time    `json:"created_at"`
	UpdatedAt         time.Time    `json:"updated_at"`
}
