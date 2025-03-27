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
	mux.HandleFunc("POST /api/users", apiCfg.handlerCreateUser)
	mux.Handle("PUT /api/users", apiCfg.authMiddleware(http.HandlerFunc(apiCfg.handlerUpdateUser)))

	mux.HandleFunc("POST /auth/refresh", apiCfg.handlerRefresh)
	mux.HandleFunc("POST /auth/revoke", apiCfg.handlerRevoke)

	mux.HandleFunc("POST /auth/login", apiCfg.handlerLogin)
	mux.HandleFunc("POST /auth/logout", apiCfg.handlerRevoke)

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
