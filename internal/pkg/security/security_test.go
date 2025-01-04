package security

import (
	"strings"
	"testing"
)

func TestHashPassword(t *testing.T) {
	password := "secure_password"

	hashedPassword, err := GenerateHash(password)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if hashedPassword == "" {
		t.Fatal("expected a hashed password, got an empty string")
	}

	// Ensure the hashed password has the correct format
	if len(hashedPassword) < 30 || !isValidHashFormat(hashedPassword) {
		t.Fatalf("unexpected hash format: %s", hashedPassword)
	}
}

func TestVerifyPassword(t *testing.T) {
	password := "secure_password"
	wrongPassword := "wrong_password"

	hashedPassword, err := GenerateHash(password)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	// Test with correct password
	isValid, err := VerifyPassword(password, hashedPassword)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if !isValid {
		t.Fatal("expected password to be valid, got invalid")
	}

	// Test with incorrect password
	isValid, err = VerifyPassword(wrongPassword, hashedPassword)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if isValid {
		t.Fatal("expected password to be invalid, got valid")
	}
}

func TestHashAndVerifyConsistency(t *testing.T) {
	password := "another_secure_password"

	// Hash the password multiple times
	hashedPassword1, err := GenerateHash(password)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	hashedPassword2, err := GenerateHash(password)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	// Hashes should not match due to different salts
	if hashedPassword1 == hashedPassword2 {
		t.Fatal("expected different hashes for the same password")
	}

	// Verify both hashes with the original password
	isValid1, err := VerifyPassword(password, hashedPassword1)
	if err != nil || !isValid1 {
		t.Fatalf("expected first hash to verify correctly, got: %v, valid: %v", err, isValid1)
	}

	isValid2, err := VerifyPassword(password, hashedPassword2)
	if err != nil || !isValid2 {
		t.Fatalf("expected second hash to verify correctly, got: %v, valid: %v", err, isValid2)
	}
}

// Helper function to check if the hash format is valid
func isValidHashFormat(hash string) bool {
	// A valid hash format should contain 6 parts separated by '$'
	parts := strings.Split(hash, "$")
	return len(parts) == 6 && parts[1] == "argon2id"
}
