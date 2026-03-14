package admin

import (
	"context"

	"yego/internal/platform/appcontext"
	apperrors "yego/internal/platform/errors"
)

type DeleteCouponUsecase interface {
	Execute(ctx context.Context, id string) apperrors.ApplicationError
}

type deleteCouponUsecase struct {
	contextFactory appcontext.Factory
}

func NewDeleteCouponUsecase(contextFactory appcontext.Factory) DeleteCouponUsecase {
	return &deleteCouponUsecase{contextFactory: contextFactory}
}

func (u *deleteCouponUsecase) Execute(ctx context.Context, id string) apperrors.ApplicationError {
	app := u.contextFactory()
	return app.Repositories.Coupon.Delete(ctx, id)
}
