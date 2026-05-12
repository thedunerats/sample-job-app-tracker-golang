package main

import (
	"testing"
	"time"
)

func TestHashPassword(t *testing.T) {
	password := "testpassword123"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	if hash == "" {
		t.Error("Hash should not be empty")
	}

	if hash == password {
		t.Error("Hash should not equal plain password")
	}
}

func TestCheckPasswordHash(t *testing.T) {
	password := "testpassword123"
	wrongPassword := "wrongpassword"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	// Test correct password
	if !CheckPasswordHash(password, hash) {
		t.Error("CheckPasswordHash should return true for correct password")
	}

	// Test wrong password
	if CheckPasswordHash(wrongPassword, hash) {
		t.Error("CheckPasswordHash should return false for wrong password")
	}
}

func TestGenerateToken(t *testing.T) {
	userID := 1
	email := "test@example.com"

	token, err := GenerateToken(userID, email)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	if token == "" {
		t.Error("Token should not be empty")
	}
}

func TestValidateToken(t *testing.T) {
	userID := 1
	email := "test@example.com"

	// Generate a valid token
	token, err := GenerateToken(userID, email)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	// Validate the token
	claims, err := ValidateToken(token)
	if err != nil {
		t.Fatalf("Failed to validate token: %v", err)
	}

	if claims.UserID != userID {
		t.Errorf("Expected UserID %d, got %d", userID, claims.UserID)
	}

	if claims.Email != email {
		t.Errorf("Expected Email %s, got %s", email, claims.Email)
	}

	// Test invalid token
	invalidToken := "invalid.token.string"
	_, err = ValidateToken(invalidToken)
	if err == nil {
		t.Error("ValidateToken should return error for invalid token")
	}
}

func TestTokenExpiration(t *testing.T) {
	// This test would need to manipulate time or wait for expiration
	// For now, we just verify that the expiration is set correctly
	userID := 1
	email := "test@example.com"

	token, err := GenerateToken(userID, email)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	claims, err := ValidateToken(token)
	if err != nil {
		t.Fatalf("Failed to validate token: %v", err)
	}

	// Check that expiration is set to approximately 24 hours in the future
	expectedExpiration := time.Now().Add(24 * time.Hour)
	actualExpiration := claims.ExpiresAt.Time

	diff := actualExpiration.Sub(expectedExpiration)
	if diff < 0 {
		diff = -diff
	}

	// Allow 1 minute tolerance
	if diff > time.Minute {
		t.Errorf("Expiration time is not set correctly. Expected ~%v, got %v", expectedExpiration, actualExpiration)
	}
}

func TestMultipleTokens(t *testing.T) {
	// Test that different users get different tokens
	token1, _ := GenerateToken(1, "user1@example.com")
	token2, _ := GenerateToken(2, "user2@example.com")

	if token1 == token2 {
		t.Error("Different users should have different tokens")
	}

	// Validate both tokens
	claims1, err := ValidateToken(token1)
	if err != nil {
		t.Fatalf("Failed to validate token1: %v", err)
	}

	claims2, err := ValidateToken(token2)
	if err != nil {
		t.Fatalf("Failed to validate token2: %v", err)
	}

	if claims1.UserID == claims2.UserID {
		t.Error("Different tokens should have different user IDs")
	}
}
