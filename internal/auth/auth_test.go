package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestMakeAndValidateJWT(t *testing.T) {
	// Create a valid JWT for a new user.
	userID := uuid.New()
	secret := "supersecret"
	expiresIn := 5 * time.Minute

	token, err := MakeJWT(userID, secret, expiresIn)
	if err != nil {
		t.Fatalf("Error generating JWT: %v", err)
	}

	// Validate the token using the correct secret.
	returnedUserID, err := ValidateJWT(token, secret)
	if err != nil {
		t.Fatalf("Error validating JWT: %v", err)
	}

	if returnedUserID != userID {
		t.Errorf("Expected userID %v, got %v", userID, returnedUserID)
	}
}

func TestExpiredJWT(t *testing.T) {
	// Generate a JWT that is already expired.
	userID := uuid.New()
	secret := "supersecret"
	// Set expiration in the past.
	expiresIn := -1 * time.Minute

	token, err := MakeJWT(userID, secret, expiresIn)
	if err != nil {
		t.Fatalf("Error generating JWT: %v", err)
	}

	// Validate should fail due to token expiration.
	_, err = ValidateJWT(token, secret)
	if err == nil {
		t.Fatal("Expected error when validating expired JWT, got nil")
	}
}

func TestWrongSecretJWT(t *testing.T) {
	// Generate a valid JWT with a secret.
	userID := uuid.New()
	secret := "supersecret"
	wrongSecret := "wrongsecret"
	expiresIn := 5 * time.Minute

	token, err := MakeJWT(userID, secret, expiresIn)
	if err != nil {
		t.Fatalf("Error generating JWT: %v", err)
	}

	// Validate with an incorrect secret; validation should fail.
	_, err = ValidateJWT(token, wrongSecret)
	if err == nil {
		t.Fatal("Expected error when validating JWT with wrong secret, got nil")
	}
}
