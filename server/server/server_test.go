package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"
)

// Tests use dev_test_server, which is a copy of production server (keeped up to date)

// All comments are on the first test only

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
		Email:    "Test123@example.com",
	}

	tests := []struct {
		name           string
		requestHeaders map[string]string
		requestBody    parametersCreateUser
		expectedStatus int
		expectResponse bool
		checkResponse  func(*testing.T, ClientUser)
	}{
		{
			name:           "Valid user creation",
			requestBody:    testUser,
			expectedStatus: 201,
			expectResponse: true,
			checkResponse: func(t *testing.T, u ClientUser) {
				if u.ID == "" {
					t.Error("Response missing 'id' field")
				}
				if u.CreatedAt == "" {
					t.Error("Response missing 'created_at' field")
				}
				if u.UpdatedAt == "" {
					t.Error("Response missing 'updated_at' field")
				}
				if u.Username != "TestUser1" {
					t.Error("Response have incorrect 'username' field")
				}
				if u.Email != "Test123@example.com" {
					t.Error("Response have incorrect 'email' field")
				}
			},
		},
		{
			name: "Duplicate user's username",
			requestBody: parametersCreateUser{
				Username: "TestUser1",
				Password: "hjldsfoeri",
				Email:    "Test123@example2.com",
			},
			expectedStatus: 409,
		},
		{
			name: "Duplicate user's email",
			requestBody: parametersCreateUser{
				Username: "Testuser2",
				Password: "hjldsfoeri",
				Email:    "Test123@example.com",
			},
			expectedStatus: 409,
		},
		{
			name: "Missing a field",
			requestBody: parametersCreateUser{
				Username: "Testuser2",
				Email:    "Test123@example.com",
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
				var responseBody ClientUser
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

func TestGetUser(t *testing.T) {
	ctx := SetupTestContext(t)
	defer ctx.Server.Shutdown(context.Background())

	ctx.CreateTestUser(t)
	ctx.LoginTestUser(t)

	testMethod := "GET"
	testEndpoint := ctx.BaseURL + "/api/users"

	tests := []struct {
		name           string
		requestHeaders map[string]string
		requestBody    map[string]string
		expectedStatus int
		expectResponse bool
		checkResponse  func(*testing.T, ClientUser)
	}{
		{
			name: "Valid",
			requestHeaders: map[string]string{
				"Authorization": fmt.Sprintf("Bearer %s", ctx.UserAcessToken),
			},
			expectedStatus: 200,
			expectResponse: true,
			checkResponse: func(t *testing.T, u ClientUser) {
				if u.ID == "" {
					t.Error("Response missing 'id' field")
				}
				if u.CreatedAt == "" {
					t.Error("Response missing 'created_at' field")
				}
				if u.UpdatedAt == "" {
					t.Error("Response missing 'updated_at' field")
				}
				if u.Username != ctx.UserUsername {
					t.Error("Response have incorrect 'username' field")
				}
				if u.Email != ctx.UserEmail {
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
				var responseBody ClientUser
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

func TestUpdateUser(t *testing.T) {
	ctx := SetupTestContext(t)
	defer ctx.Server.Shutdown(context.Background())

	ctx.CreateTestUser(t)
	ctx.LoginTestUser(t)

	testMethod := "PUT"
	testEndpoint := ctx.BaseURL + "/api/users"

	tests := []struct {
		name           string
		requestHeaders map[string]string
		requestBody    parametersCreateUser
		expectedStatus int
		expectResponse bool
		checkResponse  func(*testing.T, ClientUser)
		checkAfter     func(*testing.T)
	}{
		{
			name: "Valid token and valid data",
			requestHeaders: map[string]string{
				"Authorization": fmt.Sprintf("Bearer %s", ctx.UserAcessToken),
			},
			requestBody: parametersCreateUser{
				Username: ctx.UserUsername,
				Password: ctx.UserPassword,
				Email:    "newemail@example.com",
			},
			expectedStatus: 200,
			expectResponse: true,
			checkResponse: func(t *testing.T, u ClientUser) {
				if u.ID == "" {
					t.Error("Response missing 'id' field")
				}
				if u.CreatedAt == "" {
					t.Error("Response field 'created_at_at' was not updated")
				}
				if u.UpdatedAt == "" || u.UpdatedAt == u.CreatedAt {
					t.Error("Response field 'updated_at' missing or not updated")
				}
				if u.Username != ctx.UserUsername {
					t.Error("Response have incorrect 'username' field")
				}
				if u.Email != "newemail@example.com" {
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
			requestBody: parametersCreateUser{
				Username: ctx.UserUsername,
				Password: ctx.UserPassword,
			},
			expectedStatus: 400,
		},
		{
			name:           "Missing token",
			expectedStatus: 401,
		},
		{
			name: "Invalid token",
			requestHeaders: map[string]string{
				"Authorization": "Bearer badaccesstoken",
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
				var responseBody ClientUser
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

	tests := []struct {
		name           string
		requestHeaders map[string]string
		requestBody    map[string]string
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

	tests := []struct {
		name           string
		requestHeaders map[string]string
		requestBody    parametersLogin
		expectedStatus int
		expectResponse bool
		checkResponse  func(*testing.T, ClientTokensAndUser)
	}{
		{
			name: "Valid login",
			requestBody: parametersLogin{
				Username: ctx.UserUsername,
				Password: ctx.UserPassword,
			},
			expectedStatus: 201,
			expectResponse: true,
			checkResponse: func(t *testing.T, r ClientTokensAndUser) {
				if r.ID == "" {
					t.Error("Response missing 'id' field")
				}
				if r.Username != ctx.UserUsername {
					t.Error("Response have incorrect 'username' field")
				}
				if r.AccessToken == "" {
					t.Error("Response missing 'access_token' field")
				}
				if !TestValidateAccessToken(r.AccessToken) {
					t.Error("Response's access_token is invalid")
				}
				if r.RefreshToken == "" {
					t.Error("Response missing 'refresh_token' field")
				}
				if !ctx.TestValidateRefreshToken(r.RefreshToken) {
					t.Error("Response's refresh_token is invalid")
				}
			},
		},
		{
			name: "Invalid password",
			requestBody: parametersLogin{
				Username: ctx.UserUsername,
				Password: "12345",
			},
			expectedStatus: 401,
		},
		{
			name: "Invalid username",
			requestBody: parametersLogin{
				Username: "Testtest12",
				Password: ctx.UserPassword,
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
				var responseBody ClientTokensAndUser
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
		name           string
		requestHeaders map[string]string
		requestBody    map[string]string
		expectedStatus int
		expectResponse bool
		checkResponse  func(*testing.T, map[string]string)
		checkAfter     func(*testing.T)
	}{
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

	tests := []struct {
		name           string
		requestHeaders map[string]string
		requestBody    map[string]string
		expectedStatus int
		expectResponse bool
		checkResponse  func(*testing.T, ClientTokens)
	}{
		{
			name: "Valid Refresh Token",
			requestHeaders: map[string]string{
				"Authorization": fmt.Sprintf("Bearer %s", ctx.UserRefreshToken),
			},
			expectedStatus: 201,
			expectResponse: true,
			checkResponse: func(t *testing.T, ct ClientTokens) {
				if ct.AccessToken == "" {
					t.Error("Response missing 'access_token' field")
				}
				if !TestValidateAccessToken(ct.AccessToken) {
					t.Error("Response's access_token is invalid")
				}
				if ct.RefreshToken == "" {
					t.Error("Response missing 'refresh_token' field")
				}
				if !ctx.TestValidateRefreshToken(ct.RefreshToken) {
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
				var responseBody ClientTokens
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
		name           string
		requestHeaders map[string]string
		requestBody    map[string]string
		expectedStatus int
		expectResponse bool
		checkResponse  func(*testing.T, map[string]string)
		checkAfter     func(*testing.T)
	}{
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
		checkResponse  func(*testing.T, ClientMedium)
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
			checkResponse: func(t *testing.T, r ClientMedium) {
				if r.ID == "" {
					t.Error("'id' response field missing")
				}
				if r.CreatedAt == "" {
					t.Error("'created_at' response field missing")
				}
				if r.UpdatedAt == "" {
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
				if r.ImageUrl != testBook.ImageUrl {
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
				var responseBody ClientMedium
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
		checkResponse  func(*testing.T, ClientMedium)
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
			checkResponse: func(t *testing.T, m ClientMedium) {
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
				if m.ImageUrl != testBook.ImageUrl {
					t.Error("incorrect ImageUrl")
				}
				if m.ID == "" {
					t.Error("'id' response field missing")
				}
				if m.CreatedAt == "" {
					t.Error("'created_at' response field missing")
				}
				if m.UpdatedAt == "" {
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
				var responseBody ClientMedium
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
		checkResponse  func(*testing.T, ClientListMedia)
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
			checkResponse: func(t *testing.T, r ClientListMedia) {
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
					if m.ImageUrl != testBook.ImageUrl && m.ImageUrl != testBook2.ImageUrl {
						t.Error("incorrect ImageUrl")
					}
					if m.ID == "" {
						t.Error("'id' response field missing")
					}
					if m.CreatedAt == "" {
						t.Error("'created_at' response field missing")
					}
					if m.UpdatedAt == "" {
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
				var responseBody ClientListMedia
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
	ctx.CreateTestMediumCustom(t, testBook2)

	testMethod := "PUT"
	testEndpoint := ctx.BaseURL + "/api/media"

	tests := []struct {
		name           string
		requestHeaders map[string]string
		requestBody    parametersUpdateMedium
		expectedStatus int
		expectResponse bool
		checkResponse  func(*testing.T, ClientMedium)
		checkAfter     func(*testing.T)
	}{
		{
			name: "Valid",
			requestHeaders: map[string]string{
				"Authorization": fmt.Sprintf("Bearer %s", ctx.UserAcessToken),
			},
			requestBody: parametersUpdateMedium{
				MediumID:    mediumId,
				Title:       "The Fellowship of the Ring",
				Creator:     "John Ronald Reuel Tolkien",
				ReleaseYear: 1954,
				ImageUrl:    "https://upload.wikimedia.org/wikipedia/en/thumb/8/8e/The_Fellowship_of_the_Ring_cover.gif/220px-The_Fellowship_of_the_Ring_cover.gif",
			},
			expectedStatus: 200,
			expectResponse: true,
			checkResponse: func(t *testing.T, m ClientMedium) {
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
				MediumID:    mediumId,
				Title:       "The Two Towers",
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
				MediumID:    "ba983bd8-36ce-4d1b-ad24-2b65269f9921",
				Title:       "The Fellowship of the Ring",
				Creator:     "John Ronald Reuel Tolkien",
				ReleaseYear: 1954,
				ImageUrl:    "https://upload.wikimedia.org/wikipedia/en/thumb/8/8e/The_Fellowship_of_the_Ring_cover.gif/220px-The_Fellowship_of_the_Ring_cover.gif",
			},
			expectedStatus: 404,
		},
		{
			name: "Malformed medium ID",
			requestHeaders: map[string]string{
				"Authorization": fmt.Sprintf("Bearer %s", ctx.UserAcessToken),
			},
			requestBody: parametersUpdateMedium{
				MediumID:    "ba983bd8-2b65269f9921",
				Title:       "The Fellowship of the Ring",
				Creator:     "John Ronald Reuel Tolkien",
				ReleaseYear: 1954,
				ImageUrl:    "https://upload.wikimedia.org/wikipedia/en/thumb/8/8e/The_Fellowship_of_the_Ring_cover.gif/220px-The_Fellowship_of_the_Ring_cover.gif",
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
				var responseBody ClientMedium
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

	dateNow := time.Now().UTC().Format(time.RFC3339)
	dateFuture := time.Now().UTC().AddDate(0, 0, 21).Format(time.RFC3339)

	datePast := time.Now().UTC().AddDate(0, 0, -21).Format(time.RFC3339)
	dateInvalid := "2024_01_01:21:21:21"

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
				MediumID:  medium1Id,
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
				if r.MediaID != medium1Id {
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
				MediumID:  medium2Id,
				StartDate: dateNow,
				EndDate:   "",
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
				MediumID:  medium3Id,
				StartDate: "",
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
				MediumID:  medium4Id,
				StartDate: "",
				EndDate:   "",
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
				MediumID:  medium5Id,
				StartDate: dateNow,
				EndDate:   datePast,
			},
			expectedStatus: 400,
		},
		{
			name: "Invalid, end date malformed",
			requestHeaders: map[string]string{
				"Authorization": fmt.Sprintf("Bearer %s", ctx.UserAcessToken),
			},
			requestBody: parametersCreateUserMediumRecord{
				MediumID:  medium5Id,
				StartDate: dateNow,
				EndDate:   dateInvalid,
			},
			expectedStatus: 400,
		},
		{
			name: "Conflict, try to create Record with same user-medium couple",
			requestHeaders: map[string]string{
				"Authorization": fmt.Sprintf("Bearer %s", ctx.UserAcessToken),
			},
			requestBody: parametersCreateUserMediumRecord{
				MediumID:  medium1Id,
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
					if r.ID != record1ID && r.ID != record2ID {
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
					if r.MediaID != medium1Id && r.MediaID != medium2Id {
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
				RecordID:  record1ID,
				StartDate: time.Now().AddDate(0, 0, 11).Format(time.RFC3339),
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
				if cr.MediaID != medium1ID {
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
				if cr.Duration != 10 {
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
				RecordID:  record1ID,
				StartDate: time.Now().AddDate(0, 0, 30).Format(time.RFC3339),
			},
			expectedStatus: 400,
		},
		{
			name: "Malformed record ID",
			requestHeaders: map[string]string{
				"Authorization": fmt.Sprintf("Bearer %s", ctx.UserAcessToken),
			},
			requestBody: parametersUpdateRecord{
				RecordID: "wrongID",
			},
			expectedStatus: 400,
		},
		{
			name: "Unknown record ID",
			requestHeaders: map[string]string{
				"Authorization": fmt.Sprintf("Bearer %s", ctx.UserAcessToken),
			},
			requestBody: parametersUpdateRecord{
				RecordID: "ba983bd8-36ce-4d1b-ad24-2b65240f9921",
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
			name: "Unknown record ID",
			requestHeaders: map[string]string{
				"Authorization": fmt.Sprintf("Bearer %s", ctx.UserAcessToken),
			},
			requestBody: parametersDeleteRecord{
				RecordID: "ba983bd8-36ce-4d1b-ad24-2b65240f9921",
			},
			expectedStatus: 404,
		},
		{
			name: "Malformed record ID",
			requestHeaders: map[string]string{
				"Authorization": fmt.Sprintf("Bearer %s", ctx.UserAcessToken),
			},
			requestBody: parametersDeleteRecord{
				RecordID: "wrongID",
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
==================================
TESTS FOR PASSWORD RESET ENDPOINTS
==================================
*/

func TestPasswordResetStep1(t *testing.T) {
	ctx := SetupTestContext(t)
	defer ctx.Server.Shutdown(context.Background())

	ctx.CreateTestUser(t)
	ctx.LoginTestUser(t)

	testMethod := "POST"
	testEndpoint := ctx.BaseURL + "/auth/password_reset"

	tests := []struct {
		name           string
		requestHeaders map[string]string
		requestBody    parametersPasswordResetRequest
		expectedStatus int
		expectResponse bool
		checkResponse  func(*testing.T, responsePasswordResetRequest)
	}{
		{
			name: "Valid",
			requestBody: parametersPasswordResetRequest{
				Email: ctx.UserEmail,
			},
			expectedStatus: 200,
			expectResponse: true,
			checkResponse: func(t *testing.T, rprr responsePasswordResetRequest) {
				if rprr.ResetLink == "" {
					t.Error("invalid 'reset_link' field")
				}
				if rprr.ResetToken == "" {
					t.Error("invalid 'reset_token' field")
				}
				if rprr.Message != "Password reset initiated" {
					t.Error("invalid 'message'")
				}
			},
		},
		{
			name: "Wrong email",
			requestBody: parametersPasswordResetRequest{
				Email: "fakeemail@example.com",
			},
			expectedStatus: 200,
			expectResponse: true,
			checkResponse: func(t *testing.T, rprr responsePasswordResetRequest) {
				if rprr.Message != "If your email exists in our system, you'll receive reset instruction" {
					t.Error("invalid 'message'")
				}
				if rprr.ResetLink != "" {
					t.Error("invalid 'reset_link'")
				}
			},
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
				var responseBody responsePasswordResetRequest
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

func TestPasswordResetStep2(t *testing.T) {
	ctx := SetupTestContext(t)
	defer ctx.Server.Shutdown(context.Background())

	ctx.CreateTestUser(t)
	ctx.LoginTestUser(t)

	resetToken := ctx.GetPasswordResetToken(t)

	testMethod := "GET"
	testEndpoint := ctx.BaseURL + "/auth/password_reset"

	tests := []struct {
		name           string
		requestHeaders map[string]string
		queryParameter string
		requestBody    map[string]string
		expectedStatus int
		expectResponse bool
		checkResponse  func(*testing.T, responseVerifyResetToken)
	}{
		{
			name:           "Valid",
			queryParameter: fmt.Sprintf("?token=%s", resetToken),
			expectedStatus: 200,
			expectResponse: true,
			checkResponse: func(t *testing.T, r responseVerifyResetToken) {
				if r.Email != ctx.UserEmail {
					t.Error("invalid 'email'")
				}
				if !r.Valid {
					t.Error("reset token invalid")
				}
			},
		},
		{
			name:           "Invalid token",
			queryParameter: "?token=wrongtoken",
			expectedStatus: 400,
		},
		{
			name:           "No queryParameter",
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
				var responseBody responseVerifyResetToken
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

func TestPasswordResetStep3(t *testing.T) {
	ctx := SetupTestContext(t)
	defer ctx.Server.Shutdown(context.Background())

	ctx.CreateTestUser(t)
	ctx.LoginTestUser(t)

	resetToken := ctx.GetPasswordResetToken(t)

	testMethod := "PUT"
	testEndpoint := ctx.BaseURL + "/auth/password_reset"

	tests := []struct {
		name           string
		requestHeaders map[string]string
		requestBody    parametersResetPassword
		expectedStatus int
		expectResponse bool
		checkResponse  func(*testing.T, ClientUser)
		checkAfter     func(*testing.T)
	}{
		{
			name: "Valid",
			requestBody: parametersResetPassword{
				Token:       resetToken,
				NewPassword: "cvbn7890",
			},
			expectedStatus: 200,
			expectResponse: true,
			checkResponse: func(t *testing.T, cu ClientUser) {
				if cu.ID == "" {
					t.Error("invalid 'id' field")
				}
				if cu.CreatedAt == "" {
					t.Error("invalid 'created_at' field")
				}
				if cu.UpdatedAt == "" || cu.UpdatedAt == cu.CreatedAt {
					t.Error("invalid 'updated_at' field")
				}
				if cu.Username != ctx.UserUsername {
					t.Error("invalid 'username' field")
				}
				if cu.Email != ctx.UserEmail {
					t.Error("invalid 'email' field")
				}
			},
			checkAfter: func(t *testing.T) {
				if !ctx.TestLoginAfterPasswordReset(t, "cvbn7890") {
					t.Error("new password doesn't work")
				}
			},
		},
		{
			name: "Invalid reset token",
			requestBody: parametersResetPassword{
				Token:       "faketoken",
				NewPassword: "12345678",
			},
			expectedStatus: 400,
			checkAfter: func(t *testing.T) {
				if ctx.TestLoginAfterPasswordReset(t, "12345678") {
					t.Error("password shouldn't have changed")
				}
			},
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
				var responseBody ClientUser
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
