package error_response

import (
	"net/http"
	"time"
)

// ErrorResponse represents a structured error response
type ErrorResponse struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
	Now     int64  `json:"now"`
}

// NewErrorResponse creates a new ErrorResponse
func NewErrorResponse(message string, status int) *ErrorResponse {
	return &ErrorResponse{
		Message: message,
		Status:  status,
		Now:     time.Now().Unix(),
	}
}

// BadRequestError represents a 400 Bad Request error
func BadRequestError(message string) *ErrorResponse {
	if message == "" {
		message = http.StatusText(http.StatusBadRequest)
	}
	return NewErrorResponse(message, http.StatusBadRequest)
}

// NotFoundError represents a 404 Not Found error
func NotFoundError(message string) *ErrorResponse {
	if message == "" {
		message = http.StatusText(http.StatusNotFound)
	}
	return NewErrorResponse(message, http.StatusNotFound)
}

// UnauthorizedError represents a 401 Unauthorized error
func UnauthorizedError(message string) *ErrorResponse {
	if message == "" {
		message = http.StatusText(http.StatusUnauthorized)
	}
	return NewErrorResponse(message, http.StatusUnauthorized)
}

// GoneError represents a 410 Gone error
func GoneError(message string) *ErrorResponse {
	if message == "" {
		message = http.StatusText(http.StatusGone)
	}
	return NewErrorResponse(message, http.StatusGone)
}

// InternalServerError represents a 500 Internal Server Error
func InternalServerError(message string) *ErrorResponse {
	if message == "" {
		message = http.StatusText(http.StatusInternalServerError)
	}
	return NewErrorResponse(message, http.StatusInternalServerError)
}

func ServiceUnavailable(message string) *ErrorResponse {
	if message == "" {
		message = http.StatusText(http.StatusServiceUnavailable)
	}
	return NewErrorResponse(message, http.StatusServiceUnavailable)
}
