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
)

// Error is the standard error interface
type Error struct {
	Type    Type   `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"error,omitempty"`
}

// newError is a helper function to create a new Error
func newError(e *Error) *Error {
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

func NewInternal(err error) *Error {
	return newError(&Error{
		Type:    Internal,
		Message: "Internal server error",
		Data:    err,
	})
}

func NewConflict(err error, resource ...map[string]any) *Error {
	message := "Resource already exists"
	if len(resource) > 0 {
		message = fmt.Sprintf("Resource: %v already exists", resource)
	}
	return newError(&Error{
		Type:    Conflict,
		Message: message,
		Data:    err,
	})
}

func NewNotFound(err error, resource ...map[string]any) *Error {
	message := "Resource not found"
	if len(resource) > 0 {
		message = fmt.Sprintf("Resource: %v not found", resource)
	}
	return newError(&Error{
		Type:    NotFound,
		Message: message,
		Data:    err,
	})
}

func NewAuthorization(err error, reason ...string) *Error {
	message := "Authorization failed"
	if len(reason) > 0 {
		message = fmt.Sprintf("Authorization failed. Reason: %v", reason[0])
	}
	return newError(&Error{
		Type:    Authorization,
		Message: message,
		Data:    err,
	})
}
