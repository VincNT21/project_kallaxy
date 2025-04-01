package server

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/VincNT21/kallaxy/server/internal/database"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func Start(envVars ...map[string]string) {
	const port = "8080"

	// Env variables are provided for testing cases
	// Check for them
	if len(envVars) > 0 && envVars[0] != nil {
		// Testing case: Use provided environment variables
		for key, value := range envVars[0] {
			os.Setenv(key, value)
		}
	} else {
		// Normal case: Load the .env file into the environement
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatalf("--FATAL ERROR-- couldn't load .env file: %v", err)
		}
	}

	// Get info from .env
	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		log.Fatal("--FATAL ERROR-- DB_URL env. variable must be set")
	}
	jwtsecret := os.Getenv("SECRET")
	if jwtsecret == "" {
		log.Fatal("--FATAL ERROR-- SECRET env. variable must be set")
	}

	// Open a connection to database
	dbConnection, err := pgxpool.New(context.Background(), dbUrl)
	if err != nil {
		log.Fatalf("--FATAL ERROR-- Couldn't open a connection to db: %v", err)
	}
	defer dbConnection.Close()

	// Create a *database.Queries to store in config struct
	db := database.New(dbConnection)

	// Init apiCfg
	apiCfg := newAPIConfig(db, jwtsecret)

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
	mux.HandleFunc("POST /auth/password_reset", apiCfg.handlerPasswordResetRequest)
	mux.HandleFunc("GET /auth/password_reset", apiCfg.handlerVerifyResetToken)
	mux.HandleFunc("PUT /auth/password_reset", apiCfg.handlerResetPassword)

	// Admin endpoint (only used on test server)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)
	mux.HandleFunc("GET /admin/user", apiCfg.handlerCheckUserExists)
	mux.HandleFunc("GET /admin/medium", apiCfg.handlerCheckMediumExists)
	mux.HandleFunc("GET /admin/record", apiCfg.handlerCheckRecordExists)

	// Create a http server that listens on defined port and use multiplexer
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("--DEBUG-- Server started on port %s\n", port)

	// Start the server and log a fatal error if it fails to start
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}

// Delete any revoked Refresh token
func (cfg *apiConfig) CleanRefreshTokens() {
	err := cfg.db.DeleteRevokedTokens(context.Background())
	if err != nil {
		log.Printf("--ERROR-- Couldn't delete revoked refresh tokens in db: %v", err)
		return
	}
	log.Println("--INFO-- Deleting revoked refresh token successful")
}
