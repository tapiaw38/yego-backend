package admin

import (
	"context"

	"yego/internal/platform/appcontext"
	apperrors "yego/internal/platform/errors"
)

type ListCouponsOutput struct {
	Coupons []*CouponOutput `json:"coupons"`
}

type ListCouponsUsecase interface {
	Execute(ctx context.Context) (*ListCouponsOutput, apperrors.ApplicationError)
}

type listCouponsUsecase struct {
	contextFactory appcontext.Factory
}

func NewListCouponsUsecase(contextFactory appcontext.Factory) ListCouponsUsecase {
	return &listCouponsUsecase{contextFactory: contextFactory}
}

func (u *listCouponsUsecase) Execute(ctx context.Context) (*ListCouponsOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	coupons, appErr := app.Repositories.Coupon.List(ctx)
	if appErr != nil {
		return nil, appErr
	}

	out := &ListCouponsOutput{Coupons: make([]*CouponOutput, 0, len(coupons))}
	for _, c := range coupons {
		out.Coupons = append(out.Coupons, toCouponOutput(c))
	}
	return out, nil
}
