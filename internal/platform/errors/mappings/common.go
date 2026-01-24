package mappings

import "net/http"

// Common error mappings
var (
	RequestBodyParsingError = ErrorDetails{
		Code:       "common:request-body-parsing-error",
		StatusCode: http.StatusBadRequest,
		Message:    "invalid request body",
	}

	InternalServerError = ErrorDetails{
		Code:       "common:internal-server-error",
		StatusCode: http.StatusInternalServerError,
		Message:    "internal server error",
	}

	UnauthorizedError = ErrorDetails{
		Code:       "common:unauthorized",
		StatusCode: http.StatusUnauthorized,
		Message:    "unauthorized access",
	}

	ForbiddenError = ErrorDetails{
		Code:       "common:forbidden",
		StatusCode: http.StatusForbidden,
		Message:    "access forbidden",
	}
)
