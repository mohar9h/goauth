package auth

import (
	"errors"
	"github.com/mohar9h/goauth/config"
	"github.com/mohar9h/goauth/internal/entity"
	"github.com/mohar9h/goauth/internal/utils"
	"strings"
	"time"
)

var ErrTokenInvalid = errors.New("token is invalid or expired")

func ValidateToken(raw string, cfg *config.Config) (*entity.Token, error) {

	if cfg == nil {
		cfg = config.DefaultConfig()
	}
	cfg.ApplyDefaults()

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
		return nil, errors.New("token expired")
	}

	go func() {
		err = cfg.Storage.TouchLastUsed(tok.ID)
		if err != nil {

		}
	}()

	return tok, nil
}
