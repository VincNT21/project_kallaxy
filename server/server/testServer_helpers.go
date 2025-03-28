package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/VincNT21/kallaxy/server/internal/auth"
	"github.com/VincNT21/kallaxy/server/internal/database"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TestContext struct {
	Server           *http.Server
	BaseURL          string
	UserAcessToken   string
	UserRefreshToken string
	UserUsername     string
	UserPassword     string
	UserEmail        string
	Client           *http.Client
}

// Setup creates a clean test environment and returns a TestContext
func SetupTestContext(t *testing.T) *TestContext {

	// Start test server with dynamic port
	server, baseURL := setupTestServer(t)

	// Create a test context
	ctx := &TestContext{
		Server:       server,
		BaseURL:      baseURL,
		UserUsername: "TestUser",
		UserPassword: "azerty1234",
		UserEmail:    "test@example.com",
		Client:       &http.Client{},
	}

	// Reset the database to start with a clean slate
	ctx.ResetDatabase(t)

	return ctx
}

// Create a test server that behave identically to production server
func setupTestServer(t *testing.T) (*http.Server, string) {
	// Create a listener on a random avalaible port
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("Failed to create listener: %v", err)
	}

	// Get the dynamic port that was assigned
	port := listener.Addr().(*net.TCPAddr).Port
	serverURL := fmt.Sprintf("http://127.0.0.1:%d", port)

	// Set up test environnement variables
	testEnv := map[string]string{
		"DB_URL": "postgresql://postgres:postgres@localhost:5432/kallaxytestdb",
		"SECRET": "test-jwt-secret",
	}

	// Open a connection to database
	dbConnection, err := pgxpool.New(context.Background(), testEnv["DB_URL"])
	if err != nil {
		log.Fatalf("Failed to connect to database, %v", err)
	}

	// Create a *database.Queries
	db := database.New(dbConnection)

	// Init apiCfg
	apiCfg := newAPIConfig(db, testEnv["SECRET"])

	// Delete revoked refresh token in database
	apiCfg.CleanRefreshTokens()

	// Create the request multiplexer (router)
	mux := http.NewServeMux()

	// Register handlers
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset) // Reset handler is only used on test server

	mux.HandleFunc("POST /api/users", apiCfg.handlerCreateUser)
	mux.Handle("PUT /api/users", apiCfg.authMiddleware(http.HandlerFunc(apiCfg.handlerUpdateUser)))

	mux.HandleFunc("POST /auth/refresh", apiCfg.handlerRefresh)
	mux.HandleFunc("POST /auth/revoke", apiCfg.handlerRevoke)

	mux.HandleFunc("POST /auth/login", apiCfg.handlerLogin)
	mux.Handle("POST /auth/logout", apiCfg.authMiddleware(http.HandlerFunc(apiCfg.handlerLogout)))

	mux.HandleFunc("POST /auth/password-reset", apiCfg.handlerPasswordResetRequest)
	mux.HandleFunc("GET /auth/password-reset", apiCfg.handlerVerifyResetToken)
	mux.HandleFunc("PUT /auth/password-reset", apiCfg.handlerResetPassword)

	// Create a http server with our multiplexer
	server := &http.Server{
		Handler: mux,
	}

	// Start server in a goroutine
	go func() {
		if err := server.Serve(listener); err != nil && err != http.ErrServerClosed {
			t.Logf("Server error: %v", err)
		}
	}()

	// Small delay to ensure server is up
	time.Sleep(100 * time.Millisecond)

	// return server and URL so tests can use it
	return server, serverURL
}

// ResetDatabase sends a request to the admin Reset endpoint
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

// CreateTestUser creates a user
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

// LoginTestUser logs in a user and stores tokens in context variables
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
	var tokens map[string]string
	err = json.NewDecoder(resp.Body).Decode(&tokens)
	if err != nil {
		t.Fatalf("Failed to decode login response body: %v", err)
	}
	ctx.UserAcessToken = tokens["access_token"]
	ctx.UserRefreshToken = tokens["refresh_token"]
}

func TestValidateAccessToken(token string) bool {
	_, err := auth.ValidateJWT(token, "test-jwt-secret")
	return err == nil
}

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
