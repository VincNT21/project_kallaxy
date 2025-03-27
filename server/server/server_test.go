package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
)

// All comments are on first test only

func TestCreateUser(t *testing.T) {
	// Setup test environnement
	ctx := SetupTestContext(t)
	defer ctx.Server.Shutdown(context.Background()) // Clean shutdown when test ends

	// Define test method and endpoint
	testMethod := "POST"
	testEndpoint := ctx.BaseURL + "/api/users"

	// Create tests table
	tests := []struct {
		name           string
		requestBody    map[string]string
		expectedStatus int
		checkResponse  func(*testing.T, map[string]string)
	}{
		{
			name: "Valid user creation",
			requestBody: map[string]string{
				"username": "Testuser1",
				"password": "12345678",
				"email":    "Test123@example.com",
			},
			expectedStatus: 201,
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
		})
	}
}

func TestLogin(t *testing.T) {
	ctx := SetupTestContext(t)
	defer ctx.Server.Shutdown(context.Background())

	ctx.CreateTestUser(t)

	testMethod := "POST"
	testEndpoint := ctx.BaseURL + "/auth/login"

	tests := []struct {
		name           string
		requestBody    map[string]string
		expectedStatus int
		checkResponse  func(*testing.T, map[string]string)
	}{
		{
			name: "Valid login",
			requestBody: map[string]string{
				"username": ctx.UserUsername,
				"password": ctx.UserPassword,
			},
			expectedStatus: 201,
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

	tests := []struct {
		name            string
		requestHeaders  map[string]string
		requestBody     map[string]string
		expectedStatus  int
		checkResponse   func(*testing.T, map[string]string)
		checkRevokation func(*testing.T)
	}{
		{
			name: "Logout ok",
			requestHeaders: map[string]string{
				"Authorization": fmt.Sprintf("Bearer %s", ctx.UserAcessToken),
			},
			expectedStatus: 204,
			checkRevokation: func(t *testing.T) {
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
			if tc.checkRevokation != nil {
				tc.checkRevokation(t)
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

	tests := []struct {
		name            string
		requestHeaders  map[string]string
		requestBody     map[string]string
		expectedStatus  int
		checkResponse   func(*testing.T, map[string]string)
		checkRevokation func(*testing.T)
	}{
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
			checkRevokation: func(t *testing.T) {
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
			if tc.checkRevokation != nil {
				tc.checkRevokation(t)
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

	tests := []struct {
		name           string
		requestHeaders map[string]string
		requestBody    map[string]string
		expectedStatus int
		checkResponse  func(*testing.T, map[string]string)
	}{
		{
			name: "Valid Refresh Token",
			requestHeaders: map[string]string{
				"Authorization": fmt.Sprintf("Bearer %s", ctx.UserRefreshToken),
			},
			expectedStatus: 201,
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

	tests := []struct {
		name            string
		requestHeaders  map[string]string
		requestBody     map[string]string
		expectedStatus  int
		checkResponse   func(*testing.T, map[string]string)
		checkRevokation func(*testing.T)
	}{
		{
			name: "Revoke ok",
			requestHeaders: map[string]string{
				"Authorization": fmt.Sprintf("Bearer %s", ctx.UserRefreshToken),
			},
			expectedStatus: 204,
			checkRevokation: func(t *testing.T) {
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
			if tc.checkRevokation != nil {
				tc.checkRevokation(t)
			}
		})
	}
}

/*====== MODEL FOR TESTS========
func TestXxx(t *testing.T) {
	ctx := SetupTestContext(t)
	defer ctx.Server.Shutdown(context.Background())

	ctx.CreateTestUser(t)
	ctx.LoginTestUser(t)

	testMethod := "xxxxxxx"
	testEndpoint := ctx.BaseURL + "xxxxxxxxxx"

	tests := []struct {
		name           string
		requestHeaders map[string]string
		requestBody    map[string]string
		expectedStatus int
		checkResponse  func(*testing.T, map[string]string)
	}{

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
		})
	}
}
===============================*/
