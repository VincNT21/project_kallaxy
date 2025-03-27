package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
	"time"
)

func TestCreateUser(t *testing.T) {
	testEnv := map[string]string{
		"DB_URL": "postgres://postgres:postgres@localhost:5432/kallaxy?sslmode=disable",
		"SECRET": "oyD+hk8wMswsng1nZ5n8IQEbHbW68+s6dyuqaBs4i3p0aX647253DCy9ldd0lcxZJfgVvFbZ1ufvcUhlK115TQ==",
	}
	// Start server in a goroutine
	go func() {
		Start(testEnv)
	}()

	// Give the server a moment to start
	time.Sleep(100 * time.Millisecond)

	// Create tests table
	tests := []struct {
		name           string
		requestBody    map[string]string
		expectedStatus int
		checkResponse  func(*testing.T, map[string]interface{})
	}{
		{
			name: "Valid user creation",
			requestBody: map[string]string{
				"username": "Testuser1",
				"password": "12345678",
				"email":    "Test123@example.com",
			},
			expectedStatus: 201,
			checkResponse: func(t *testing.T, body map[string]interface{}) {
				if _, exists := body["id"]; !exists {
					t.Error("Reponse missing 'id' field")
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
			req, _ := http.NewRequest("POST", "http://localhost:8080/api/users", bytes.NewBuffer(requestBody))
			req.Header.Set("Content-Type", "application/json")

			// Send the request
			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				t.Fatalf("Failed to send request: %v", err)
			}
			defer resp.Body.Close()

			// Check status code
			if resp.StatusCode != tc.expectedStatus {
				t.Errorf("Expected status code %d, got %d", tc.expectedStatus, resp.StatusCode)
			}

			// If JSON response expected, check the response body
			if resp.StatusCode == 201 {
				var responseBody map[string]interface{}
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
