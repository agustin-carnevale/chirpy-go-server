package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestCheckPasswordHash(t *testing.T) {
	// First, we need to create some hashed passwords for testing
	password1 := "correctPassword123!"
	password2 := "anotherPassword456!"
	hash1, _ := HashPassword(password1)
	hash2, _ := HashPassword(password2)

	tests := []struct {
		name     string
		password string
		hash     string
		wantErr  bool
	}{
		{
			name:     "Correct password",
			password: password1,
			hash:     hash1,
			wantErr:  false,
		},
		{
			name:     "Incorrect password",
			password: "wrongPassword",
			hash:     hash1,
			wantErr:  true,
		},
		{
			name:     "Password doesn't match different hash",
			password: password1,
			hash:     hash2,
			wantErr:  true,
		},
		{
			name:     "Empty password",
			password: "",
			hash:     hash1,
			wantErr:  true,
		},
		{
			name:     "Invalid hash",
			password: password1,
			hash:     "invalidhash",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CheckPasswordHash(tt.password, tt.hash)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckPasswordHash() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestValidateJWT ensures that token validation works correctly
func TestValidateJWT(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "supersecretkey"
	expiresIn := time.Minute * 15

	// Generate a valid token
	token, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Fatalf("Failed to create JWT: %v", err)
	}

	// Generate a expired token
	expiredToken, err := MakeJWT(userID, tokenSecret, -time.Second)
	if err != nil {
		t.Fatalf("Failed to create JWT: %v", err)
	}

	tests := []struct {
		name        string
		tokenString string
		tokenSecret string
		expectErr   bool
	}{
		{"Valid token", token, tokenSecret, false},
		{"Invalid token secret", token, "wrongsecret", true},
		{"Malformed token", "invalid.token.string", tokenSecret, true},
		{"Token already expired", expiredToken, tokenSecret, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parsedUserID, err := ValidateJWT(tt.tokenString, tt.tokenSecret)
			if (err != nil) != tt.expectErr {
				t.Errorf("ValidateJWT() error = %v, expectErr %v", err, tt.expectErr)
			}
			if !tt.expectErr && parsedUserID != userID {
				t.Errorf("Expected user ID %v, got %v", userID, parsedUserID)
			}
		})
	}
}
