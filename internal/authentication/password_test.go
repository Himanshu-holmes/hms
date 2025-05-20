package authentication

import (
	"fmt"
	"testing"
)

// TestHashPassword checks that hashing a password returns a valid bcrypt hash without error.
func TestHashPassword(t *testing.T) {
	password := "mySecurePassword123!"
	hash, err := HashPassword(password)
	fmt.Println(hash)
	if err != nil {
		t.Fatalf("HashPassword returned an error: %v", err)
	}

	if len(hash) == 0 {
		t.Error("Expected non-empty hash from HashPassword")
	}
}

// TestCheckPasswordHash verifies that the hash of a password matches the original password.
func TestCheckPasswordHash(t *testing.T) {
	password := "testPassword!"
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword returned an error: %v", err)
	}

	if !CheckPasswordHash(password, hash) {
		t.Error("CheckPasswordHash returned false for valid password and hash")
	}
}

// TestCheckPasswordHash_Invalid checks that an incorrect password does not match the hash.
func TestCheckPasswordHash_Invalid(t *testing.T) {
	password := "correctPassword"
	wrongPassword := "wrongPassword"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword returned an error: %v", err)
	}

	if CheckPasswordHash(wrongPassword, hash) {
		t.Error("CheckPasswordHash returned true for incorrect password")
	}
}
