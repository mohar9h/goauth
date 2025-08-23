// Package goauth provides a secure, flexible authentication system
// for Go applications with support for personal access tokens.
//
// Example:
//
//	client, err := goauth.NewClient(
//	    goauth.WithSigningKey("your-secret-key"),
//	    goauth.WithStorage(database),
//	)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	token, err := client.CreateToken(&goauth.TokenOptions{
//	    UserID: 123,
//	    Abilities: []string{"read:posts", "write:comments"},
//	})
package goauth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/mohar9h/goauth/config"
	"github.com/mohar9h/goauth/internal/auth"
	"github.com/mohar9h/goauth/internal/entity"
	"github.com/mohar9h/goauth/internal/storage"
	"github.com/mohar9h/goauth/internal/utils"
	"gorm.io/gorm"
)

// Client represents the main authentication client
type Client struct {
	config  *config.Config
	storage storage.Driver
}

// Option is a functional option for configuring the client
type Option func(*Client) error

// WithSigningKey sets the signing key for token generation
func WithSigningKey(key string) Option {
	return func(c *Client) error {
		if key == "" {
			return fmt.Errorf("signing key cannot be empty")
		}
		c.config.SigningKey = key
		return nil
	}
}

// WithTokenLength sets the length of generated tokens
func WithTokenLength(length int) Option {
	return func(c *Client) error {
		if length < 16 {
			return fmt.Errorf("token length must be at least 16 characters")
		}
		c.config.TokenLength = length
		return nil
	}
}

// WithTokenExpiration sets the token expiration duration
func WithTokenExpiration(duration time.Duration) Option {
	return func(c *Client) error {
		c.config.ExpireAt = duration
		return nil
	}
}

// WithStorage sets the storage driver
func WithStorage(driver storage.Driver) Option {
	return func(c *Client) error {
		if driver == nil {
			return fmt.Errorf("storage driver cannot be nil")
		}
		c.storage = driver
		return nil
	}
}

// WithGormStorage sets up GORM-based storage
func WithGormStorage(db *gorm.DB) Option {
	return func(c *Client) error {
		if db == nil {
			return fmt.Errorf("database connection cannot be nil")
		}
		c.storage = storage.NewGormDriver(db)
		return nil
	}
}

// WithMemoryStorage sets up in-memory storage (for testing)
func WithMemoryStorage() Option {
	return func(c *Client) error {
		c.storage = storage.NewMemoryDriver()
		return nil
	}
}

// NewClient creates a new authentication client with the given options
func NewClient(opts ...Option) (*Client, error) {
	// Generate secure default signing key if not provided
	defaultKey := generateSecureKey()

	client := &Client{
		config: &config.Config{
			TokenLength:      32, // Increased from 20 for better security
			TokenPrefix:      "",
			ExpireAt:         24 * time.Hour, // Default 24 hour expiration
			SigningMethod:    "HS256",
			SigningKey:       defaultKey,
			AbilityDelimiter: ":",
		},
		storage: storage.NewMemoryDriver(), // Default to memory storage
	}

	// Apply options
	for _, opt := range opts {
		if err := opt(client); err != nil {
			return nil, fmt.Errorf("failed to apply option: %w", err)
		}
	}

	// Validate configuration
	if err := client.config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return client, nil
}

// CreateToken generates a new token using the client's configuration
func (c *Client) CreateToken(ctx context.Context, opts *TokenOptions) (string, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	if opts == nil {
		return "", fmt.Errorf("token options cannot be nil")
	}

	if opts.UserId <= 0 {
		return "", fmt.Errorf("user ID must be positive")
	}

	// Create auth options with client config
	authOpts := &auth.TokenOptions{
		UserId:    opts.UserId,
		Name:      opts.Name,
		Abilities: opts.Abilities,
		Config:    c.config,
	}

	// Use client's storage
	authOpts.Config.Storage = c.storage

	// Check for context cancellation
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
	}

	return auth.CreateToken(authOpts)
}

// ValidateToken checks if the given token is valid and returns token info
func (c *Client) ValidateToken(ctx context.Context, raw string) (*entity.PersonalAccessToken, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	if raw == "" {
		return nil, fmt.Errorf("token cannot be empty")
	}

	// Check for context cancellation
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	return auth.ValidateToken(raw, c.config)
}

// RevokeToken removes a token from storage
func (c *Client) RevokeToken(ctx context.Context, raw string) error {
	if ctx == nil {
		ctx = context.Background()
	}

	if raw == "" {
		return fmt.Errorf("token cannot be empty")
	}

	// Check for context cancellation
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	return auth.RevokeToken(raw, c.config)
}

// GetTokenInfo retrieves token information without validation
func (c *Client) GetTokenInfo(ctx context.Context, raw string) (*entity.PersonalAccessToken, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	if raw == "" {
		return nil, fmt.Errorf("token cannot be empty")
	}

	// Check for context cancellation
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	// Extract token hash without validation
	parts := strings.Split(raw, "|")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid token format")
	}

	hashed := utils.HashToken(parts[1])
	return c.storage.FindByHash(hashed)
}

// generateSecureKey generates a cryptographically secure signing key
func generateSecureKey() string {
	// Try to get from environment first
	if key := os.Getenv("GOAUTH_SIGNING_KEY"); key != "" {
		return key
	}

	// Generate a secure random key
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		// Fallback to timestamp-based key if crypto/rand fails
		key = []byte(fmt.Sprintf("fallback-key-%d", time.Now().UnixNano()))
	}

	return base64.StdEncoding.EncodeToString(key)
}

// Legacy compatibility types and functions
type TokenOptions = auth.TokenOptions
type TokenResult = auth.Result
type Config = config.Config
type PersonalAccessToken = entity.PersonalAccessToken

var (
	legacyClient *Client
	legacyOnce   sync.Once
)

// getLegacyClient returns a singleton client for legacy functions
func getLegacyClient() (*Client, error) {
	var err error
	legacyOnce.Do(func() {
		legacyClient, err = NewClient()
	})
	return legacyClient, err
}

// Legacy functions for backward compatibility
// Deprecated: Use NewClient() and client methods instead
func CreateToken(opts *TokenOptions) (string, error) {
	client, err := getLegacyClient()
	if err != nil {
		return "", err
	}
	return client.CreateToken(context.Background(), opts)
}

// Deprecated: Use NewClient() and client methods instead
func ValidateToken(raw string, cfg *config.Config) (*entity.PersonalAccessToken, error) {
	client, err := getLegacyClient()
	if err != nil {
		return nil, err
	}
	return client.ValidateToken(context.Background(), raw)
}

// Deprecated: Use NewClient() and client methods instead
func RevokeToken(raw string, cfg *config.Config) error {
	client, err := getLegacyClient()
	if err != nil {
		return err
	}
	return client.RevokeToken(context.Background(), raw)
}

// Deprecated: Use NewClient() and client methods instead
func SetupGorm(db *gorm.DB) *config.Config {
	cfg := config.DefaultConfig()
	cfg.Storage = storage.NewGormDriver(db)
	return cfg
}
