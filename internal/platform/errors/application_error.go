package apperrors

import (
	"log"

	"github.com/gin-gonic/gin"
	"wappi/internal/platform/errors/mappings"
)

// ApplicationError represents a structured error for the application
type ApplicationError interface {
	error
	Code() string
	StatusCode() int
	Message() string
	Log(c *gin.Context)
	OriginalError() error
}

type applicationError struct {
	code          string
	statusCode    int
	message       string
	originalError error
}

// NewApplicationError creates a new application error from error details
func NewApplicationError(details mappings.ErrorDetails, originalError error) ApplicationError {
	return &applicationError{
		code:          details.Code,
		statusCode:    details.StatusCode,
		message:       details.Message,
		originalError: originalError,
	}
}

func (e *applicationError) Error() string {
	return e.message
}

func (e *applicationError) Code() string {
	return e.code
}

func (e *applicationError) StatusCode() int {
	return e.statusCode
}

func (e *applicationError) Message() string {
	return e.message
}

func (e *applicationError) OriginalError() error {
	return e.originalError
}

func (e *applicationError) Log(c *gin.Context) {
	if e.originalError != nil {
		log.Printf("[ERROR] %s: %s - Original: %v", e.code, e.message, e.originalError)
	} else {
		log.Printf("[ERROR] %s: %s", e.code, e.message)
	}
}

// MarshalJSON implements json.Marshaler for API responses
func (e *applicationError) MarshalJSON() ([]byte, error) {
	return []byte(`{"code":"` + e.code + `","message":"` + e.message + `"}`), nil
}
