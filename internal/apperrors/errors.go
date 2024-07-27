// Package apperrors provides common error definitions used across the application.
package apperrors

import "errors"

// ErrInvalidRequestData is an error that indicates that the request data is invalid.
var ErrInvalidRequestData = errors.New("invalid request data")
