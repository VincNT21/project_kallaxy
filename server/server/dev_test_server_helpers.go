package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"testing"

	"github.com/VincNT21/kallaxy/server/internal/auth"
	"github.com/jackc/pgx/v5/pgtype"
)

// Send a request to the admin Reset endpoint
func (ctx *TestContext) ResetDatabase(t *testing.T) {
	req, err := http.NewRequest("POST", ctx.BaseURL+"/admin/reset", nil)
	if err != nil {
		t.Fatalf("Failed to create reset request: %v", err)
	}
	resp, err := ctx.Client.Do(req)
	if err != nil {
		t.Fatalf("Failed to reset database: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Fatalf("Failed to reset database. Status: %d", resp.StatusCode)
	}
}

// Creates a user
func (ctx *TestContext) CreateTestUser(t *testing.T) {
	// Create user via API request
	payload := fmt.Sprintf(`{"username":"%s", "password":"%s", "email":"%s"}`, ctx.UserUsername, ctx.UserPassword, ctx.UserEmail)
	resp, err := ctx.Client.Post(ctx.BaseURL+"/api/users", "application/json", strings.NewReader(payload))
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 201 {
		t.Fatalf("Failed to create test user. Status: %d", resp.StatusCode)
	}
}

// Log in a user and stores tokens in context variables
func (ctx *TestContext) LoginTestUser(t *testing.T) {
	// Login via API request
	payload := fmt.Sprintf(`{"username":"%s", "password":"%s"}`, ctx.UserUsername, ctx.UserPassword)
	resp, err := ctx.Client.Post(ctx.BaseURL+"/auth/login", "application/json", strings.NewReader(payload))
	if err != nil {
		t.Fatalf("Failed to login test user: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 201 {
		t.Fatalf("Failed to login test user. Status: %d", resp.StatusCode)
	}
	var response map[string]string
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Failed to decode login response body: %v", err)
	}
	ctx.UserAcessToken = response["access_token"]
	ctx.UserRefreshToken = response["refresh_token"]
	ctx.UserID = response["id"]
}

// Call auth.ValidateJWT
func TestValidateAccessToken(token string) bool {
	_, err := auth.ValidateJWT(token, "test-jwt-secret")
	return err == nil
}

// Chech if a Refresh token is valid
func (ctx *TestContext) TestValidateRefreshToken(token string) bool {
	// To check if validate, make a request to Refresh endpoint
	req, err := http.NewRequest("POST", ctx.BaseURL+"/auth/refresh", nil)
	if err != nil {
		log.Printf("ERROR with TestValidateRefreshToken() create request: %v\n", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	resp, err := ctx.Client.Do(req)
	if err != nil {
		log.Printf("ERROR with TestValidateRefreshToken() do request: %v\n", err)
	}
	defer resp.Body.Close()
	return resp.StatusCode == 201
}

// Check if logged user exist by reaching admin/user endpoint
func (ctx *TestContext) TestIfUserExist() bool {
	// To check if a user still exists in DB, make a request to GET /api/users
	body := map[string]string{
		"user_id": ctx.UserID,
	}
	requestBody, err := json.Marshal(body)
	if err != nil {
		log.Printf("ERROR with TestIfUserExist() Marshal request body: %v\n", err)
	}
	req, err := http.NewRequest("GET", ctx.BaseURL+"/admin/user", bytes.NewBuffer(requestBody))
	if err != nil {
		log.Printf("ERROR with TestIfUserExist() create request: %v\n", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := ctx.Client.Do(req)
	if err != nil {
		log.Printf("ERROR with TestIfUserExists() do request: %v\n", err)
	}
	defer resp.Body.Close()
	return resp.StatusCode == 200
}

// Create a medium for testing use, return medium ID if needed
func (ctx *TestContext) CreateTestMedium(t *testing.T, testBook parametersCreateMedium) pgtype.UUID {
	// Create medium via API request
	reqBody, err := json.Marshal(testBook)
	if err != nil {
		t.Fatalf("Failed to marshal body request for test medium: %v", err)
	}
	req, err := http.NewRequest("POST", ctx.BaseURL+"/api/media", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatalf("Failed to create test medium: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", ctx.UserAcessToken))
	resp, err := ctx.Client.Do(req)
	if err != nil {
		t.Fatalf("Failed to create test medium: %v", err)
	}

	defer resp.Body.Close()
	if resp.StatusCode != 201 {
		t.Fatalf("Failed to create test medium. Status: %d", resp.StatusCode)
	}

	// Parse response
	var responseBody Medium
	err = json.NewDecoder(resp.Body).Decode(&responseBody)
	if err != nil {
		t.Fatalf("Failed to decode response body for test medium: %v", err)
	}

	return responseBody.ID
}

// Check if a medium exist by reaching admin/medium endpoint
func (ctx *TestContext) TestIfMediumExist(mediumID pgtype.UUID) bool {
	// To check if a user still exists in DB, make a request to GET /api/users
	body := parametersCheckMediumExists{
		MediumID: mediumID,
	}
	requestBody, err := json.Marshal(body)
	if err != nil {
		log.Printf("ERROR with TestIfMediumExists() Marshal request body: %v\n", err)
	}
	req, err := http.NewRequest("GET", ctx.BaseURL+"/admin/medium", bytes.NewBuffer(requestBody))
	if err != nil {
		log.Printf("ERROR with TestIfMediumExists() create request: %v\n", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := ctx.Client.Do(req)
	if err != nil {
		log.Printf("ERROR with TestIfMediumExists() do request: %v\n", err)
	}
	defer resp.Body.Close()
	return resp.StatusCode == 200
}
