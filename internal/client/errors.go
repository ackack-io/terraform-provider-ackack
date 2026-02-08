// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"
	"net/http"
)

// APIError represents an error response from the ackack.io API.
type APIError struct {
	StatusCode int
	Message    string
	ErrorField string
}

func (e *APIError) Error() string {
	if e.ErrorField != "" {
		return fmt.Sprintf("API error (status %d): %s - %s", e.StatusCode, e.Message, e.ErrorField)
	}
	return fmt.Sprintf("API error (status %d): %s", e.StatusCode, e.Message)
}

// IsNotFoundError returns true if the error is a 404 Not Found error.
func IsNotFoundError(err error) bool {
	if apiErr, ok := err.(*APIError); ok {
		return apiErr.StatusCode == http.StatusNotFound
	}
	return false
}

// IsRateLimitError returns true if the error is a 429 Too Many Requests error.
func IsRateLimitError(err error) bool {
	if apiErr, ok := err.(*APIError); ok {
		return apiErr.StatusCode == http.StatusTooManyRequests
	}
	return false
}

// IsUnauthorizedError returns true if the error is a 401 Unauthorized error.
func IsUnauthorizedError(err error) bool {
	if apiErr, ok := err.(*APIError); ok {
		return apiErr.StatusCode == http.StatusUnauthorized
	}
	return false
}

// IsForbiddenError returns true if the error is a 403 Forbidden error.
func IsForbiddenError(err error) bool {
	if apiErr, ok := err.(*APIError); ok {
		return apiErr.StatusCode == http.StatusForbidden
	}
	return false
}

// IsBadRequestError returns true if the error is a 400 Bad Request error.
func IsBadRequestError(err error) bool {
	if apiErr, ok := err.(*APIError); ok {
		return apiErr.StatusCode == http.StatusBadRequest
	}
	return false
}
