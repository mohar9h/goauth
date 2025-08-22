// Package auth internal/auth/revoke.go
package auth

import (
	"fmt"
	"github.com/mohar9h/goauth/config"
	"strconv"
)

func RevokeToken(raw string, cfg *config.Config) error {
	token, err := ValidateToken(raw, cfg)
	if err != nil {
		return err
	}

	if err = cfg.Storage.RevokeToken(strconv.FormatInt(token.ID, 10)); err != nil {
		return fmt.Errorf("failed to revoke token: %w", err)
	}

	return nil
}
