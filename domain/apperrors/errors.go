package apperrors

import (
	"errors"
	"fmt"
	"github.com/j03hanafi/seternak-backend/utils/consts"
	"github.com/spf13/viper"
	"net/http"
)

// Type holds a type string and integer code for the error
type Type string

// Set of valid errorTypes
const (
	Internal      Type = "E000" // Server (500) and fallback errors
	Conflict      Type = "E001" // Already exists (eg, create account with existent email) - 409
	NotFound      Type = "E002" // For not finding resource
	Authorization Type = "E003" // Authentication Failures
	BadRequest    Type = "E004" // Validation errors / BadInput
)

// Error is the standard error interface
type Error struct {
	Type    Type   `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"error,omitempty"`
}

// newError is a helper function to create a new Error
func newError(err error, e *Error) *Error {
	// Check if the error is an Error type
	var appErr *Error
	if errors.As(err, &appErr) {
		e.Data = appErr.Data
	} else {
		e.Data = err
	}

	// Check if app is in production mode
	if viper.GetString("APP_ENV") == consts.ProductionMode {
		e.Data = nil
	}

	return e
}

// Error satisfies standard error interface
func (e Error) Error() string {
	return e.Message
}

// Status is a mapping errors to status codes
func (e Error) Status() int {
	switch e.Type {
	case Internal:
		return http.StatusInternalServerError
	case Conflict:
		return http.StatusConflict
	case NotFound:
		return http.StatusNotFound
	case Authorization:
		return http.StatusUnauthorized
	case BadRequest:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}

// Status checks the runtime type
// of the error and returns a http
// status code if the error is Error
func Status(err error) int {
	var e *Error
	if errors.As(err, &e) {
		return e.Status()
	}
	return http.StatusInternalServerError
}

/*
	Error "Factories"
*/

// NewInternal returns a 500 Internal Server Error
func NewInternal(err error) *Error {
	return newError(err, &Error{
		Type:    Internal,
		Message: "Internal server error",
	})
}

// NewConflict returns a 409 Conflict Error
func NewConflict(err error, resource ...any) *Error {
	message := "Resource already exists"
	if len(resource) > 0 {
		message = fmt.Sprintf("Resource: %v already exists", resource)
	}
	return newError(err, &Error{
		Type:    Conflict,
		Message: message,
	})
}

// NewNotFound returns a 404 Not Found Error
func NewNotFound(err error, resource ...any) *Error {
	message := "Resource not found"
	if len(resource) > 0 {
		message = fmt.Sprintf("Resource: %v not found", resource)
	}
	return newError(err, &Error{
		Type:    NotFound,
		Message: message,
	})
}

// NewAuthorization returns a 401 Unauthorized Error
func NewAuthorization(err error, reason ...string) *Error {
	message := "Authorization failed"
	if len(reason) > 0 {
		message = fmt.Sprintf("Authorization failed. Reason: %v", reason[0])
	}
	return newError(err, &Error{
		Type:    Authorization,
		Message: message,
	})
}

// NewBadRequest returns a 400 Bad Request Error
func NewBadRequest(err error, reason ...string) *Error {
	message := "Bad request"
	if len(reason) > 0 {
		message = fmt.Sprintf("Bad request. Reason: %v", reason[0])
	}
	return newError(err, &Error{
		Type:    BadRequest,
		Message: message,
	})
}
