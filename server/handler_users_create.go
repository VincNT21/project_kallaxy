package server

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/VincNT21/kallaxy/server/internal/database"
)

func (cfg *apiConfig) CreateUser(w http.ResponseWriter, r *http.Request) {
	// Parameters struct match what we'll get from request
	type parameters struct {
		username string `json:"username"`
		password string `json:"password"`
		email    string `json:"email"`
	}
	// Response struct match what we'll use for response
	type response struct {
		username string `json:"username"`
		email    string `json:"username"`
	}

	// Get body from request
	decoder := json.NewDecoder(r.Body)
	params := parameters{}

	err := decoder.Decode(&params)
	if err != nil {
		log.Print("Error decoding body from request")
	}

	// Call query function
	user, err := cfg.db.CreateUser(context.Background(), database.CreateUserParams{})
	if err != nil {
		log.Print("Error creating user")
	}

}
