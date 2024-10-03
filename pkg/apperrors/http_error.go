package apperrors

import (
	"errors"
	"fmt"
	"net/http"
)

// Type holds a type string and integer code for the error
type Type string

// Application-wide error codes
const (
	Authorization       Type = "AUTHORIZATION" // Authentication Failures -
	BadRequest          Type = "BAD_REQUEST"   // Validation errors / Bad input
	Conflict            Type = "CONFLICT"      // Already exists (eg, create account with existent email) - 409
	Internal            Type = "INTERNAL"      // Server (500) and fallback errors
	NotFound            Type = "NOT_FOUND"     // For not finding resource
	PayloadTooLarge     Type = "PAYLOAD_TOO_LARGE"
	UnprocessableEntity Type = "UNPROCESSABLE_ENTITY" // 422
	ErrInvalidClaims    Type = "INVALID_CLAIMS"       // Invalid JWT claims
)

// Error holds a custom error for the application
// implements the error interface
type Error struct {
	Type    Type   `json:"type"`
	Message string `json:"message"`
}

// Error satisfies the error interface
func (e Error) Error() string {
	return e.Message
}

// Status is a mapping errors to status codes
// Of course, this must be tailored to your application
func (e Error) Status() int {
	switch e.Type {
	case Authorization:
		return http.StatusUnauthorized
	case BadRequest:
		return http.StatusBadRequest
	case Conflict:
		return http.StatusConflict
	case Internal:
		return http.StatusInternalServerError
	case NotFound:
		return http.StatusNotFound
	case PayloadTooLarge:
		return http.StatusRequestEntityTooLarge
	case UnprocessableEntity:
		return http.StatusUnprocessableEntity
	case ErrInvalidClaims:
		return http.StatusUnauthorized
	default:
		return http.StatusInternalServerError
	}
}

// Status checks the runtime type of the error and returns an http status code if the error is of type Error
func Status(err error) int {
	var e *Error
	if errors.As(err, &e) {
		return e.Status()
	}
	return http.StatusInternalServerError
}

/*
* Error "constructors"
 */

// NewAuthorization to create a 401
func NewAuthorization(reason string) *Error {
	return &Error{
		Type:    Authorization,
		Message: reason,
	}
}

// NewBadRequest to create 400 errors (validation, for example)
func NewBadRequest(reason string) *Error {
	return &Error{
		Type:    BadRequest,
		Message: fmt.Sprintf("Bad request. Reason: %v", reason),
	}
}

// NewConflict to create an error for 409
func NewConflict(message string) *Error {
	return &Error{
		Type:    Conflict,
		Message: message,
	}
}

// NewInternal for 500 errors and unknown errors
func NewInternal() *Error {
	return &Error{
		Type:    Internal,
		Message: "Internal server error.",
	}
}

// NewNotFound to create an error for 404
func NewNotFound(message string) *Error {
	return &Error{
		Type:    NotFound,
		Message: message,
	}
}

// NewPayloadTooLarge to create an error for 413
func NewPayloadTooLarge(maxBodySize int64, contentLength int64) *Error {
	return &Error{
		Type:    PayloadTooLarge,
		Message: fmt.Sprintf("Max payload size of %v exceeded. Actual payload size: %v", maxBodySize, contentLength),
	}
}

func NewUnprocessableEntity(reason string) *Error {
	return &Error{
		Type:    UnprocessableEntity,
		Message: fmt.Sprintf("Unprocessable entity. Reason: %v", reason),
	}
}

// NewInvalidClaims to create an error for 401
func NewInvalidClaims(reason string) *Error {
	return &Error{
		Type:    ErrInvalidClaims,
		Message: fmt.Sprintf("Invalid claims. Reason: %v", reason),
	}
}
