# GoAuth - Secure Authentication Package

A secure, flexible authentication system for Go applications with support for personal access tokens.

## Features

- 🔐 **Secure Token Generation**: Cryptographically secure random tokens with configurable length
- 🗄️ **Multiple Storage Backends**: Support for GORM (PostgreSQL, MySQL, SQLite) and in-memory storage
- ⏰ **Token Expiration**: Configurable token expiration with automatic cleanup
- 🔄 **Context Support**: Full context.Context support for cancellation and timeouts
- 🛡️ **Security First**: Environment-based configuration, secure defaults, and proper error handling
- 🧪 **Testable**: Clean interfaces and dependency injection for easy testing
- 📊 **Analytics**: Token usage tracking with last-used timestamps

## Installation

```bash
go get github.com/mohar9h/goauth
```

## Quick Start

### Basic Usage

```go
package main

import (
    "context"
    "log"
    
    "github.com/mohar9h/goauth"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

func main() {
    // Setup database connection
    db, err := gorm.Open(postgres.Open("dsn"), &gorm.Config{})
    if err != nil {
        log.Fatal(err)
    }
    
    // Create client with GORM storage
    client, err := goauth.NewClient(
        goauth.WithSigningKey("your-secret-key"),
        goauth.WithGormStorage(db),
        goauth.WithTokenExpiration(24 * time.Hour),
    )
    if err != nil {
        log.Fatal(err)
    }
    
    // Create a token
    token, err := client.CreateToken(context.Background(), &goauth.TokenOptions{
        UserId:     123,
        Name:       stringPtr("API Token"),
        Abilities:  []string{"read:posts", "write:comments"},
    })
    if err != nil {
        log.Fatal(err)
    }
    
    log.Printf("Generated token: %s", token)
    
    // Validate token
    tokenInfo, err := client.ValidateToken(context.Background(), token)
    if err != nil {
        log.Fatal(err)
    }
    
    log.Printf("Token belongs to user: %d", tokenInfo.UserId)
    
    // Revoke token
    err = client.RevokeToken(context.Background(), token)
    if err != nil {
        log.Fatal(err)
    }
}

func stringPtr(s string) *string {
    return &s
}
```

### Environment Configuration

Set the signing key via environment variable for production:

```bash
export GOAUTH_SIGNING_KEY="your-super-secret-key-here"
```

### Advanced Configuration

```go
client, err := goauth.NewClient(
    goauth.WithSigningKey("your-secret-key"),
    goauth.WithTokenLength(64),                    // 64-character tokens
    goauth.WithTokenExpiration(7 * 24 * time.Hour), // 7 days
    goauth.WithGormStorage(db),
)
```

### In-Memory Storage (for testing)

```go
client, err := goauth.NewClient(
    goauth.WithSigningKey("test-key"),
    goauth.WithMemoryStorage(),
)
```

## API Reference

### Client

#### `NewClient(opts ...Option) (*Client, error)`

Creates a new authentication client with the given options.

#### `client.CreateToken(ctx context.Context, opts *TokenOptions) (string, error)`

Generates a new personal access token.

#### `client.ValidateToken(ctx context.Context, raw string) (*PersonalAccessToken, error)`

Validates a token and returns its information.

#### `client.RevokeToken(ctx context.Context, raw string) error`

Revokes a token, making it invalid.

#### `client.GetTokenInfo(ctx context.Context, raw string) (*PersonalAccessToken, error)`

Retrieves token information without validation.

### Options

#### `WithSigningKey(key string) Option`

Sets the signing key for token generation. Required for security.

#### `WithTokenLength(length int) Option`

Sets the length of generated tokens (minimum 16 characters).

#### `WithTokenExpiration(duration time.Duration) Option`

Sets the token expiration duration.

#### `WithGormStorage(db *gorm.DB) Option`

Sets up GORM-based storage for tokens.

#### `WithMemoryStorage() Option`

Sets up in-memory storage (useful for testing).

### Types

#### `TokenOptions`

```go
type TokenOptions struct {
    UserId    int64     // User ID (required)
    Name      *string   // Token name (optional)
    Abilities []string  // Token abilities/permissions
}
```

#### `PersonalAccessToken`

```go
type PersonalAccessToken struct {
    ID         int64      `gorm:"primaryKey;autoIncrement"`
    UserId     int64      `gorm:"index"`
    Token      string     `gorm:"index;size:100"`
    Name       *string    `gorm:"size:100"`
    Abilities  string     `gorm:"type:text"`
    CreatedAt  time.Time  `gorm:"autoCreateTime"`
    ExpiresAt  *time.Time `gorm:"index"`
    LastUsedAt *time.Time
}
```

## Database Schema

The package automatically creates the following table:

```sql
CREATE TABLE personal_access_tokens (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    token VARCHAR(100) NOT NULL,
    name VARCHAR(100),
    abilities TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMP,
    last_used_at TIMESTAMP,
    
    INDEX idx_user_id (user_id),
    INDEX idx_token (token),
    INDEX idx_expires_at (expires_at)
);
```

## Security Considerations

1. **Signing Key**: Always use a strong, randomly generated signing key in production
2. **Environment Variables**: Store sensitive configuration in environment variables
3. **Token Length**: Use at least 32 characters for token length in production
4. **Expiration**: Set reasonable token expiration times
5. **HTTPS**: Always use HTTPS in production to protect tokens in transit
6. **Storage**: Use secure database connections and proper access controls

## Testing

```go
func TestTokenCreation(t *testing.T) {
    client, err := goauth.NewClient(
        goauth.WithSigningKey("test-key"),
        goauth.WithMemoryStorage(),
    )
    require.NoError(t, err)
    
    token, err := client.CreateToken(context.Background(), &goauth.TokenOptions{
        UserId:    123,
        Abilities: []string{"read:posts"},
    })
    require.NoError(t, err)
    require.NotEmpty(t, token)
    
    // Validate token
    tokenInfo, err := client.ValidateToken(context.Background(), token)
    require.NoError(t, err)
    require.Equal(t, int64(123), tokenInfo.UserId)
}
```

## Migration from Legacy API

The package maintains backward compatibility with the legacy API:

```go
// Old way (deprecated)
token, err := goauth.CreateToken(&goauth.TokenOptions{...})

// New way (recommended)
client, err := goauth.NewClient(goauth.WithSigningKey("key"))
token, err := client.CreateToken(context.Background(), &goauth.TokenOptions{...})
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## License

MIT License - see LICENSE file for details.
