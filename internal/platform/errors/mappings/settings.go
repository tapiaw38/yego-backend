package mappings

import "net/http"

// Settings-related error mappings
var (
	SettingsGetError = ErrorDetails{
		Code:       "settings:get-error",
		StatusCode: http.StatusInternalServerError,
		Message:    "failed to get settings",
	}

	SettingsCreateError = ErrorDetails{
		Code:       "settings:create-error",
		StatusCode: http.StatusInternalServerError,
		Message:    "failed to create settings",
	}

	SettingsUpdateError = ErrorDetails{
		Code:       "settings:update-error",
		StatusCode: http.StatusInternalServerError,
		Message:    "failed to update settings",
	}
)
