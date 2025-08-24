package auth_test

import (
	"context"
	"testing"
	"time"

	"github.com/mohar9h/goauth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name    string
		opts    []goauth.Option
		wantErr bool
	}{
		{
			name:    "default client",
			opts:    []goauth.Option{},
			wantErr: false,
		},
		{
			name: "client with signing key",
			opts: []goauth.Option{
				goauth.WithSigningKey("test-key-123"),
			},
			wantErr: false,
		},
		{
			name: "client with memory storage",
			opts: []goauth.Option{
				goauth.WithSigningKey("test-key-123"),
				goauth.WithMemoryStorage(),
			},
			wantErr: false,
		},
		{
			name: "client with token length",
			opts: []goauth.Option{
				goauth.WithSigningKey("test-key-123"),
				goauth.WithTokenLength(64),
			},
			wantErr: false,
		},
		{
			name: "client with expiration",
			opts: []goauth.Option{
				goauth.WithSigningKey("test-key-123"),
				goauth.WithTokenExpiration(1 * time.Hour),
			},
			wantErr: false,
		},
		{
			name: "invalid token length",
			opts: []goauth.Option{
				goauth.WithSigningKey("test-key-123"),
				goauth.WithTokenLength(10), // Too short
			},
			wantErr: true,
		},
		{
			name: "empty signing key",
			opts: []goauth.Option{
				goauth.WithSigningKey(""),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := goauth.NewClient(tt.opts...)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.NotNil(t, client)
		})
	}
}

func TestCreateToken(t *testing.T) {
	client, err := goauth.NewClient(
		goauth.WithSigningKey("test-key-123"),
		goauth.WithMemoryStorage(),
	)
	require.NoError(t, err)

	tests := []struct {
		name    string
		opts    *goauth.TokenOptions
		wantErr bool
	}{
		{
			name: "valid token creation",
			opts: &goauth.TokenOptions{
				UserId:    123,
				Name:      stringPtr("Test Token"),
				Abilities: []string{"read:posts", "write:comments"},
			},
			wantErr: false,
		},
		{
			name: "token without name",
			opts: &goauth.TokenOptions{
				UserId:    456,
				Abilities: []string{"read:posts"},
			},
			wantErr: false,
		},
		{
			name: "token with wildcard abilities",
			opts: &goauth.TokenOptions{
				UserId:    789,
				Abilities: []string{"*"},
			},
			wantErr: false,
		},
		{
			name:    "nil options",
			opts:    nil,
			wantErr: true,
		},
		{
			name: "invalid user ID",
			opts: &goauth.TokenOptions{
				UserId:    0,
				Abilities: []string{"read:posts"},
			},
			wantErr: true,
		},
		{
			name: "negative user ID",
			opts: &goauth.TokenOptions{
				UserId:    -1,
				Abilities: []string{"read:posts"},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := client.CreateToken(context.Background(), tt.opts)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.NotEmpty(t, token)
			assert.Contains(t, token, "|")
		})
	}
}

func TestValidateToken(t *testing.T) {
	client, err := goauth.NewClient(
		goauth.WithSigningKey("test-key-123"),
		goauth.WithMemoryStorage(),
	)
	require.NoError(t, err)

	// Create a valid token first
	token, err := client.CreateToken(context.Background(), &goauth.TokenOptions{
		UserId:    123,
		Name:      stringPtr("Test Token"),
		Abilities: []string{"read:posts"},
	})
	require.NoError(t, err)

	tests := []struct {
		name    string
		token   string
		wantErr bool
	}{
		{
			name:    "valid token",
			token:   token,
			wantErr: false,
		},
		{
			name:    "valid token with bearer prefix",
			token:   "Bearer " + token,
			wantErr: false,
		},
		{
			name:    "empty token",
			token:   "",
			wantErr: true,
		},
		{
			name:    "invalid token format",
			token:   "invalid-token",
			wantErr: true,
		},
		{
			name:    "non-existent token",
			token:   "1|nonexistenttoken",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokenInfo, err := client.ValidateToken(context.Background(), tt.token)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.NotNil(t, tokenInfo)
			assert.Equal(t, int64(123), tokenInfo.UserId)
			assert.Equal(t, "read:posts", tokenInfo.Abilities)
		})
	}
}

func TestRevokeToken(t *testing.T) {
	client, err := goauth.NewClient(
		goauth.WithSigningKey("test-key-123"),
		goauth.WithMemoryStorage(),
	)
	require.NoError(t, err)

	// Create a valid token first
	token, err := client.CreateToken(context.Background(), &goauth.TokenOptions{
		UserId:    123,
		Abilities: []string{"read:posts"},
	})
	require.NoError(t, err)

	// Verify token is valid
	_, err = client.ValidateToken(context.Background(), token)
	require.NoError(t, err)

	// Revoke token
	err = client.RevokeToken(context.Background(), token)
	require.NoError(t, err)

	// Verify token is no longer valid
	_, err = client.ValidateToken(context.Background(), token)
	assert.Error(t, err)
}

func TestGetTokenInfo(t *testing.T) {
	client, err := goauth.NewClient(
		goauth.WithSigningKey("test-key-123"),
		goauth.WithMemoryStorage(),
	)
	require.NoError(t, err)

	// Create a valid token first
	token, err := client.CreateToken(context.Background(), &goauth.TokenOptions{
		UserId:    123,
		Name:      stringPtr("Test Token"),
		Abilities: []string{"read:posts"},
	})
	require.NoError(t, err)

	// Get token info
	tokenInfo, err := client.GetTokenInfo(context.Background(), token)
	require.NoError(t, err)
	assert.NotNil(t, tokenInfo)
	assert.Equal(t, int64(123), tokenInfo.UserId)
	assert.Equal(t, "read:posts", tokenInfo.Abilities)
}

func TestTokenExpiration(t *testing.T) {
	client, err := goauth.NewClient(
		goauth.WithSigningKey("test-key-123"),
		goauth.WithMemoryStorage(),
		goauth.WithTokenExpiration(1*time.Millisecond), // Very short expiration
	)
	require.NoError(t, err)

	// Create a token
	token, err := client.CreateToken(context.Background(), &goauth.TokenOptions{
		UserId:    123,
		Abilities: []string{"read:posts"},
	})
	require.NoError(t, err)

	// Token should be valid immediately
	_, err = client.ValidateToken(context.Background(), token)
	require.NoError(t, err)

	// Wait for expiration
	time.Sleep(10 * time.Millisecond)

	// Token should be expired
	_, err = client.ValidateToken(context.Background(), token)
	assert.Error(t, err)
}

func TestContextCancellation(t *testing.T) {
	client, err := goauth.NewClient(
		goauth.WithSigningKey("test-key-123"),
		goauth.WithMemoryStorage(),
	)
	require.NoError(t, err)

	// Create a cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	// Try to create token with cancelled context
	_, err = client.CreateToken(ctx, &goauth.TokenOptions{
		UserId:    123,
		Abilities: []string{"read:posts"},
	})
	assert.Error(t, err)
	assert.Equal(t, context.Canceled, err)
}

// Helper function to create string pointers
func stringPtr(s string) *string {
	return &s
}
