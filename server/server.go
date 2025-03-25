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

func Start() {
	const port = "8080"
	// Load the .env file into the environement
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("couldn't load .env file: %v", err)
	}

	// Get info from .env
	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		log.Fatal("DB_URL must be set")
	}

	// Open a connection to database
	dbConnection, err := pgxpool.New(context.Background(), dbUrl)
	if err != nil {
		log.Fatalf("Couldn't open a connection to db: %v", err)
	}
	defer dbConnection.Close()

	// Create a *database.Queries to store in config struct
	db := database.New(dbConnection)

	// Init apiCfg
	apiCfg := newAPIConfig(db)

	// Create the request multiplexer (router)
	mux := http.NewServeMux()

	// Register handlers
	apiCfg.db.CreateUser(context.Background(), database.CreateUserParams{})

	// Create a http server that listens on defined port and use multiplexer
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	// Start the server and log a fatal error if it fails to start
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	} else {
		log.Printf("Server started on port %s\n", port)
	}

}
