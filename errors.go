package fimage

import (
	"errors"
	"fmt"
)

// Common errors returned by the SDK.
var (
	// ErrUnauthorized is returned when the API token is invalid or missing.
	ErrUnauthorized = errors.New("unauthorized: invalid or missing API token")

	// ErrNotFound is returned when the requested resource is not found.
	ErrNotFound = errors.New("not found: the requested resource does not exist")

	// ErrBadRequest is returned when the request is invalid.
	ErrBadRequest = errors.New("bad request: invalid request parameters")

	// ErrForbidden is returned when access to a resource is forbidden.
	ErrForbidden = errors.New("forbidden: access denied")

	// ErrQuotaExceeded is returned when storage quota is exceeded.
	ErrQuotaExceeded = errors.New("quota exceeded: storage limit reached")

	// ErrFileTooLarge is returned when the uploaded file exceeds the size limit.
	ErrFileTooLarge = errors.New("file too large: exceeds maximum file size")

	// ErrInvalidFormat is returned when the file format is not allowed.
	ErrInvalidFormat = errors.New("invalid format: file type not allowed")
)

// APIError represents an error returned by the F-Image API.
type APIError struct {
	// StatusCode is the HTTP status code.
	StatusCode int

	// Message is the error message from the API.
	Message string
}

// Error implements the error interface.
func (e *APIError) Error() string {
	return fmt.Sprintf("f-image API error (status %d): %s", e.StatusCode, e.Message)
}

// IsNotFound returns true if the error is a not found error.
func IsNotFound(err error) bool {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.StatusCode == 404
	}
	return errors.Is(err, ErrNotFound)
}

// IsUnauthorized returns true if the error is an unauthorized error.
func IsUnauthorized(err error) bool {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.StatusCode == 401
	}
	return errors.Is(err, ErrUnauthorized)
}

// IsForbidden returns true if the error is a forbidden error.
func IsForbidden(err error) bool {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.StatusCode == 403
	}
	return errors.Is(err, ErrForbidden)
}

// IsBadRequest returns true if the error is a bad request error.
func IsBadRequest(err error) bool {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.StatusCode == 400
	}
	return errors.Is(err, ErrBadRequest)
}

// IsQuotaExceeded returns true if the error is a quota exceeded error.
func IsQuotaExceeded(err error) bool {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.StatusCode == 402 || apiErr.StatusCode == 413
	}
	return errors.Is(err, ErrQuotaExceeded)
}
