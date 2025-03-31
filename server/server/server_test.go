package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

// Tests use dev_test_server, which is a updated copy of production server

// All comments are on the first test only

type testCase struct {
	name           string
	requestHeaders map[string]string
	requestBody    map[string]string
	expectedStatus int
	expectResponse bool
	checkResponse  func(*testing.T, map[string]string)
	checkAfter     func(*testing.T)
}

/*=========================
TESTS FOR USERS ENDPOINTS
=========================*/

func TestCreateUser(t *testing.T) {
	ctx := SetupTestContext(t)
	defer ctx.Server.Shutdown(context.Background())

	ctx.CreateTestUser(t)
	ctx.LoginTestUser(t)

	testMethod := "POST"
	testEndpoint := ctx.BaseURL + "/api/users"

	testUser := parametersCreateUser{
		Username: "TestUser1",
		Password: "12345678",
		Email: "Test123@example.com",
	}

	tests := []struct {
		name           string
		requestHeaders map[string]string
		requestBody    parametersCreateUser
		expectedStatus int
		expectResponse bool
		checkResponse  func(*testing.T, User)
	}{
		{
			name: "Valid user creation",
			requestBody: testUser,
			expectedStatus: 201,
			expectResponse: true,
			checkResponse: func(t *testing.T, u User) {
				if u.ID
				if _, exists := body["id"]; !exists {
					t.Error("Response missing 'id' field")
				}
				if _, exists := body["created_at"]; !exists {
					t.Error("Response missing 'created_at' field")
				}
				if _, exists := body["updated_at"]; !exists {
					t.Error("Response missing 'updated_at' field")
				}
				if username, exists := body["username"]; !exists || username != "Testuser1" {
					t.Error("Response have incorrect 'username' field")
				}
				if email, exists := body["email"]; !exists || email != "Test123@example.com" {
					t.Error("Response have incorrect 'email' field")
				}
			},
		},
		{
			name: "Duplicate user's username",
			requestBody: map[string]string{
				"username": "Testuser1",
				"password": "hjldsfoeri",
				"email":    "Test123@example2.com",
			},
			expectedStatus: 409,
		},
		{
			name: "Duplicate user's email",
			requestBody: map[string]string{
				"username": "Testuser2",
				"password": "hjldsfoeri",
				"email":    "Test123@example.com",
			},
			expectedStatus: 409,
		},
		{
			name: "Missing a field",
			requestBody: map[string]string{
				"username": "Testuser2",
				"password": "",
				"email":    "Test123@example.com",
			},
			expectedStatus: 400,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			requestBody, _ := json.Marshal(tc.requestBody)
			req, _ := http.NewRequest(testMethod, testEndpoint, bytes.NewBuffer(requestBody))
			req.Header.Set("Content-Type", "application/json")
			if tc.requestHeaders != nil {
				for headerKey, headerValue := range tc.requestHeaders {
					req.Header.Set(headerKey, headerValue)
				}
			}
			resp, err := ctx.Client.Do(req)
			if err != nil {
				t.Fatalf("Failed to send request: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tc.expectedStatus {
				t.Errorf("Expected status code %d, got %d", tc.expectedStatus, resp.StatusCode)
			}

			if tc.expectResponse {
				var responseBody User
				err := json.NewDecoder(resp.Body).Decode(&responseBody)
				if err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}
				if tc.checkResponse != nil {
					tc.checkResponse(t, responseBody)
				}
			}
		})
	}
}

func TestCreateUserOld(t *testing.T) {
	// Setup test environnement
	ctx := SetupTestContext(t)
	defer ctx.Server.Shutdown(context.Background()) // Clean shutdown when test ends

	// Define test method and endpoint
	testMethod := "POST"
	testEndpoint := ctx.BaseURL + "/api/users"

	// Create tests table
	tests := []testCase{
		{
			name: "Valid user creation",
			requestBody: map[string]string{
				"username": "Testuser1",
				"password": "12345678",
				"email":    "Test123@example.com",
			},
			expectedStatus: 201,
			expectResponse: true,
			checkResponse: func(t *testing.T, body map[string]string) {
				if _, exists := body["id"]; !exists {
					t.Error("Response missing 'id' field")
				}
				if _, exists := body["created_at"]; !exists {
					t.Error("Response missing 'created_at' field")
				}
				if _, exists := body["updated_at"]; !exists {
					t.Error("Response missing 'updated_at' field")
				}
				if username, exists := body["username"]; !exists || username != "Testuser1" {
					t.Error("Response have incorrect 'username' field")
				}
				if email, exists := body["email"]; !exists || email != "Test123@example.com" {
					t.Error("Response have incorrect 'email' field")
				}
			},
		},
		{
			name: "Duplicate user's username",
			requestBody: map[string]string{
				"username": "Testuser1",
				"password": "hjldsfoeri",
				"email":    "Test123@example2.com",
			},
			expectedStatus: 409,
		},
		{
			name: "Duplicate user's email",
			requestBody: map[string]string{
				"username": "Testuser2",
				"password": "hjldsfoeri",
				"email":    "Test123@example.com",
			},
			expectedStatus: 409,
		},
		{
			name: "Missing a field",
			requestBody: map[string]string{
				"username": "Testuser2",
				"password": "",
				"email":    "Test123@example.com",
			},
			expectedStatus: 400,
		},
	}

	// Test loop
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create request body
			requestBody, _ := json.Marshal(tc.requestBody)

			// Create HTTP request
			req, _ := http.NewRequest(testMethod, testEndpoint, bytes.NewBuffer(requestBody))
			req.Header.Set("Content-Type", "application/json")

			// Send the request
			resp, err := ctx.Client.Do(req)
			if err != nil {
				t.Fatalf("Failed to send request: %v", err)
			}
			defer resp.Body.Close()

			// Check status code
			if resp.StatusCode != tc.expectedStatus {
				t.Errorf("Expected status code %d, got %d", tc.expectedStatus, resp.StatusCode)
			}

			// If JSON response expected, check the response body
			if tc.expectResponse {
				var responseBody map[string]string
				err := json.NewDecoder(resp.Body).Decode(&responseBody)
				if err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}
				if tc.checkResponse != nil {
					tc.checkResponse(t, responseBody)
				}
			}
			if tc.checkAfter != nil {
				tc.checkAfter(t)
			}
		})
	}
}

func TestGetUser(t *testing.T) {
	ctx := SetupTestContext(t)
	defer ctx.Server.Shutdown(context.Background())

	ctx.CreateTestUser(t)
	ctx.LoginTestUser(t)

	testMethod := "GET"
	testEndpoint := ctx.BaseURL + "/api/users"

	tests := []testCase{
		{
			name: "Valid",
			requestHeaders: map[string]string{
				"Authorization": fmt.Sprintf("Bearer %s", ctx.UserAcessToken),
			},
			expectedStatus: 200,
			expectResponse: true,
			checkResponse: func(t *testing.T, m map[string]string) {
				if _, exists := m["id"]; !exists {
					t.Error("Response missing 'id' field")
				}
				if username, exists := m["username"]; !exists || username != ctx.UserUsername {
					t.Error("Response field 'username' incorrect")
				}
				if email, exists := m["email"]; !exists || email != ctx.UserEmail {
					t.Error("Response have incorrect 'email' field")
				}
			},
		},
		{
			name:           "No access_token",
			expectedStatus: 401,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			requestBody, _ := json.Marshal(tc.requestBody)
			req, _ := http.NewRequest(testMethod, testEndpoint, bytes.NewBuffer(requestBody))
			req.Header.Set("Content-Type", "application/json")
			if tc.requestHeaders != nil {
				for headerKey, headerValue := range tc.requestHeaders {
					req.Header.Set(headerKey, headerValue)
				}
			}
			resp, err := ctx.Client.Do(req)
			if err != nil {
				t.Fatalf("Failed to send request: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tc.expectedStatus {
				t.Errorf("Expected status code %d, got %d", tc.expectedStatus, resp.StatusCode)
			}

			if tc.expectResponse {
				var responseBody map[string]string
				err := json.NewDecoder(resp.Body).Decode(&responseBody)
				if err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}
				if tc.checkResponse != nil {
					tc.checkResponse(t, responseBody)
				}
			}
			if tc.checkAfter != nil {
				tc.checkAfter(t)
			}
		})
	}
}

func TestUpdateUser(t *testing.T) {
	ctx := SetupTestContext(t)
	defer ctx.Server.Shutdown(context.Background())

	ctx.CreateTestUser(t)
	ctx.LoginTestUser(t)

	testMethod := "PUT"
	testEndpoint := ctx.BaseURL + "/api/users"

	tests := []testCase{
		{
			name: "Valid token and valid data",
			requestHeaders: map[string]string{
				"Authorization": fmt.Sprintf("Bearer %s", ctx.UserAcessToken),
			},
			requestBody: map[string]string{
				"username": ctx.UserUsername,
				"password": ctx.UserPassword,
				"email":    "newemail@example.com",
			},
			expectedStatus: 200,
			expectResponse: true,
			checkResponse: func(t *testing.T, m map[string]string) {
				if _, exists := m["id"]; !exists {
					t.Error("Response missing 'id' field")
				}
				if _, exists := m["updated_at"]; !exists {
					t.Error("Response missing 'updated_at' field")
				}
				if m["updated_at"] == m["created_at"] {
					t.Error("Response field 'updated_at' was not updated")
				}
				if username, exists := m["username"]; !exists || username != ctx.UserUsername {
					t.Error("Response have incorrect 'username' field")
				}
				if email, exists := m["email"]; !exists || email != "newemail@example.com" {
					t.Error("Response have incorrect 'email' field")
				}
			},
			checkAfter: func(t *testing.T) {
				if ctx.TestValidateRefreshToken(ctx.UserRefreshToken) {
					t.Error("Refresh token is still valid")
				}
			},
		},
		{
			name: "Valid token but missing a field",
			requestHeaders: map[string]string{
				"Authorization": fmt.Sprintf("Bearer %s", ctx.UserAcessToken),
			},
			requestBody: map[string]string{
				"username": ctx.UserUsername,
				"password": ctx.UserPassword,
				"email":    "",
			},
			expectedStatus: 400,
		},
		{
			name: "Missing token",
			requestBody: map[string]string{
				"username": ctx.UserUsername,
				"password": ctx.UserPassword,
				"email":    "newemail@example.com",
			},
			expectedStatus: 401,
		},
		{
			name: "Invalid token",
			requestHeaders: map[string]string{
				"Authorization": "Bearer badaccesstoken",
			},
			requestBody: map[string]string{
				"username": ctx.UserUsername,
				"password": ctx.UserPassword,
				"email":    "newemail@example.com",
			},
			expectedStatus: 401,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			requestBody, _ := json.Marshal(tc.requestBody)
			req, _ := http.NewRequest(testMethod, testEndpoint, bytes.NewBuffer(requestBody))
			req.Header.Set("Content-Type", "application/json")
			if tc.requestHeaders != nil {
				for headerKey, headerValue := range tc.requestHeaders {
					req.Header.Set(headerKey, headerValue)
				}
			}
			resp, err := ctx.Client.Do(req)
			if err != nil {
				t.Fatalf("Failed to send request: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tc.expectedStatus {
				t.Errorf("Expected status code %d, got %d", tc.expectedStatus, resp.StatusCode)
			}

			if tc.expectResponse {
				var responseBody map[string]string
				err := json.NewDecoder(resp.Body).Decode(&responseBody)
				if err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}
				if tc.checkResponse != nil {
					tc.checkResponse(t, responseBody)
				}
			}
			if tc.checkAfter != nil {
				tc.checkAfter(t)
			}
		})
	}
}

func TestDeleteUser(t *testing.T) {
	ctx := SetupTestContext(t)
	defer ctx.Server.Shutdown(context.Background())

	ctx.CreateTestUser(t)
	ctx.LoginTestUser(t)

	testMethod := "DELETE"
	testEndpoint := ctx.BaseURL + "/api/users"

	tests := []testCase{
		{
			name: "Valid",
			requestHeaders: map[string]string{
				"Authorization": fmt.Sprintf("Bearer %s", ctx.UserAcessToken),
			},
			expectedStatus: 200,
			checkAfter: func(t *testing.T) {
				if ctx.TestIfUserExist() {
					t.Error("User still exists in database")
				}
			},
		},
		{
			name:           "No access_token",
			expectedStatus: 401,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			requestBody, _ := json.Marshal(tc.requestBody)
			req, _ := http.NewRequest(testMethod, testEndpoint, bytes.NewBuffer(requestBody))
			req.Header.Set("Content-Type", "application/json")
			if tc.requestHeaders != nil {
				for headerKey, headerValue := range tc.requestHeaders {
					req.Header.Set(headerKey, headerValue)
				}
			}
			resp, err := ctx.Client.Do(req)
			if err != nil {
				t.Fatalf("Failed to send request: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tc.expectedStatus {
				t.Errorf("Expected status code %d, got %d", tc.expectedStatus, resp.StatusCode)
			}

			if tc.expectResponse {
				var responseBody map[string]string
				err := json.NewDecoder(resp.Body).Decode(&responseBody)
				if err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}
				if tc.checkResponse != nil {
					tc.checkResponse(t, responseBody)
				}
			}
			if tc.checkAfter != nil {
				tc.checkAfter(t)
			}
		})
	}
}

/*===================================
TESTS FOR AUTHENTIFICATION ENDPOINTS
===================================*/

func TestLogin(t *testing.T) {
	ctx := SetupTestContext(t)
	defer ctx.Server.Shutdown(context.Background())

	ctx.CreateTestUser(t)

	testMethod := "POST"
	testEndpoint := ctx.BaseURL + "/auth/login"

	tests := []testCase{
		{
			name: "Valid login",
			requestBody: map[string]string{
				"username": ctx.UserUsername,
				"password": ctx.UserPassword,
			},
			expectedStatus: 201,
			expectResponse: true,
			checkResponse: func(t *testing.T, m map[string]string) {
				if _, exists := m["id"]; !exists {
					t.Error("Response missing 'id' field")
				}
				if username, exists := m["username"]; !exists || username != ctx.UserUsername {
					t.Error("Response have incorrect 'username' field")
				}
				access_token, exists := m["access_token"]
				if !exists {
					t.Error("Response missing 'access_token' field")
				}
				if !TestValidateAccessToken(access_token) {
					t.Error("Response's access_token is invalid")
				}
				refreshToken, exists := m["refresh_token"]
				if !exists {
					t.Error("Response missing 'refresh_token' field")
				}
				if !ctx.TestValidateRefreshToken(refreshToken) {
					t.Error("Response's refresh_token is invalid")
				}
			},
		},
		{
			name: "Invalid password",
			requestBody: map[string]string{
				"username": ctx.UserUsername,
				"password": "12345",
			},
			expectedStatus: 401,
		},
		{
			name: "Invalid username",
			requestBody: map[string]string{
				"username": "Testtest12",
				"password": ctx.UserPassword,
			},
			expectedStatus: 401,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			requestBody, _ := json.Marshal(tc.requestBody)
			req, _ := http.NewRequest(testMethod, testEndpoint, bytes.NewBuffer(requestBody))
			req.Header.Set("Content-Type", "application/json")
			resp, err := ctx.Client.Do(req)
			if err != nil {
				t.Fatalf("Failed to send request: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tc.expectedStatus {
				t.Errorf("Expected status code %d, got %d", tc.expectedStatus, resp.StatusCode)
			}

			if tc.expectResponse {
				var responseBody map[string]string
				err := json.NewDecoder(resp.Body).Decode(&responseBody)
				if err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}
				if tc.checkResponse != nil {
					tc.checkResponse(t, responseBody)
				}
			}
			if tc.checkAfter != nil {
				tc.checkAfter(t)
			}
		})
	}

}

func TestLogout(t *testing.T) {
	ctx := SetupTestContext(t)
	defer ctx.Server.Shutdown(context.Background())

	ctx.CreateTestUser(t)
	ctx.LoginTestUser(t)

	testMethod := "POST"
	testEndpoint := ctx.BaseURL + "/auth/logout"

	tests := []testCase{
		{
			name: "Logout ok",
			requestHeaders: map[string]string{
				"Authorization": fmt.Sprintf("Bearer %s", ctx.UserAcessToken),
			},
			expectedStatus: 204,
			checkAfter: func(t *testing.T) {
				if ctx.TestValidateRefreshToken(ctx.UserRefreshToken) {
					t.Error("Refresh token is still valid")
				}
			},
		},
		{
			name: "Invalid token provided",
			requestHeaders: map[string]string{
				"Authorization": "Bearer wrongToken",
			},
			expectedStatus: 401,
		},
		{
			name:           "Missing Authorization header",
			expectedStatus: 401,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			requestBody, _ := json.Marshal(tc.requestBody)
			req, _ := http.NewRequest(testMethod, testEndpoint, bytes.NewBuffer(requestBody))
			req.Header.Set("Content-Type", "application/json")
			if tc.requestHeaders != nil {
				for headerKey, headerValue := range tc.requestHeaders {
					req.Header.Set(headerKey, headerValue)
				}
			}
			resp, err := ctx.Client.Do(req)
			if err != nil {
				t.Fatalf("Failed to send request: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tc.expectedStatus {
				t.Errorf("Expected status code %d, got %d", tc.expectedStatus, resp.StatusCode)
			}

			if resp.StatusCode == 200 || resp.StatusCode == 201 {
				var responseBody map[string]string
				err := json.NewDecoder(resp.Body).Decode(&responseBody)
				if err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}
				if tc.checkResponse != nil {
					tc.checkResponse(t, responseBody)
				}
			}
			if tc.checkAfter != nil {
				tc.checkAfter(t)
			}
		})
	}
}

func TestRefreshTokens(t *testing.T) {
	ctx := SetupTestContext(t)
	defer ctx.Server.Shutdown(context.Background())

	ctx.CreateTestUser(t)
	ctx.LoginTestUser(t)

	testMethod := "POST"
	testEndpoint := ctx.BaseURL + "/auth/refresh"

	tests := []testCase{
		{
			name: "Valid Refresh Token",
			requestHeaders: map[string]string{
				"Authorization": fmt.Sprintf("Bearer %s", ctx.UserRefreshToken),
			},
			expectedStatus: 201,
			expectResponse: true,
			checkResponse: func(t *testing.T, m map[string]string) {
				access_token, exists := m["access_token"]
				if !exists {
					t.Error("Response missing 'access_token' field")
				}
				if !TestValidateAccessToken(access_token) {
					t.Error("Response's access_token is invalid")
				}
				refreshToken, exists := m["refresh_token"]
				if !exists {
					t.Error("Response missing 'refresh_token' field")
				}
				if !ctx.TestValidateRefreshToken(refreshToken) {
					t.Error("Response's refresh_token is invalid")
				}
			},
		},
		{
			name: "Invalid Refresh Token",
			requestHeaders: map[string]string{
				"Authorization": "Bearer invalidrefreshtoken",
			},
			expectedStatus: 401,
		},
		{
			name:           "Refresh Token not provided",
			expectedStatus: 401,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			requestBody, _ := json.Marshal(tc.requestBody)
			req, _ := http.NewRequest(testMethod, testEndpoint, bytes.NewBuffer(requestBody))
			req.Header.Set("Content-Type", "application/json")
			if tc.requestHeaders != nil {
				for headerKey, headerValue := range tc.requestHeaders {
					req.Header.Set(headerKey, headerValue)
				}
			}
			resp, err := ctx.Client.Do(req)
			if err != nil {
				t.Fatalf("Failed to send request: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tc.expectedStatus {
				t.Errorf("Expected status code %d, got %d", tc.expectedStatus, resp.StatusCode)
			}

			if tc.expectResponse {
				var responseBody map[string]string
				err := json.NewDecoder(resp.Body).Decode(&responseBody)
				if err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}
				if tc.checkResponse != nil {
					tc.checkResponse(t, responseBody)
				}
			}
		})
	}
}

func TestRevokeToken(t *testing.T) {
	ctx := SetupTestContext(t)
	defer ctx.Server.Shutdown(context.Background())

	ctx.CreateTestUser(t)
	ctx.LoginTestUser(t)

	testMethod := "POST"
	testEndpoint := ctx.BaseURL + "/auth/revoke"

	tests := []testCase{
		{
			name: "Revoke ok",
			requestHeaders: map[string]string{
				"Authorization": fmt.Sprintf("Bearer %s", ctx.UserRefreshToken),
			},
			expectedStatus: 204,
			checkAfter: func(t *testing.T) {
				if ctx.TestValidateRefreshToken(ctx.UserRefreshToken) {
					t.Error("Refresh token is still valid")
				}
			},
		},
		{
			name:           "Missing Authorization header",
			expectedStatus: 401,
		},
		{
			name: "Invalid token",
			requestHeaders: map[string]string{
				"Authorization": "Bearer wrongrefreshtoken",
			},
			expectedStatus: 404,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			requestBody, _ := json.Marshal(tc.requestBody)
			req, _ := http.NewRequest(testMethod, testEndpoint, bytes.NewBuffer(requestBody))
			req.Header.Set("Content-Type", "application/json")
			if tc.requestHeaders != nil {
				for headerKey, headerValue := range tc.requestHeaders {
					req.Header.Set(headerKey, headerValue)
				}
			}
			resp, err := ctx.Client.Do(req)
			if err != nil {
				t.Fatalf("Failed to send request: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tc.expectedStatus {
				t.Errorf("Expected status code %d, got %d", tc.expectedStatus, resp.StatusCode)
			}

			if tc.expectResponse {
				var responseBody map[string]string
				err := json.NewDecoder(resp.Body).Decode(&responseBody)
				if err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}
				if tc.checkResponse != nil {
					tc.checkResponse(t, responseBody)
				}
			}
			if tc.checkAfter != nil {
				tc.checkAfter(t)
			}
		})
	}
}

/*
=========================
TESTS FOR MEDIA ENDPOINTS
=========================
*/

func TestCreateMedium(t *testing.T) {
	ctx := SetupTestContext(t)
	defer ctx.Server.Shutdown(context.Background())

	ctx.CreateTestUser(t)
	ctx.LoginTestUser(t)

	testMethod := "POST"
	testEndpoint := ctx.BaseURL + "/api/media"

	testBook := parametersCreateMedium{
		Title:       "The Fellowship of the Ring",
		MediaType:   "book",
		Creator:     "J.R.R Tolkien",
		ReleaseYear: 1954,
		ImageUrl:    "https://upload.wikimedia.org/wikipedia/en/thumb/8/8e/The_Fellowship_of_the_Ring_cover.gif/220px-The_Fellowship_of_the_Ring_cover.gif",
	}

	tests := []struct {
		name           string
		requestHeaders map[string]string
		requestBody    parametersCreateMedium
		expectedStatus int
		expectResponse bool
		checkResponse  func(*testing.T, Medium)
		checkAfter     func(*testing.T)
	}{
		{
			name: "Valid, with no metadata",
			requestHeaders: map[string]string{
				"Authorization": fmt.Sprintf("Bearer %s", ctx.UserAcessToken),
			},
			requestBody:    testBook,
			expectedStatus: 201,
			expectResponse: true,
			checkResponse: func(t *testing.T, r Medium) {
				nilID := pgtype.UUID{}
				if r.ID == nilID {
					t.Error("'id' response field missing")
				}
				nilTime := pgtype.Timestamp{}
				if r.CreatedAt == nilTime {
					t.Error("'created_at' response field missing")
				}
				if r.UpdatedAt == nilTime {
					t.Error("'updated_at' response field missing")
				}
				if r.Title != testBook.Title {
					t.Error("'title' response field incorrect")
				}
				if r.MediaType != testBook.MediaType {
					t.Error("'media_type' response field incorrect")
				}
				if r.Creator != testBook.Creator {
					t.Error("'creator' response field incorrect")
				}
				if r.ReleaseYear != testBook.ReleaseYear {
					t.Error("'release_year' response field incorrect")
				}
				imageUrl := pgtype.Text{String: testBook.ImageUrl, Valid: true}
				if r.ImageUrl != imageUrl {
					t.Error("'image_url' response field incorrect")
				}
			},
		},
		{
			name:           "No access_token",
			expectedStatus: 401,
		},
		{
			name: "Duplicate Title",
			requestHeaders: map[string]string{
				"Authorization": fmt.Sprintf("Bearer %s", ctx.UserAcessToken),
			},
			requestBody:    testBook,
			expectedStatus: 409,
		},
		{
			name: "Missing a needed field",
			requestHeaders: map[string]string{
				"Authorization": fmt.Sprintf("Bearer %s", ctx.UserAcessToken),
			},
			requestBody: parametersCreateMedium{
				Title:       "The Fellowship of the Ring2",
				MediaType:   "book",
				ReleaseYear: 1954,
				ImageUrl:    "https://upload.wikimedia.org/wikipedia/en/thumb/8/8e/The_Fellowship_of_the_Ring_cover.gif/220px-The_Fellowship_of_the_Ring_cover.gif",
			},
			expectedStatus: 400,
		},
		{
			name: "Empty image_url field",
			requestHeaders: map[string]string{
				"Authorization": fmt.Sprintf("Bearer %s", ctx.UserAcessToken),
			},
			requestBody: parametersCreateMedium{
				Title:       "The Fellowship of the Ring3",
				MediaType:   "book",
				Creator:     "J.R.R Tolkien",
				ReleaseYear: 1954,
				ImageUrl:    "",
			},
			expectedStatus: 201,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			requestBody, _ := json.Marshal(tc.requestBody)
			req, _ := http.NewRequest(testMethod, testEndpoint, bytes.NewBuffer(requestBody))
			req.Header.Set("Content-Type", "application/json")
			if tc.requestHeaders != nil {
				for headerKey, headerValue := range tc.requestHeaders {
					req.Header.Set(headerKey, headerValue)
				}
			}
			resp, err := ctx.Client.Do(req)
			if err != nil {
				t.Fatalf("Failed to send request: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tc.expectedStatus {
				t.Errorf("Expected status code %d, got %d", tc.expectedStatus, resp.StatusCode)
			}

			if tc.expectResponse {
				var responseBody Medium
				err := json.NewDecoder(resp.Body).Decode(&responseBody)
				if err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}
				if tc.checkResponse != nil {
					tc.checkResponse(t, responseBody)
				}
			}
		})
	}
}

func TestGetMediaByTitle(t *testing.T) {
	ctx := SetupTestContext(t)
	defer ctx.Server.Shutdown(context.Background())

	ctx.CreateTestUser(t)
	ctx.LoginTestUser(t)

	testBook := parametersCreateMedium{
		Title:       "The Fellowship of the Ring",
		MediaType:   "book",
		Creator:     "J.R.R Tolkien",
		ReleaseYear: 1954,
		ImageUrl:    "https://upload.wikimedia.org/wikipedia/en/thumb/8/8e/The_Fellowship_of_the_Ring_cover.gif/220px-The_Fellowship_of_the_Ring_cover.gif",
	}

	ctx.CreateTestMediumCustom(t, testBook)

	testMethod := "GET"
	testEndpoint := ctx.BaseURL + "/api/media"

	tests := []struct {
		name           string
		requestHeaders map[string]string
		queryParameter string
		requestBody    map[string]string
		expectedStatus int
		expectResponse bool
		checkResponse  func(*testing.T, Medium)
		checkAfter     func(*testing.T)
	}{
		{
			name: "Valid",
			requestHeaders: map[string]string{
				"Authorization": fmt.Sprintf("Bearer %s", ctx.UserAcessToken),
			},
			queryParameter: "?title=The+Fellowship+Of+The+Ring",
			expectedStatus: 200,
			expectResponse: true,
			checkResponse: func(t *testing.T, m Medium) {
				if m.Title != testBook.Title {
					t.Error("incorrect Title")
				}
				if m.MediaType != testBook.MediaType {
					t.Error("incorrect MediaType")
				}
				if m.Creator != testBook.Creator {
					t.Error("incorrect Creator")
				}
				if m.ReleaseYear != testBook.ReleaseYear {
					t.Error("incorrect ReleaseYear")
				}
				imageUrl := pgtype.Text{String: testBook.ImageUrl, Valid: true}
				if m.ImageUrl != imageUrl {
					t.Error("incorrect ImageUrl")
				}
				nilID := pgtype.UUID{}
				if m.ID == nilID {
					t.Error("'id' response field missing")
				}
				nilTime := pgtype.Timestamp{}
				if m.CreatedAt == nilTime {
					t.Error("'created_at' response field missing")
				}
				if m.UpdatedAt == nilTime {
					t.Error("'updated_at' response field missing")
				}
			},
		},
		{
			name:           "No access_token",
			expectedStatus: 401,
		},
		{
			name: "Invalid title",
			requestHeaders: map[string]string{
				"Authorization": fmt.Sprintf("Bearer %s", ctx.UserAcessToken),
			},
			queryParameter: "?title=The+Fellowshi+Of+The+Ring",
			expectedStatus: 404,
		},
		{
			name: "Invalid query parameter",
			requestHeaders: map[string]string{
				"Authorization": fmt.Sprintf("Bearer %s", ctx.UserAcessToken),
			},
			queryParameter: "?til=The+Fellowshi+Of+The+Ring",
			expectedStatus: 400,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			requestBody, _ := json.Marshal(tc.requestBody)
			req, _ := http.NewRequest(testMethod, testEndpoint+tc.queryParameter, bytes.NewBuffer(requestBody))
			req.Header.Set("Content-Type", "application/json")
			if tc.requestHeaders != nil {
				for headerKey, headerValue := range tc.requestHeaders {
					req.Header.Set(headerKey, headerValue)
				}
			}
			resp, err := ctx.Client.Do(req)
			if err != nil {
				t.Fatalf("Failed to send request: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tc.expectedStatus {
				t.Errorf("Expected status code %d, got %d", tc.expectedStatus, resp.StatusCode)
			}

			if tc.expectResponse {
				var responseBody Medium
				err := json.NewDecoder(resp.Body).Decode(&responseBody)
				if err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}
				if tc.checkResponse != nil {
					tc.checkResponse(t, responseBody)
				}
			}
			if tc.checkAfter != nil {
				tc.checkAfter(t)
			}
		})
	}
}

func TestGetMediaByType(t *testing.T) {
	ctx := SetupTestContext(t)
	defer ctx.Server.Shutdown(context.Background())

	ctx.CreateTestUser(t)
	ctx.LoginTestUser(t)

	testBook := parametersCreateMedium{
		Title:       "The Fellowship of the Ring",
		MediaType:   "book",
		Creator:     "J.R.R Tolkien",
		ReleaseYear: 1954,
		ImageUrl:    "https://upload.wikimedia.org/wikipedia/en/8/8e/The_Fellowship_of_the_Ring_cover.gif",
	}
	testBook2 := parametersCreateMedium{
		Title:       "The Two Towers",
		MediaType:   "book",
		Creator:     "J.R.R Tolkien",
		ReleaseYear: 1954,
		ImageUrl:    "https://upload.wikimedia.org/wikipedia/en/a/a1/The_Two_Towers_cover.gif",
	}

	ctx.CreateTestMediumCustom(t, testBook)
	ctx.CreateTestMediumCustom(t, testBook2)

	testMethod := "GET"
	testEndpoint := ctx.BaseURL + "/api/media"

	tests := []struct {
		name           string
		requestHeaders map[string]string
		queryParameter string
		requestBody    map[string]string
		expectedStatus int
		expectResponse bool
		checkResponse  func(*testing.T, responseGetMediaByType)
		checkAfter     func(*testing.T)
	}{
		{
			name: "Valid",
			requestHeaders: map[string]string{
				"Authorization": fmt.Sprintf("Bearer %s", ctx.UserAcessToken),
			},
			queryParameter: "?type=book",
			expectedStatus: 200,
			expectResponse: true,
			checkResponse: func(t *testing.T, r responseGetMediaByType) {
				for _, m := range r.Media {
					if m.Title != testBook.Title && m.Title != testBook2.Title {
						t.Error("incorrect Title")
					}
					if m.MediaType != testBook.MediaType {
						t.Error("incorrect MediaType")
					}
					if m.Creator != testBook.Creator && m.Creator != testBook2.Creator {
						t.Error("incorrect Creator")
					}
					if m.ReleaseYear != testBook.ReleaseYear && m.ReleaseYear != testBook2.ReleaseYear {
						t.Error("incorrect ReleaseYear")
					}
					imageUrl1 := pgtype.Text{String: testBook.ImageUrl, Valid: true}
					imageUrl2 := pgtype.Text{String: testBook2.ImageUrl, Valid: true}
					if m.ImageUrl != imageUrl1 && m.ImageUrl != imageUrl2 {
						t.Error("incorrect ImageUrl")
					}
					nilID := pgtype.UUID{}
					if m.ID == nilID {
						t.Error("'id' response field missing")
					}
					nilTime := pgtype.Timestamp{}
					if m.CreatedAt == nilTime {
						t.Error("'created_at' response field missing")
					}
					if m.UpdatedAt == nilTime {
						t.Error("'updated_at' response field missing")
					}
				}

			},
		},
		{
			name:           "No access_token",
			expectedStatus: 401,
		},
		{
			name: "Invalid media type",
			requestHeaders: map[string]string{
				"Authorization": fmt.Sprintf("Bearer %s", ctx.UserAcessToken),
			},
			queryParameter: "?type=comic",
			expectedStatus: 404,
		},
		{
			name: "Invalid query parameter",
			requestHeaders: map[string]string{
				"Authorization": fmt.Sprintf("Bearer %s", ctx.UserAcessToken),
			},
			queryParameter: "?typ=book",
			expectedStatus: 400,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			requestBody, _ := json.Marshal(tc.requestBody)
			req, _ := http.NewRequest(testMethod, testEndpoint+tc.queryParameter, bytes.NewBuffer(requestBody))
			req.Header.Set("Content-Type", "application/json")
			if tc.requestHeaders != nil {
				for headerKey, headerValue := range tc.requestHeaders {
					req.Header.Set(headerKey, headerValue)
				}
			}
			resp, err := ctx.Client.Do(req)
			if err != nil {
				t.Fatalf("Failed to send request: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tc.expectedStatus {
				t.Errorf("Expected status code %d, got %d", tc.expectedStatus, resp.StatusCode)
			}

			if tc.expectResponse {
				var responseBody responseGetMediaByType
				err := json.NewDecoder(resp.Body).Decode(&responseBody)
				if err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}
				if tc.checkResponse != nil {
					tc.checkResponse(t, responseBody)
				}
			}
			if tc.checkAfter != nil {
				tc.checkAfter(t)
			}
		})
	}
}

func TestUpdateMedium(t *testing.T) {
	ctx := SetupTestContext(t)
	defer ctx.Server.Shutdown(context.Background())

	ctx.CreateTestUser(t)
	ctx.LoginTestUser(t)

	testBook := parametersCreateMedium{
		Title:       "The Fellowship of the Ring",
		MediaType:   "book",
		Creator:     "J.R.R Tolkien",
		ReleaseYear: 1954,
		ImageUrl:    "https://upload.wikimedia.org/wikipedia/en/thumb/8/8e/The_Fellowship_of_the_Ring_cover.gif/220px-The_Fellowship_of_the_Ring_cover.gif",
	}
	testBook2 := parametersCreateMedium{
		Title:       "The Two Towers",
		MediaType:   "book",
		Creator:     "J.R.R Tolkien",
		ReleaseYear: 1954,
		ImageUrl:    "https://upload.wikimedia.org/wikipedia/en/a/a1/The_Two_Towers_cover.gif",
	}

	mediumId := ctx.CreateTestMediumCustom(t, testBook)
	wrongID := pgtype.UUID{Valid: false}
	ctx.CreateTestMediumCustom(t, testBook2)

	testMethod := "PUT"
	testEndpoint := ctx.BaseURL + "/api/media"

	tests := []struct {
		name           string
		requestHeaders map[string]string
		requestBody    parametersUpdateMedium
		expectedStatus int
		expectResponse bool
		checkResponse  func(*testing.T, Medium)
		checkAfter     func(*testing.T)
	}{
		{
			name: "Valid",
			requestHeaders: map[string]string{
				"Authorization": fmt.Sprintf("Bearer %s", ctx.UserAcessToken),
			},
			requestBody: parametersUpdateMedium{
				ID:          mediumId,
				Title:       "The Fellowship of the Ring",
				MediaType:   "book",
				Creator:     "John Ronald Reuel Tolkien",
				ReleaseYear: 1954,
				ImageUrl:    "https://upload.wikimedia.org/wikipedia/en/thumb/8/8e/The_Fellowship_of_the_Ring_cover.gif/220px-The_Fellowship_of_the_Ring_cover.gif",
			},
			expectedStatus: 200,
			expectResponse: true,
			checkResponse: func(t *testing.T, m Medium) {
				if m.Title != testBook.Title {
					t.Error("Invalid 'title' field")
				}
				if m.UpdatedAt == m.CreatedAt {
					t.Error("'updated_at' field not updated")
				}
				if m.Creator != "John Ronald Reuel Tolkien" {
					t.Error("'creator' field not updated")
				}
			},
		},
		{
			name:           "No access_token",
			expectedStatus: 401,
		},
		{
			name: "Title duplicate",
			requestHeaders: map[string]string{
				"Authorization": fmt.Sprintf("Bearer %s", ctx.UserAcessToken),
			},
			requestBody: parametersUpdateMedium{
				ID:          mediumId,
				Title:       "The Two Towers",
				MediaType:   "book",
				Creator:     "John Ronald Reuel Tolkien",
				ReleaseYear: 1954,
				ImageUrl:    "https://upload.wikimedia.org/wikipedia/en/thumb/8/8e/The_Fellowship_of_the_Ring_cover.gif/220px-The_Fellowship_of_the_Ring_cover.gif",
			},
			expectedStatus: 409,
		},
		{
			name: "Wrong medium ID",
			requestHeaders: map[string]string{
				"Authorization": fmt.Sprintf("Bearer %s", ctx.UserAcessToken),
			},
			requestBody: parametersUpdateMedium{
				ID:          wrongID,
				Title:       "The Fellowship of the Ring",
				MediaType:   "book",
				Creator:     "John Ronald Reuel Tolkien",
				ReleaseYear: 1954,
				ImageUrl:    "https://upload.wikimedia.org/wikipedia/en/thumb/8/8e/The_Fellowship_of_the_Ring_cover.gif/220px-The_Fellowship_of_the_Ring_cover.gif",
			},
			expectedStatus: 404,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			requestBody, _ := json.Marshal(tc.requestBody)
			req, _ := http.NewRequest(testMethod, testEndpoint, bytes.NewBuffer(requestBody))
			req.Header.Set("Content-Type", "application/json")
			if tc.requestHeaders != nil {
				for headerKey, headerValue := range tc.requestHeaders {
					req.Header.Set(headerKey, headerValue)
				}
			}
			resp, err := ctx.Client.Do(req)
			if err != nil {
				t.Fatalf("Failed to send request: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tc.expectedStatus {
				t.Errorf("Expected status code %d, got %d", tc.expectedStatus, resp.StatusCode)
			}

			if tc.expectResponse {
				var responseBody Medium
				err := json.NewDecoder(resp.Body).Decode(&responseBody)
				if err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}
				if tc.checkResponse != nil {
					tc.checkResponse(t, responseBody)
				}
			}
			if tc.checkAfter != nil {
				tc.checkAfter(t)
			}
		})
	}
}

func TestDeleteMedium(t *testing.T) {
	ctx := SetupTestContext(t)
	defer ctx.Server.Shutdown(context.Background())

	ctx.CreateTestUser(t)
	ctx.LoginTestUser(t)

	testBook := parametersCreateMedium{
		Title:       "The Fellowship of the Ring",
		MediaType:   "book",
		Creator:     "J.R.R Tolkien",
		ReleaseYear: 1954,
		ImageUrl:    "https://upload.wikimedia.org/wikipedia/en/8/8e/The_Fellowship_of_the_Ring_cover.gif",
	}
	mediumID := ctx.CreateTestMediumCustom(t, testBook)

	testMethod := "DELETE"
	testEndpoint := ctx.BaseURL + "/api/media"

	tests := []struct {
		name           string
		requestHeaders map[string]string
		requestBody    parametersDeleteMedium
		expectedStatus int
		expectResponse bool
		checkResponse  func(*testing.T, map[string]string)
		checkAfter     func(*testing.T)
	}{
		{
			name: "Valid",
			requestHeaders: map[string]string{
				"Authorization": fmt.Sprintf("Bearer %s", ctx.UserAcessToken),
			},
			requestBody: parametersDeleteMedium{
				MediumID: mediumID,
			},
			expectedStatus: 200,
			checkAfter: func(t *testing.T) {
				if ctx.TestIfMediumExist(mediumID) {
					t.Error("Medium still exists in database")
				}
			},
		},
		{
			name:           "No access_token",
			expectedStatus: 401,
		},
		{
			name: "Wrong medium ID (already deleted)",
			requestHeaders: map[string]string{
				"Authorization": fmt.Sprintf("Bearer %s", ctx.UserAcessToken),
			},
			requestBody: parametersDeleteMedium{
				MediumID: mediumID,
			},
			expectedStatus: 404,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			requestBody, _ := json.Marshal(tc.requestBody)
			req, _ := http.NewRequest(testMethod, testEndpoint, bytes.NewBuffer(requestBody))
			req.Header.Set("Content-Type", "application/json")
			if tc.requestHeaders != nil {
				for headerKey, headerValue := range tc.requestHeaders {
					req.Header.Set(headerKey, headerValue)
				}
			}
			resp, err := ctx.Client.Do(req)
			if err != nil {
				t.Fatalf("Failed to send request: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tc.expectedStatus {
				t.Errorf("Expected status code %d, got %d", tc.expectedStatus, resp.StatusCode)
			}

			if tc.expectResponse {
				var responseBody map[string]string
				err := json.NewDecoder(resp.Body).Decode(&responseBody)
				if err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}
				if tc.checkResponse != nil {
					tc.checkResponse(t, responseBody)
				}
			}
			if tc.checkAfter != nil {
				tc.checkAfter(t)
			}
		})
	}
}

/*
============================
TESTS FOR RECORDS ENDPOINTS
============================
*/

func TestCreateRecord(t *testing.T) {
	ctx := SetupTestContext(t)
	defer ctx.Server.Shutdown(context.Background())

	ctx.CreateTestUser(t)
	ctx.LoginTestUser(t)

	medium1Id := ctx.CreateTestMediumRandom(t)
	medium2Id := ctx.CreateTestMediumRandom(t)
	medium3Id := ctx.CreateTestMediumRandom(t)
	medium4Id := ctx.CreateTestMediumRandom(t)
	medium5Id := ctx.CreateTestMediumRandom(t)

	dateNow := pgtype.Timestamp{
		Valid: true,
		Time:  time.Now().UTC(),
	}
	dateFuture := pgtype.Timestamp{
		Valid: true,
		Time:  time.Now().UTC().AddDate(0, 0, 21),
	}
	datePast := pgtype.Timestamp{
		Valid: true,
		Time:  time.Now().UTC().AddDate(0, 0, -21),
	}
	dateInvalid := pgtype.Timestamp{
		Valid: false,
	}

	testMethod := "POST"
	testEndpoint := ctx.BaseURL + "/api/records"

	tests := []struct {
		name           string
		requestHeaders map[string]string
		requestBody    parametersCreateUserMediumRecord
		expectedStatus int
		expectResponse bool
		checkResponse  func(*testing.T, ClientRecord)
		checkAfter     func(*testing.T)
	}{
		{
			name: "Valid, start date and end date",
			requestHeaders: map[string]string{
				"Authorization": fmt.Sprintf("Bearer %s", ctx.UserAcessToken),
			},
			requestBody: parametersCreateUserMediumRecord{
				MediaID:   medium1Id,
				StartDate: dateNow,
				EndDate:   dateFuture,
			},
			expectedStatus: 201,
			expectResponse: true,
			checkResponse: func(t *testing.T, r ClientRecord) {
				if r.ID == "" {
					t.Error("'id' response field missing")
				}
				if r.CreatedAt == "" {
					t.Error("'created_at' response field missing")
				}
				if r.UpdatedAt == "" {
					t.Error("'updated_at' response field missing")
				}
				if r.UserID != ctx.UserID.String() {
					t.Error("'user_id' response field incorrect")
				}
				if r.MediaID != medium1Id.String() {
					t.Error("'media_id' response field incorrect")
				}
				if !r.IsFinished {
					t.Error("'is_finished' response field incorrect")
				}
				if r.StartDate == "" {
					t.Error("'start_date' response field missing")
				}
				if r.EndDate == "" {
					t.Error("'end_date_date' response field missing")
				}
				if r.Duration != 21 {
					t.Error("'duration' response field incorrect")
				}
			},
		},
		{
			name: "Valid, no end date",
			requestHeaders: map[string]string{
				"Authorization": fmt.Sprintf("Bearer %s", ctx.UserAcessToken),
			},
			requestBody: parametersCreateUserMediumRecord{
				MediaID:   medium2Id,
				StartDate: dateNow,
				EndDate:   dateInvalid,
			},
			expectedStatus: 201,
			expectResponse: true,
			checkResponse: func(t *testing.T, r ClientRecord) {
				if r.ID == "" {
					t.Error("'id' response field missing")
				}
				if r.StartDate == "" {
					t.Error("'start_date' response field missing")
				}
				if r.EndDate != "" {
					t.Error("'end_date' response field incorrect")
				}
				if r.IsFinished {
					t.Error("'is_finished' response field incorrect")
				}
				if r.Duration != 0 {
					t.Error("'duration' response field incorrect")
				}
			},
		},
		{
			name: "Valid, no start date",
			requestHeaders: map[string]string{
				"Authorization": fmt.Sprintf("Bearer %s", ctx.UserAcessToken),
			},
			requestBody: parametersCreateUserMediumRecord{
				MediaID:   medium3Id,
				StartDate: dateInvalid,
				EndDate:   dateNow,
			},
			expectedStatus: 201,
			expectResponse: true,
			checkResponse: func(t *testing.T, r ClientRecord) {
				if r.ID == "" {
					t.Error("'id' response field missing")
				}
				if r.StartDate != "" {
					t.Error("'start_date' response field incorrect")
				}
				if r.EndDate == "" {
					t.Error("'end_date' response field missing")
				}
				if r.IsFinished {
					t.Error("'is_finished' response field incorrect")
				}
				if r.Duration != 0 {
					t.Error("'duration' response field incorrect")
				}
			},
		},
		{
			name: "Valid, no date",
			requestHeaders: map[string]string{
				"Authorization": fmt.Sprintf("Bearer %s", ctx.UserAcessToken),
			},
			requestBody: parametersCreateUserMediumRecord{
				MediaID:   medium4Id,
				StartDate: dateInvalid,
				EndDate:   dateInvalid,
			},
			expectedStatus: 201,
			expectResponse: true,
			checkResponse: func(t *testing.T, r ClientRecord) {
				if r.ID == "" {
					t.Error("'id' response field missing")
				}
				if r.StartDate != "" {
					t.Error("'start_date' response field incorrect")
				}
				if r.EndDate != "" {
					t.Error("'end_date' response field incorrect")
				}
				if r.IsFinished {
					t.Error("'is_finished' response field incorrect")
				}
				if r.Duration != 0 {
					t.Error("'duration' response field incorrect")
				}
			},
		},
		{
			name:           "No access_token",
			expectedStatus: 401,
		},
		{
			name: "Invalid, end date in past",
			requestHeaders: map[string]string{
				"Authorization": fmt.Sprintf("Bearer %s", ctx.UserAcessToken),
			},
			requestBody: parametersCreateUserMediumRecord{
				MediaID:   medium5Id,
				StartDate: dateNow,
				EndDate:   datePast,
			},
			expectedStatus: 400,
		},
		{
			name: "Conflict, try to create Record with same user-medium couple",
			requestHeaders: map[string]string{
				"Authorization": fmt.Sprintf("Bearer %s", ctx.UserAcessToken),
			},
			requestBody: parametersCreateUserMediumRecord{
				MediaID:   medium1Id,
				StartDate: dateNow,
				EndDate:   dateFuture,
			},
			expectedStatus: 409,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			requestBody, _ := json.Marshal(tc.requestBody)
			req, _ := http.NewRequest(testMethod, testEndpoint, bytes.NewBuffer(requestBody))
			req.Header.Set("Content-Type", "application/json")
			if tc.requestHeaders != nil {
				for headerKey, headerValue := range tc.requestHeaders {
					req.Header.Set(headerKey, headerValue)
				}
			}
			resp, err := ctx.Client.Do(req)
			if err != nil {
				t.Fatalf("Failed to send request: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tc.expectedStatus {
				t.Errorf("Expected status code %d, got %d", tc.expectedStatus, resp.StatusCode)
			}

			if tc.expectResponse {
				var responseBody ClientRecord
				err := json.NewDecoder(resp.Body).Decode(&responseBody)
				if err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}
				if tc.checkResponse != nil {
					tc.checkResponse(t, responseBody)
				}
			}
			if tc.checkAfter != nil {
				tc.checkAfter(t)
			}
		})
	}
}

func TestGetRecordsByUserID(t *testing.T) {
	ctx := SetupTestContext(t)
	defer ctx.Server.Shutdown(context.Background())

	ctx.CreateTestUser(t)
	ctx.LoginTestUser(t)
	medium1Id := ctx.CreateTestMediumRandom(t)
	medium2Id := ctx.CreateTestMediumRandom(t)
	record1ID := ctx.CreateTestRecord(t, medium1Id)
	record2ID := ctx.CreateTestRecord(t, medium2Id)

	testMethod := "GET"
	testEndpoint := ctx.BaseURL + "/api/records"

	tests := []struct {
		name           string
		requestHeaders map[string]string
		requestBody    map[string]string
		expectedStatus int
		expectResponse bool
		checkResponse  func(*testing.T, ClientRecords)
		checkAfter     func(*testing.T)
	}{
		{
			name: "Valid, with two records",
			requestHeaders: map[string]string{
				"Authorization": fmt.Sprintf("Bearer %s", ctx.UserAcessToken),
			},
			expectedStatus: 200,
			expectResponse: true,
			checkResponse: func(t *testing.T, cr ClientRecords) {
				for _, r := range cr.Records {
					if r.ID != record1ID.String() && r.ID != record2ID.String() {
						t.Error("'id' response field missing")
					}
					if r.CreatedAt == "" {
						t.Error("'created_at' response field missing")
					}
					if r.UpdatedAt == "" {
						t.Error("'updated_at' response field missing")
					}
					if r.UserID != ctx.UserID.String() {
						t.Error("'user_id' response field incorrect")
					}
					if r.MediaID != medium1Id.String() && r.MediaID != medium2Id.String() {
						t.Error("'media_id' response field incorrect")
					}
					if !r.IsFinished {
						t.Error("'is_finished' response field incorrect")
					}
					if r.StartDate == "" {
						t.Error("'start_date' response field missing")
					}
					if r.EndDate == "" {
						t.Error("'end_date_date' response field missing")
					}
					if r.Duration != 21 {
						t.Error("'duration' response field incorrect")
					}
				}
			},
		},
		{
			name:           "No access_token",
			expectedStatus: 401,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			requestBody, _ := json.Marshal(tc.requestBody)
			req, _ := http.NewRequest(testMethod, testEndpoint, bytes.NewBuffer(requestBody))
			req.Header.Set("Content-Type", "application/json")
			if tc.requestHeaders != nil {
				for headerKey, headerValue := range tc.requestHeaders {
					req.Header.Set(headerKey, headerValue)
				}
			}
			resp, err := ctx.Client.Do(req)
			if err != nil {
				t.Fatalf("Failed to send request: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tc.expectedStatus {
				t.Errorf("Expected status code %d, got %d", tc.expectedStatus, resp.StatusCode)
			}

			if tc.expectResponse {
				var responseBody ClientRecords
				err := json.NewDecoder(resp.Body).Decode(&responseBody)
				if err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}
				if tc.checkResponse != nil {
					tc.checkResponse(t, responseBody)
				}
			}
			if tc.checkAfter != nil {
				tc.checkAfter(t)
			}
		})
	}
}

func TestUpdateRecord(t *testing.T) {
	ctx := SetupTestContext(t)
	defer ctx.Server.Shutdown(context.Background())

	ctx.CreateTestUser(t)
	ctx.LoginTestUser(t)
	medium1ID := ctx.CreateTestMediumRandom(t)
	record1ID := ctx.CreateTestRecord(t, medium1ID)

	testMethod := "PUT"
	testEndpoint := ctx.BaseURL + "/api/records"

	tests := []struct {
		name           string
		requestHeaders map[string]string
		requestBody    parametersUpdateRecord
		expectedStatus int
		expectResponse bool
		checkResponse  func(*testing.T, ClientRecord)
		checkAfter     func(*testing.T)
	}{
		{
			name: "Valid, changed start date",
			requestHeaders: map[string]string{
				"Authorization": fmt.Sprintf("Bearer %s", ctx.UserAcessToken),
			},
			requestBody: parametersUpdateRecord{
				RecordID: record1ID,
				StartDate: pgtype.Timestamp{
					Time:  time.Now().AddDate(0, 0, 11),
					Valid: true,
				},
			},
			expectedStatus: 200,
			expectResponse: true,
			checkResponse: func(t *testing.T, cr ClientRecord) {
				if cr.ID == "" {
					t.Error("'id' response field missing")
				}
				if cr.CreatedAt == "" {
					t.Error("'created_at' response field missing")
				}
				if cr.UpdatedAt == "" {
					t.Error("'updated_at' response field missing")
				}
				if cr.UserID != ctx.UserID.String() {
					t.Error("'user_id' response field incorrect")
				}
				if cr.MediaID != medium1ID.String() {
					t.Error("'media_id' response field incorrect")
				}
				if !cr.IsFinished {
					t.Error("'is_finished' response field incorrect")
				}
				if cr.StartDate == "" {
					t.Error("'start_date' response field missing")
				}
				if cr.EndDate == "" {
					t.Error("'end_date_date' response field missing")
				}
				if cr.Duration != 9 {
					t.Error("'duration' response field incorrect.")
				}
			},
		},
		{
			name:           "No access_token",
			expectedStatus: 401,
		},
		{
			name: "New start date is after end date",
			requestHeaders: map[string]string{
				"Authorization": fmt.Sprintf("Bearer %s", ctx.UserAcessToken),
			},
			requestBody: parametersUpdateRecord{
				RecordID: record1ID,
				StartDate: pgtype.Timestamp{
					Time:  time.Now().AddDate(0, 0, 30),
					Valid: true,
				},
			},
			expectedStatus: 400,
		},
		{
			name: "Invalid record ID",
			requestHeaders: map[string]string{
				"Authorization": fmt.Sprintf("Bearer %s", ctx.UserAcessToken),
			},
			requestBody: parametersUpdateRecord{
				RecordID: pgtype.UUID{Valid: false},
			},
			expectedStatus: 404,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			requestBody, _ := json.Marshal(tc.requestBody)
			req, _ := http.NewRequest(testMethod, testEndpoint, bytes.NewBuffer(requestBody))
			req.Header.Set("Content-Type", "application/json")
			if tc.requestHeaders != nil {
				for headerKey, headerValue := range tc.requestHeaders {
					req.Header.Set(headerKey, headerValue)
				}
			}
			resp, err := ctx.Client.Do(req)
			if err != nil {
				t.Fatalf("Failed to send request: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tc.expectedStatus {
				t.Errorf("Expected status code %d, got %d", tc.expectedStatus, resp.StatusCode)
			}

			if tc.expectResponse {
				var responseBody ClientRecord
				err := json.NewDecoder(resp.Body).Decode(&responseBody)
				if err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}
				if tc.checkResponse != nil {
					tc.checkResponse(t, responseBody)
				}
			}
			if tc.checkAfter != nil {
				tc.checkAfter(t)
			}
		})
	}
}

func TestDeleteRecord(t *testing.T) {
	ctx := SetupTestContext(t)
	defer ctx.Server.Shutdown(context.Background())

	ctx.CreateTestUser(t)
	ctx.LoginTestUser(t)
	mediumID := ctx.CreateTestMediumRandom(t)
	recordID := ctx.CreateTestRecord(t, mediumID)

	testMethod := "DELETE"
	testEndpoint := ctx.BaseURL + "/api/records"

	tests := []struct {
		name           string
		requestHeaders map[string]string
		requestBody    parametersDeleteRecord
		expectedStatus int
		expectResponse bool
		checkResponse  func(*testing.T, map[string]string)
		checkAfter     func(*testing.T)
	}{
		{
			name: "Valid",
			requestHeaders: map[string]string{
				"Authorization": fmt.Sprintf("Bearer %s", ctx.UserAcessToken),
			},
			requestBody: parametersDeleteRecord{
				RecordID: recordID,
			},
			expectedStatus: 200,
			expectResponse: false,
			checkAfter: func(t *testing.T) {
				if ctx.TestIfRecordExist(recordID) {
					t.Error("record still exists")
				}
			},
		},
		{
			name:           "No access_token",
			expectedStatus: 401,
		},
		{
			name: "Invalid record ID",
			requestHeaders: map[string]string{
				"Authorization": fmt.Sprintf("Bearer %s", ctx.UserAcessToken),
			},
			requestBody: parametersDeleteRecord{
				RecordID: pgtype.UUID{Valid: false},
			},
			expectedStatus: 404,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			requestBody, _ := json.Marshal(tc.requestBody)
			req, _ := http.NewRequest(testMethod, testEndpoint, bytes.NewBuffer(requestBody))
			req.Header.Set("Content-Type", "application/json")
			if tc.requestHeaders != nil {
				for headerKey, headerValue := range tc.requestHeaders {
					req.Header.Set(headerKey, headerValue)
				}
			}
			resp, err := ctx.Client.Do(req)
			if err != nil {
				t.Fatalf("Failed to send request: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tc.expectedStatus {
				t.Errorf("Expected status code %d, got %d", tc.expectedStatus, resp.StatusCode)
			}

			if tc.expectResponse {
				var responseBody map[string]string
				err := json.NewDecoder(resp.Body).Decode(&responseBody)
				if err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}
				if tc.checkResponse != nil {
					tc.checkResponse(t, responseBody)
				}
			}
			if tc.checkAfter != nil {
				tc.checkAfter(t)
			}
		})
	}
}
