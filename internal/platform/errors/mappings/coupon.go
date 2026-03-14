package mappings

import "net/http"

var (
	CouponCreateError = ErrorDetails{
		Code:       "coupon:create-error",
		StatusCode: http.StatusInternalServerError,
		Message:    "failed to create coupon",
	}

	CouponGetError = ErrorDetails{
		Code:       "coupon:get-error",
		StatusCode: http.StatusInternalServerError,
		Message:    "failed to get coupon",
	}

	CouponNotFoundError = ErrorDetails{
		Code:       "coupon:not-found",
		StatusCode: http.StatusNotFound,
		Message:    "coupon not found",
	}

	CouponListError = ErrorDetails{
		Code:       "coupon:list-error",
		StatusCode: http.StatusInternalServerError,
		Message:    "failed to list coupons",
	}

	CouponUpdateError = ErrorDetails{
		Code:       "coupon:update-error",
		StatusCode: http.StatusInternalServerError,
		Message:    "failed to update coupon",
	}

	CouponDeleteError = ErrorDetails{
		Code:       "coupon:delete-error",
		StatusCode: http.StatusInternalServerError,
		Message:    "failed to delete coupon",
	}

	CouponDuplicateCodeError = ErrorDetails{
		Code:       "coupon:duplicate-code",
		StatusCode: http.StatusConflict,
		Message:    "coupon code already exists",
	}
)
