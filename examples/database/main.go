package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/mohar9h/goauth"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	// Setup SQLite database (replace with your database)
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto-migrate the token table
	err = db.AutoMigrate(&goauth.PersonalAccessToken{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Create client with database storage
	client, err := goauth.NewClient(
		goauth.WithSigningKey("your-production-secret-key"),
		goauth.WithGormStorage(db),
		goauth.WithTokenLength(64),
		goauth.WithTokenExpiration(7*24*time.Hour), // 7 days
	)
	if err != nil {
		log.Fatal("Failed to create client:", err)
	}

	// Create multiple tokens for different users
	users := []struct {
		id        int64
		name      string
		abilities []string
	}{
		{1, "Admin Token", []string{"*"}},
		{2, "Read-Only Token", []string{"read:posts", "read:comments"}},
		{3, "Moderator Token", []string{"read:posts", "write:comments", "delete:comments"}},
	}

	var tokens []string
	for _, user := range users {
		token, err := client.CreateToken(context.Background(), &goauth.TokenOptions{
			UserId:    user.id,
			Name:      &user.name,
			Abilities: user.abilities,
		})
		if err != nil {
			log.Fatalf("Failed to create token for user %d: %v", user.id, err)
		}
		tokens = append(tokens, token)
		fmt.Printf("Created token for user %d: %s\n", user.id, token)
	}

	// Validate all tokens
	for i, token := range tokens {
		tokenInfo, err := client.ValidateToken(context.Background(), token)
		if err != nil {
			log.Fatalf("Failed to validate token %d: %v", i+1, err)
		}
		fmt.Printf("Token %d belongs to user %d with abilities: %s\n",
			i+1, tokenInfo.UserId, tokenInfo.Abilities)
	}

	// Demonstrate token expiration
	fmt.Println("\n--- Testing Token Expiration ---")

	// Create a token with very short expiration
	shortExpiryClient, err := goauth.NewClient(
		goauth.WithSigningKey("test-key"),
		goauth.WithGormStorage(db),
		goauth.WithTokenExpiration(1*time.Millisecond),
	)
	if err != nil {
		log.Fatal("Failed to create short expiry client:", err)
	}

	expiringToken, err := shortExpiryClient.CreateToken(context.Background(), &goauth.TokenOptions{
		UserId:    999,
		Abilities: []string{"test"},
	})
	if err != nil {
		log.Fatal("Failed to create expiring token:", err)
	}

	// Token should be valid immediately
	_, err = shortExpiryClient.ValidateToken(context.Background(), expiringToken)
	if err != nil {
		log.Fatal("Token should be valid immediately:", err)
	}
	fmt.Println("Token is valid immediately after creation")

	// Wait for expiration
	time.Sleep(10 * time.Millisecond)

	// Token should be expired
	_, err = shortExpiryClient.ValidateToken(context.Background(), expiringToken)
	if err != nil {
		fmt.Println("Token expired as expected:", err)
	} else {
		fmt.Println("Token should have expired but is still valid")
	}

	// Demonstrate context cancellation
	fmt.Println("\n--- Testing Context Cancellation ---")

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	_, err = client.CreateToken(ctx, &goauth.TokenOptions{
		UserId:    888,
		Abilities: []string{"test"},
	})
	if err != nil {
		fmt.Println("Context cancellation worked as expected:", err)
	}

	// Clean up - revoke all tokens
	fmt.Println("\n--- Cleaning Up ---")
	for i, token := range tokens {
		err := client.RevokeToken(context.Background(), token)
		if err != nil {
			log.Fatalf("Failed to revoke token %d: %v", i+1, err)
		}
		fmt.Printf("Revoked token %d\n", i+1)
	}

	// Verify tokens are revoked
	for i, token := range tokens {
		_, err := client.ValidateToken(context.Background(), token)
		if err != nil {
			fmt.Printf("Token %d is properly revoked: %v\n", i+1, err)
		} else {
			fmt.Printf("Token %d should be revoked but is still valid\n", i+1)
		}
	}

	fmt.Println("\nExample completed successfully!")
}
