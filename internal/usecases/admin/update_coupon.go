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

type UpdateCouponInput struct {
	Code              *string  `json:"code"`
	Description       *string  `json:"description"`
	DiscountType      *string  `json:"discount_type"`
	DiscountValue     *float64 `json:"discount_value"`
	MaxUses           *int     `json:"max_uses"`
	UsageLimitPerUser *int     `json:"usage_limit_per_user"`
	MinOrderAmount    *float64 `json:"min_order_amount"`
	ValidFrom         *string  `json:"valid_from"`
	ValidUntil        *string  `json:"valid_until"`
	Active            *bool    `json:"active"`
	IconURL           *string  `json:"icon_url"`
	CoverURL          *string  `json:"cover_url"`
}

type UpdateCouponUsecase interface {
	Execute(ctx context.Context, id string, input UpdateCouponInput) (*CouponOutput, apperrors.ApplicationError)
}

type updateCouponUsecase struct {
	contextFactory appcontext.Factory
}

func NewUpdateCouponUsecase(contextFactory appcontext.Factory) UpdateCouponUsecase {
	return &updateCouponUsecase{contextFactory: contextFactory}
}

func (u *updateCouponUsecase) Execute(ctx context.Context, id string, input UpdateCouponInput) (*CouponOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	existing, appErr := app.Repositories.Coupon.GetByID(ctx, id)
	if appErr != nil {
		return nil, appErr
	}

	if input.Code != nil {
		existing.Code = *input.Code
	}
	if input.Description != nil {
		existing.Description = input.Description
	}
	if input.DiscountType != nil {
		dt := domain.DiscountType(strings.ToUpper(*input.DiscountType))
		if dt != domain.DiscountTypePercentage && dt != domain.DiscountTypeFixed {
			return nil, apperrors.NewApplicationError(mappings.RequestBodyParsingError, nil)
		}
		existing.DiscountType = dt
	}
	if input.DiscountValue != nil {
		existing.DiscountValue = *input.DiscountValue
	}
	if input.MaxUses != nil {
		existing.MaxUses = input.MaxUses
	}
	if input.UsageLimitPerUser != nil {
		existing.UsageLimitPerUser = *input.UsageLimitPerUser
	}
	if input.MinOrderAmount != nil {
		existing.MinOrderAmount = input.MinOrderAmount
	}
	if input.Active != nil {
		existing.Active = *input.Active
	}
	if input.IconURL != nil {
		existing.IconURL = input.IconURL
	}
	if input.CoverURL != nil {
		existing.CoverURL = input.CoverURL
	}
	if input.ValidFrom != nil {
		if *input.ValidFrom == "" {
			existing.ValidFrom = nil
		} else {
			t, err := time.Parse(time.RFC3339, *input.ValidFrom)
			if err != nil {
				return nil, apperrors.NewApplicationError(mappings.RequestBodyParsingError, err)
			}
			existing.ValidFrom = &t
		}
	}
	if input.ValidUntil != nil {
		if *input.ValidUntil == "" {
			existing.ValidUntil = nil
		} else {
			t, err := time.Parse(time.RFC3339, *input.ValidUntil)
			if err != nil {
				return nil, apperrors.NewApplicationError(mappings.RequestBodyParsingError, err)
			}
			existing.ValidUntil = &t
		}
	}

	updated, appErr := app.Repositories.Coupon.Update(ctx, existing)
	if appErr != nil {
		return nil, appErr
	}
	return toCouponOutput(updated), nil
}
