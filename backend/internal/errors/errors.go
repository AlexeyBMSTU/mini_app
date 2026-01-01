package errors

import (
	"encoding/json"
	"fmt"
	"mini-app-backend/internal/logger"
	"net/http"
)

type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

func (e *AppError) Error() string {
	return e.Message
}

func NewAppError(code int, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}

func NewAppErrorWithDetails(code int, message, details string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Details: details,
	}
}

var (
	ErrBadRequest           = NewAppError(http.StatusBadRequest, "Bad request")
	ErrUnauthorized         = NewAppError(http.StatusUnauthorized, "Unauthorized")
	ErrForbidden           = NewAppError(http.StatusForbidden, "Forbidden")
	ErrNotFound            = NewAppError(http.StatusNotFound, "Not found")
	ErrMethodNotAllowed    = NewAppError(http.StatusMethodNotAllowed, "Method not allowed")
	ErrInternalServerError = NewAppError(http.StatusInternalServerError, "Internal server error")
	ErrServiceUnavailable  = NewAppError(http.StatusServiceUnavailable, "Service unavailable")
)

func SendErrorResponse(w http.ResponseWriter, err error) {
	var appErr *AppError
	var statusCode int
	var message string

	if err, ok := err.(*AppError); ok {
		appErr = err
		statusCode = err.Code
		message = err.Message
	} else {
		appErr = ErrInternalServerError
		statusCode = ErrInternalServerError.Code
		message = ErrInternalServerError.Message
		logger.Errorf("Unexpected error: %v", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := map[string]interface{}{
		"error": message,
	}

	if appErr.Details != "" {
		response["details"] = appErr.Details
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.Errorf("Error encoding error response: %v", err)
	}
}

func HandleHTTPError(w http.ResponseWriter, err error) {
	if err == nil {
		return
	}

	switch e := err.(type) {
	case *AppError:
		SendErrorResponse(w, e)
	default:
		SendErrorResponse(w, ErrInternalServerError)
	}
}

func WrapError(err error, message string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", message, err)
}

func IsNotFound(err error) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Code == http.StatusNotFound
	}
	return false
}

func IsBadRequest(err error) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Code == http.StatusBadRequest
	}
	return false
}

func IsUnauthorized(err error) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Code == http.StatusUnauthorized
	}
	return false
}

func IsForbidden(err error) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Code == http.StatusForbidden
	}
	return false
}

func IsInternalServerError(err error) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Code == http.StatusInternalServerError
	}
	return false
}