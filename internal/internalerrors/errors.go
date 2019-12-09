package internalerrors

import "errors"

var (
	// ErrInvalidConfiguration will be used when a configuration is invalid.
	ErrInvalidConfiguration = errors.New("configuration is invalid")
)
