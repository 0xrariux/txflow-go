package txflow

import "github.com/0xrariux/txflow-go/sdkerr"

type APIError = sdkerr.APIError

var (
	ErrPendingConfirmation = sdkerr.ErrPendingConfirmation
	ErrRateLimited         = sdkerr.ErrRateLimited
	ErrUnauthorized        = sdkerr.ErrUnauthorized
	ErrInvalidRequest      = sdkerr.ErrInvalidRequest
	ErrMalformedResponse   = sdkerr.ErrMalformedResponse
	ErrResponseTooLarge    = sdkerr.ErrResponseTooLarge
	ErrNotFound            = sdkerr.ErrNotFound
)
