package auth

import "testing"

func TestHashPassword(t *testing.T) {
	password := "testPasswordForRoshan123"
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword returned error: %v", err)
	}

	if hash == "" {
		t.Fatal("Expected non-empty hash, but got empty string")
	}

	if err := CheckPasswordHash(password, hash); err != nil {
		t.Fatalf("Expected password to match, got error: %v", err)
	}
}

func TestCheckPasswordHashValid(t *testing.T) {
	password := "myTestPassword"
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Error hashing password: %v", err)
	}

	err = CheckPasswordHash(password, hash)
	if err != nil {
		t.Errorf("Expected password to match, got error: %v", err)
	}
}

func TestCheckPasswordHashInvalid(t *testing.T) {
	password := "correctPassword"
	wrongPassword := "wrongPassword"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Error handling password: %v", err)
	}

	err = CheckPasswordHash(wrongPassword, hash)
	if err == nil {
		t.Error("Expected error for incorrect password, got nil")
	}

}
