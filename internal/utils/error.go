// Package utils Package internal/utils/error.go
package utils

import "errors"

var (
	ErrTokenExpired            = errors.New("token expired")
	ErrTokenNotFound           = errors.New("token not found")
	ErrTokenInvalidFormat      = errors.New("invalid token format")
	ErrSigningKeyCannotBeEmpty = errors.New("signing key cannot be empty")
	ErrTokenLengthMustBe       = errors.New("token length must be at least 16 characters")
	ErrStorageDriverNil        = errors.New("storage driver cannot be nil")
	ErrDatabaseConnectionNil   = errors.New("database connection cannot be nil")
)
