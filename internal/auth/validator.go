// Package auth internal/auth/validator.go
package auth

import (
	"errors"
	"strings"
	"time"

	"github.com/mohar9h/goauth/config"
	"github.com/mohar9h/goauth/internal/entity"
	"github.com/mohar9h/goauth/internal/utils"
)

var ErrTokenInvalid = errors.New("token is invalid or expired")

func ValidateToken(raw string, cfg *config.Config) (*entity.PersonalAccessToken, error) {

	if cfg == nil {
		cfg = config.DefaultConfig()
	}
	cfg.ApplyDefaults()

	if after, ok := strings.CutPrefix(raw, "Bearer "); ok {
		raw = after
	}

	parts := strings.Split(raw, "|")
	if len(parts) != 2 {
		return nil, ErrTokenInvalid
	}

	hashed := utils.HashToken(parts[1])

	tok, err := cfg.Storage.FindByHash(hashed)
	if err != nil {
		return nil, err
	}

	if tok.Token != hashed {
		return nil, ErrTokenInvalid
	}

	if tok.ExpiresAt != nil && time.Now().After(*tok.ExpiresAt) {
		return nil, utils.ErrTokenExpired
	}

	// Update last used time asynchronously
	go func() {
		if err := cfg.Storage.TouchLastUsed(tok.ID); err != nil {
			// Log error but don't fail validation
			// In a production environment, you might want to use a proper logger
			_ = err // Suppress unused variable warning
		}
	}()

	return tok, nil
}
