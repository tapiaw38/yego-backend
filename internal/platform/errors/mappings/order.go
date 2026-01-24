package mappings

import "net/http"

// Order-related error mappings
var (
	OrderNotFoundError = ErrorDetails{
		Code:       "order:not-found",
		StatusCode: http.StatusNotFound,
		Message:    "order not found",
	}

	OrderCreateError = ErrorDetails{
		Code:       "order:create-error",
		StatusCode: http.StatusInternalServerError,
		Message:    "failed to create order",
	}

	OrderUpdateError = ErrorDetails{
		Code:       "order:update-error",
		StatusCode: http.StatusInternalServerError,
		Message:    "failed to update order",
	}

	OrderInvalidStatusError = ErrorDetails{
		Code:       "order:invalid-status",
		StatusCode: http.StatusBadRequest,
		Message:    "invalid order status",
	}

	OrderInvalidIDError = ErrorDetails{
		Code:       "order:invalid-id",
		StatusCode: http.StatusBadRequest,
		Message:    "invalid order ID format",
	}

	// Order Token errors
	OrderTokenCreateError = ErrorDetails{
		Code:       "order:token:create-error",
		StatusCode: http.StatusInternalServerError,
		Message:    "failed to create order token",
	}

	OrderTokenNotFoundError = ErrorDetails{
		Code:       "order:token:not-found",
		StatusCode: http.StatusNotFound,
		Message:    "order token not found",
	}

	OrderTokenExpiredError = ErrorDetails{
		Code:       "order:token:expired",
		StatusCode: http.StatusBadRequest,
		Message:    "order token has expired",
	}

	OrderTokenAlreadyClaimedError = ErrorDetails{
		Code:       "order:token:already-claimed",
		StatusCode: http.StatusConflict,
		Message:    "order has already been claimed",
	}

	OrderAlreadyAssignedError = ErrorDetails{
		Code:       "order:already-assigned",
		StatusCode: http.StatusConflict,
		Message:    "order is already assigned to a user",
	}
)
