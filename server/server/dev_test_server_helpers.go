package server

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"testing"
	"time"

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

	type response struct {
		ID           pgtype.UUID `json:"id"`
		AccessToken  string      `json:"access_token"`
		RefreshToken string      `json:"refresh_token"`
	}
	var data response

	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		t.Fatalf("Failed to decode login response body: %v", err)
	}
	ctx.UserAcessToken = data.AccessToken
	ctx.UserRefreshToken = data.RefreshToken

	ctx.UserID = data.ID
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
	type parameters struct {
		UserID pgtype.UUID `json:"user_id"`
	}

	body := parameters{
		UserID: ctx.UserID,
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

// Create a medium (with custom fields) for testing use, return medium ID if needed
func (ctx *TestContext) CreateTestMediumCustom(t *testing.T, testBook parametersCreateMedium) pgtype.UUID {
	// Create medium via API request
	reqBody, err := json.Marshal(testBook)
	if err != nil {
		t.Fatalf("Failed to marshal body request for test medium: %v", err)
	}
	req, err := http.NewRequest("POST", ctx.BaseURL+"/api/media", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatalf("Failed to create test medium request: %v", err)
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

// Create a pre-set medium for testing use, return medium ID if needed
func (ctx *TestContext) CreateTestMediumRandom(t *testing.T) pgtype.UUID {
	randTitle := rand.Text()

	testBook := parametersCreateMedium{
		Title:       randTitle,
		MediaType:   "book",
		Creator:     "Test",
		ReleaseYear: 2025,
		ImageUrl:    "",
	}

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
	// To check if a user still exists in DB, make a request to GET /api/medium
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

// Create a record for testing use (start date = NOW, end date = NOW + 21 days)
func (ctx *TestContext) CreateTestRecord(t *testing.T, mediumID pgtype.UUID) pgtype.UUID {
	startDate := pgtype.Timestamp{
		Time:  time.Now().UTC(),
		Valid: true,
	}
	endDate := pgtype.Timestamp{
		Time:  time.Now().UTC().AddDate(0, 0, 21),
		Valid: true,
	}

	// Create Record via API request
	request := parametersCreateUserMediumRecord{
		MediaID:   mediumID,
		StartDate: startDate,
		EndDate:   endDate,
	}

	reqBody, err := json.Marshal(request)
	if err != nil {
		t.Fatalf("Failed to marshal body request for test create record: %v", err)
	}
	req, err := http.NewRequest("POST", ctx.BaseURL+"/api/records", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatalf("Failed to create test record request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", ctx.UserAcessToken))
	resp, err := ctx.Client.Do(req)
	if err != nil {
		t.Fatalf("Failed to create test record: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 201 {
		t.Fatalf("Failed to create test record. Status: %d", resp.StatusCode)
	}

	// Parse response
	type TestRecord struct {
		ID         pgtype.UUID `json:"id"`
		CreatedAt  string      `json:"created_at"`
		UpdatedAt  string      `json:"updated_at"`
		UserID     string      `json:"user_id"`
		MediaID    string      `json:"media_id"`
		IsFinished bool        `json:"is_finished"`
		StartDate  string      `json:"start_date"`
		EndDate    string      `json:"end_date"`
		Duration   int32       `json:"duration"`
	}

	var responseBody TestRecord
	err = json.NewDecoder(resp.Body).Decode(&responseBody)
	if err != nil {
		t.Fatalf("Failed to decode response body for test create records: %v", err)
	}

	return responseBody.ID

}

// Check if a record exist by reaching admin/record endpoint
func (ctx *TestContext) TestIfRecordExist(recordID pgtype.UUID) bool {
	// To check if a user still exists in DB, make a request to GET /api/medium
	body := parametersCheckRecordExists{
		RecordID: recordID,
	}
	requestBody, err := json.Marshal(body)
	if err != nil {
		log.Printf("ERROR with TestIfRecordExists() Marshal request body: %v\n", err)
	}
	req, err := http.NewRequest("GET", ctx.BaseURL+"/admin/record", bytes.NewBuffer(requestBody))
	if err != nil {
		log.Printf("ERROR with TestIfRecordExists() create request: %v\n", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := ctx.Client.Do(req)
	if err != nil {
		log.Printf("ERROR with TestIfRecordExists() do request: %v\n", err)
	}
	defer resp.Body.Close()
	return resp.StatusCode == 200
}
