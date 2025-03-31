package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/VincNT21/kallaxy/server/internal/database"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TestContext struct {
	Server           *http.Server
	BaseURL          string
	UserAcessToken   string
	UserRefreshToken string
	UserID           pgtype.UUID
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
		"DB_URL": "postgresql://postgres:postgres@localhost:5432/kallaxytest",
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

	// Users endpoints
	mux.HandleFunc("POST /api/users", apiCfg.handlerCreateUser)
	mux.Handle("GET /api/users", apiCfg.authMiddleware(http.HandlerFunc(apiCfg.handlerGetUserByID)))
	mux.Handle("PUT /api/users", apiCfg.authMiddleware(http.HandlerFunc(apiCfg.handlerUpdateUser)))
	mux.Handle("DELETE /api/users", apiCfg.authMiddleware(http.HandlerFunc(apiCfg.handlerDeleteUser)))

	// Media endpoints
	mux.Handle("POST /api/media", apiCfg.authMiddleware(http.HandlerFunc(apiCfg.handlerCreateMedium)))
	mux.Handle("GET /api/media", apiCfg.authMiddleware(http.HandlerFunc(apiCfg.handlerGetMedia)))
	mux.Handle("PUT /api/media", apiCfg.authMiddleware(http.HandlerFunc(apiCfg.handlerUpdateMedium)))
	mux.Handle("DELETE /api/media", apiCfg.authMiddleware(http.HandlerFunc(apiCfg.handlerDeleteMedium)))

	// Records endpoints
	mux.Handle("POST /api/records", apiCfg.authMiddleware(http.HandlerFunc(apiCfg.handlerCreateUserMediumRecord)))
	mux.Handle("GET /api/records", apiCfg.authMiddleware(http.HandlerFunc(apiCfg.handlerGetRecordsByUserID)))
	mux.Handle("PUT /api/records", apiCfg.authMiddleware(http.HandlerFunc(apiCfg.handlerUpdateRecord)))
	mux.Handle("DELETE /api/records", apiCfg.authMiddleware(http.HandlerFunc(apiCfg.handlerDeleteRecord)))

	// Authentification endpoints
	mux.HandleFunc("POST /auth/login", apiCfg.handlerLogin)
	mux.Handle("POST /auth/logout", apiCfg.authMiddleware(http.HandlerFunc(apiCfg.handlerLogout)))
	mux.HandleFunc("POST /auth/refresh", apiCfg.handlerRefresh)
	mux.HandleFunc("POST /auth/revoke", apiCfg.handlerRevoke)

	// Reset Password endpoints
	mux.HandleFunc("POST /auth/password-reset", apiCfg.handlerPasswordResetRequest)
	mux.HandleFunc("GET /auth/password-reset", apiCfg.handlerVerifyResetToken)
	mux.HandleFunc("PUT /auth/password-reset", apiCfg.handlerResetPassword)

	// Admin endpoint (only used on test server)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)
	mux.HandleFunc("GET /admin/user", apiCfg.handlerCheckUserExists)
	mux.HandleFunc("GET /admin/medium", apiCfg.handlerCheckMediumExists)
	mux.HandleFunc("GET /admin/record", apiCfg.handlerCheckRecordExists)

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
