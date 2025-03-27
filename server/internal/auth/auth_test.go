package auth

import (
	"net/http"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func TestCheckPasswordHash(t *testing.T) {
	// Create some hashed passwords for testing
	password1 := "correctPassword123!"
	password2 := "anotherPassword456*"
	hash1, _ := HashPassword(password1)
	hash2, _ := HashPassword(password2)

	// Create tests table
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

	// Test loop

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CheckPasswordHash(tt.password, tt.hash)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckPasswordHash() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateJWT(t *testing.T) {
	// Create a user Id and a valid Token
	userIDString := "81c1cb0d-bbdb-4faa-aede-bd371a4ab722"
	var userID pgtype.UUID
	userID.Scan(userIDString)
	validToken, err := MakeJWT(userID, "secret", time.Hour)

	// If error with MakeJWT
	if err != nil {
		t.Fatalf("Error with MakeJWT: %v", err)
	}

	// Tests table
	tests := []struct {
		name        string
		tokenString string
		tokenSecret string
		wantUserID  pgtype.UUID
		wantErr     bool
	}{
		{
			name:        "Valid token",
			tokenString: validToken,
			tokenSecret: "secret",
			wantUserID:  userID,
			wantErr:     false,
		},
		{
			name:        "Invalid token",
			tokenString: "invalidToken",
			tokenSecret: "secret",
			wantUserID:  pgtype.UUID{},
			wantErr:     true,
		},
		{
			name:        "Wrong secret",
			tokenString: validToken,
			tokenSecret: "wrongSecret",
			wantUserID:  pgtype.UUID{},
			wantErr:     true,
		},
	}

	// Test loop
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotUserID, err := ValidateJWT(tt.tokenString, tt.tokenSecret)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateJWT() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotUserID != tt.wantUserID {
				t.Errorf("ValidateJWT() gotUserID = %v, want %v", gotUserID, tt.wantUserID)
			}
		})
	}

}

func TestGetBearer(t *testing.T) {
	// Create a HTTP header with Bearer authorization and one malformed
	req, _ := http.NewRequest("GET", "https://api.example.com", nil)
	req.Header.Set("Authorization", "Bearer token123")
	req2, _ := http.NewRequest("POST", "https://api.example.com", nil)
	req2.Header.Set("Authorization", "InvalidBearer token")

	// Tests table
	tests := []struct {
		name          string
		requestHeader http.Header
		wantToken     string
		wantErr       bool
	}{
		{
			name:          "Valid Bearer token",
			requestHeader: req.Header,
			wantToken:     "token123",
			wantErr:       false,
		},
		{
			name:          "Missing authorization header",
			requestHeader: http.Header{},
			wantToken:     "",
			wantErr:       true,
		},
		{
			name:          "Malformed authorization header",
			requestHeader: req2.Header,
			wantToken:     "",
			wantErr:       true,
		},
	}

	// Tests loop
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotToken, err := GetBearerToken(tt.requestHeader)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBearerToken() error = %v, want %v", err, tt.wantErr)
			}

			if gotToken != tt.wantToken {
				t.Errorf("GetBearerToken() gotToken = %v, want %v", gotToken, tt.wantToken)
			}
		})
	}
}
