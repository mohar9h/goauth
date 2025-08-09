// Package goauth package goauth.go
package goauth

import (
	"github.com/mohar9h/goauth/config"
	"github.com/mohar9h/goauth/internal/auth"
	"github.com/mohar9h/goauth/internal/entity"
	"github.com/mohar9h/goauth/internal/storage"
	"gorm.io/gorm"
)

type TokenOptions = auth.TokenOptions
type TokenResult = auth.Result

type PersonalAccessToken = entity.PersonalAccessToken

// CreateToken generates a new token using given options and configuration.
func CreateToken(opts *TokenOptions) (string, error) {
	return auth.CreateToken(opts)
}

// ValidateToken checks if the given token string is valid and returns token info.
func ValidateToken(raw string, cfg *config.Config) (*entity.PersonalAccessToken, error) {
	return auth.ValidateToken(raw, cfg)
}

// RevokeToken removes a token from the store by its raw value.
func RevokeToken(raw string, cfg *config.Config) error {
	return auth.RevokeToken(raw, cfg)
}

// SetupGorm sets up a GORM driver with default config.
func SetupGorm(db *gorm.DB) *config.Config {
	cfg := config.DefaultConfig()
	cfg.Storage = storage.NewGormDriver(db)
	return cfg
}
