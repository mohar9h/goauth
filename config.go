package goauth

import (
	"crypto/rsa"
	"errors"
	"time"
)

// Config holds the global settings for the auth package.
type Config struct {
	TokenLength      int             // Length of random tokens (e.g., 32)
	TokenPrefix      string          // Prefix for random tokens (e.g., "pk_")
	ExpireAt         time.Duration   // Token TTL (0 = unlimited)
	SigningKey       string          // For HMAC JWT (HS256)
	SigningMethod    string          // "HS256", "RS256"
	PrivateKey       *rsa.PrivateKey // For RSA signing (optional)
	PublicKey        *rsa.PublicKey  // For RSA verification (optional)
	Storage          storage.Driver  // Optional: for random tokens
	AbilityDelimiter string          // e.g., ":" for "read:posts"
}

// Validate checks if the config is minimally valid.
func (c *Config) Validate() error {
	if c.SigningMethod != "HS256" && c.SigningMethod != "RS256" {
		return errors.New("unsupported signing method")
	}
	if c.SigningMethod == "HS256" && c.SigningKey == "" {
		return errors.New("missing HMAC signing key")
	}
	if c.SigningMethod == "RS256" && (c.PrivateKey == nil || c.PublicKey == nil) {
		return errors.New("missing RSA key pair")
	}
	if c.TokenLength < 16 {
		return errors.New("auth length too short")
	}
	return nil
}

// DefaultConfig returns a default config.
func DefaultConfig() *Config {
	return &Config{
		TokenLength:      20,
		TokenPrefix:      "",
		ExpireAt:         0,
		SigningMethod:    "HS256",
		SigningKey:       "test-key",
		AbilityDelimiter: ":",
		Storage:          storage2.NewMemoryDriver(),
	}
}

func (c *Config) ApplyDefaults() {
	def := DefaultConfig()

	if c.TokenLength == 0 {
		c.TokenLength = def.TokenLength
	}
	if c.TokenPrefix == "" {
		c.TokenPrefix = def.TokenPrefix
	}
	if c.ExpireAt == 0 {
		c.ExpireAt = def.ExpireAt
	}
	if c.SigningMethod == "" {
		c.SigningMethod = def.SigningMethod
	}
	if c.SigningKey == "" {
		c.SigningKey = def.SigningKey
	}
	if c.AbilityDelimiter == "" {
		c.AbilityDelimiter = def.AbilityDelimiter
	}
	if c.Storage == nil {
		c.Storage = storage.NewMemoryDriver()
	}
}
