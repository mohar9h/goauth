package auth

import (
	"fmt"
	"github.com/mohar9h/goauth/config"
	"github.com/mohar9h/goauth/internal/storage"
)

func CreateToken(opts *TokenOptions) (string, error) {
	if opts == nil {
		return "", fmt.Errorf("options required")
	}

	cfg := opts.Config
	if cfg == nil {
		cfg = config.DefaultConfig()
	}
	cfg.ApplyDefaults()
	cfg.Storage = storage.NewGormDriver(opts.DB)

	if err := cfg.Validate(); err != nil {
		return "", fmt.Errorf("invalid config: %w", err)
	}

	gen := NewGenerator(opts, cfg)
	result, err := gen.Create()
	if err != nil {
		return "", err
	}

	return result.PlainText, nil
}
