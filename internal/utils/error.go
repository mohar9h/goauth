// Package utils Package internal/utils/error.go
package utils

import "errors"

var (
	ErrTokenExpired       = errors.New("token expired")
	ErrTokenNotFound      = errors.New("token not found")
	ErrTokenInvalidFormat = errors.New("invalid token format")
	ErrSessionInactive    = errors.New("session inactive")
	ErrStorageFailure     = errors.New("storage operation failed")
)
