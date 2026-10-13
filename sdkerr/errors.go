// Package sdkerr defines errors shared by services and transports.
package sdkerr

import (
	"errors"
	"fmt"
)

var (
	ErrPendingConfirmation = errors.New("txflow: protocol pending TxFlow confirmation")
	ErrRateLimited         = errors.New("txflow: rate limited")
	ErrUnauthorized        = errors.New("txflow: unauthorized")
	ErrInvalidRequest      = errors.New("txflow: invalid request")
	ErrMalformedResponse   = errors.New("txflow: malformed response")
	ErrResponseTooLarge    = errors.New("txflow: response too large")
	ErrNotFound            = errors.New("txflow: not found")
	ErrSequenceGap         = errors.New("txflow: sequence gap; snapshot required")
	ErrStaleUpdate         = errors.New("txflow: stale or duplicate update")
	ErrBookUnavailable     = errors.New("txflow: book unavailable")
)

// APIError describes an HTTP failure. Code and Message require a future
// confirmed protocol decoder; arbitrary server bodies are not echoed.
type APIError struct {
	Code       string
	Message    string
	StatusCode int
	RequestID  string
}

func (e *APIError) Error() string { return fmt.Sprintf("txflow: HTTP status %d", e.StatusCode) }
func (e *APIError) Is(target error) bool {
	switch e.StatusCode {
	case 429:
		return target == ErrRateLimited
	case 401, 403:
		return target == ErrUnauthorized
	case 404:
		return target == ErrNotFound
	case 400, 422:
		return target == ErrInvalidRequest
	}
	return false
}
