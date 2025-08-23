package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/mohar9h/goauth"
)

func main() {
	// Create a new client with secure defaults
	client, err := goauth.NewClient(
		goauth.WithSigningKey("your-secret-key-here"),
		goauth.WithMemoryStorage(), // Use in-memory storage for this example
		goauth.WithTokenExpiration(24*time.Hour),
	)
	if err != nil {
		log.Fatal("Failed to create client:", err)
	}

	// Create a token for user 123
	token, err := client.CreateToken(context.Background(), &goauth.TokenOptions{
		UserId:    123,
		Name:      stringPtr("API Token for User 123"),
		Abilities: []string{"read:posts", "write:comments", "delete:own_posts"},
	})
	if err != nil {
		log.Fatal("Failed to create token:", err)
	}

	fmt.Printf("Generated token: %s\n", token)

	// Validate the token
	tokenInfo, err := client.ValidateToken(context.Background(), token)
	if err != nil {
		log.Fatal("Failed to validate token:", err)
	}

	fmt.Printf("Token belongs to user: %d\n", tokenInfo.UserId)
	fmt.Printf("Token abilities: %s\n", tokenInfo.Abilities)
	fmt.Printf("Token created at: %s\n", tokenInfo.CreatedAt.Format(time.RFC3339))

	// Get token info without validation
	info, err := client.GetTokenInfo(context.Background(), token)
	if err != nil {
		log.Fatal("Failed to get token info:", err)
	}

	fmt.Printf("Token name: %s\n", *info.Name)

	// Revoke the token
	err = client.RevokeToken(context.Background(), token)
	if err != nil {
		log.Fatal("Failed to revoke token:", err)
	}

	fmt.Println("Token revoked successfully")

	// Try to validate the revoked token (should fail)
	_, err = client.ValidateToken(context.Background(), token)
	if err != nil {
		fmt.Printf("Token validation failed as expected: %v\n", err)
	}
}

// Helper function to create string pointers
func stringPtr(s string) *string {
	return &s
}
