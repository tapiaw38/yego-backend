package admin

import (
	"context"
	"strings"
	"time"

	"yego/internal/domain"
	"yego/internal/platform/appcontext"
	apperrors "yego/internal/platform/errors"
	"yego/internal/platform/errors/mappings"
)

type CreateCouponInput struct {
	Code              string   `json:"code" binding:"required"`
	Description       *string  `json:"description"`
	DiscountType      string   `json:"discount_type" binding:"required"`
	DiscountValue     float64  `json:"discount_value" binding:"required"`
	MaxUses           *int     `json:"max_uses"`
	UsageLimitPerUser int      `json:"usage_limit_per_user"`
	MinOrderAmount    *float64 `json:"min_order_amount"`
	ValidFrom         *string  `json:"valid_from"`
	ValidUntil        *string  `json:"valid_until"`
	Active            bool     `json:"active"`
	IconURL           *string  `json:"icon_url"`
	CoverURL          *string  `json:"cover_url"`
}

type CreateCouponUsecase interface {
	Execute(ctx context.Context, input CreateCouponInput) (*CouponOutput, apperrors.ApplicationError)
}

type createCouponUsecase struct {
	contextFactory appcontext.Factory
}

func NewCreateCouponUsecase(contextFactory appcontext.Factory) CreateCouponUsecase {
	return &createCouponUsecase{contextFactory: contextFactory}
}

func (u *createCouponUsecase) Execute(ctx context.Context, input CreateCouponInput) (*CouponOutput, apperrors.ApplicationError) {
	discountType := domain.DiscountType(strings.ToUpper(input.DiscountType))
	if discountType != domain.DiscountTypePercentage && discountType != domain.DiscountTypeFixed {
		return nil, apperrors.NewApplicationError(mappings.RequestBodyParsingError, nil)
	}

	usageLimit := input.UsageLimitPerUser
	if usageLimit <= 0 {
		usageLimit = 1
	}

	coupon := &domain.Coupon{
		Code:              input.Code,
		Description:       input.Description,
		DiscountType:      discountType,
		DiscountValue:     input.DiscountValue,
		MaxUses:           input.MaxUses,
		UsageLimitPerUser: usageLimit,
		MinOrderAmount:    input.MinOrderAmount,
		Active:            input.Active,
		IconURL:           input.IconURL,
		CoverURL:          input.CoverURL,
	}

	if input.ValidFrom != nil {
		t, err := time.Parse(time.RFC3339, *input.ValidFrom)
		if err != nil {
			return nil, apperrors.NewApplicationError(mappings.RequestBodyParsingError, err)
		}
		coupon.ValidFrom = &t
	}
	if input.ValidUntil != nil {
		t, err := time.Parse(time.RFC3339, *input.ValidUntil)
		if err != nil {
			return nil, apperrors.NewApplicationError(mappings.RequestBodyParsingError, err)
		}
		coupon.ValidUntil = &t
	}

	app := u.contextFactory()
	created, appErr := app.Repositories.Coupon.Create(ctx, coupon)
	if appErr != nil {
		return nil, appErr
	}

	return toCouponOutput(created), nil
}
