package mappings

import "net/http"

// Profile-related error mappings
var (
	ProfileNotFoundError = ErrorDetails{
		Code:       "profile:not-found",
		StatusCode: http.StatusNotFound,
		Message:    "profile not found",
	}

	ProfileCreateError = ErrorDetails{
		Code:       "profile:create-error",
		StatusCode: http.StatusInternalServerError,
		Message:    "failed to create profile",
	}

	ProfileUpdateError = ErrorDetails{
		Code:       "profile:update-error",
		StatusCode: http.StatusInternalServerError,
		Message:    "failed to update profile",
	}

	ProfileTokenNotFoundError = ErrorDetails{
		Code:       "profile:token-not-found",
		StatusCode: http.StatusNotFound,
		Message:    "profile token not found or expired",
	}

	ProfileTokenExpiredError = ErrorDetails{
		Code:       "profile:token-expired",
		StatusCode: http.StatusBadRequest,
		Message:    "profile token has expired",
	}

	ProfileTokenUsedError = ErrorDetails{
		Code:       "profile:token-used",
		StatusCode: http.StatusBadRequest,
		Message:    "profile token has already been used",
	}

	ProfileTokenCreateError = ErrorDetails{
		Code:       "profile:token-create-error",
		StatusCode: http.StatusInternalServerError,
		Message:    "failed to create profile token",
	}

	LocationCreateError = ErrorDetails{
		Code:       "profile:location-create-error",
		StatusCode: http.StatusInternalServerError,
		Message:    "failed to create location",
	}

	InvalidUserIDError = ErrorDetails{
		Code:       "profile:invalid-user-id",
		StatusCode: http.StatusBadRequest,
		Message:    "invalid user ID format",
	}
)
