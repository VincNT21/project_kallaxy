package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"
)

/* ====================
Password managing
====================*/

// Hash a given password
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// Compare Hash and password
func CheckPasswordHash(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

/* ====================
JWT/Refresh token managing
====================*/

type TokenType string

const TokenTypeAccess TokenType = "kallaxy"

// Create a JSON Web Token
func MakeJWT(userID pgtype.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	signingKey := []byte(tokenSecret)
	// Create a new token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    string(TokenTypeAccess),
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		Subject:   userID.String(),
	})

	// Sign the token with the secret and return it
	signedToken, err := token.SignedString(signingKey)
	if err != nil {
		return "", fmt.Errorf("error with signing token: %v", err)
	}
	return signedToken, nil
}

// Validate a JSON Web Token
func ValidateJWT(tokenString, tokenSecret string) (pgtype.UUID, error) {
	claimsStruct := jwt.RegisteredClaims{}

	// Validate the signature of the JWT and extract the claims into a *jwt.Token struct
	token, err := jwt.ParseWithClaims(
		tokenString,
		&claimsStruct,
		func(token *jwt.Token) (interface{}, error) { return []byte(tokenSecret), nil },
	)
	if err != nil {
		return pgtype.UUID{}, fmt.Errorf("error with ParseWithClaims(): %v", err)
	}

	// Get access to token issuer and compare it to local const TokenTypeAccess
	issuer, err := token.Claims.GetIssuer()
	if err != nil {
		return pgtype.UUID{}, fmt.Errorf("error with Claims.GetIssuer(): %v", err)
	}
	if issuer != string(TokenTypeAccess) {
		return pgtype.UUID{}, errors.New("invalid issuer")
	}

	// Get access to user's id from the claims
	stringId, err := token.Claims.GetSubject()
	if err != nil {
		return pgtype.UUID{}, fmt.Errorf("error with Claims.GetSubject(): %v", err)
	}

	// Parse id into pgtype.UUID type
	var userID pgtype.UUID
	err = userID.Scan(stringId)
	if err != nil {
		return pgtype.UUID{}, fmt.Errorf("error with userID.Scan(): %v", err)
	}

	return userID, nil
}

// Generate a refresh token ()= a random 256 bites token encoded in hexa)
func MakeRefreshToken() (string, error) {
	// Generate 32 bytes of random data
	randomData := make([]byte, 32)
	_, err := rand.Read(randomData)
	if err != nil {
		return "", fmt.Errorf("error with rand.read(): %v", err)
	}

	// Convert it to hexa string
	refreshToken := hex.EncodeToString(randomData)

	return refreshToken, nil
}

// Get JWT/Refresh token from headers in request
func GetBearerToken(headers http.Header) (string, error) {
	// Get "Authorization" header
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("no auth header in request")
	}

	// Split it to get Bearer token
	splitAuth := strings.Split(authHeader, " ")
	if len(splitAuth) != 2 || splitAuth[0] != "Bearer" {
		return "", errors.New("malformed Authorization header")
	}

	return splitAuth[1], nil
}
