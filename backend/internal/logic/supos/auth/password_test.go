package auth

import "testing"

func TestHashAndVerifyPassword(t *testing.T) {
	hash, err := hashPassword("Tier0@123")
	if err != nil {
		t.Fatalf("hashPassword() error = %v", err)
	}
	if hash == "" {
		t.Fatal("hashPassword() returned empty hash")
	}
	if !verifyPassword(hash, "Tier0@123") {
		t.Fatal("verifyPassword() should accept the original password")
	}
	if verifyPassword(hash, "wrong-password") {
		t.Fatal("verifyPassword() should reject a different password")
	}
}
