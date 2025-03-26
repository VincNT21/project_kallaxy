package server

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/VincNT21/kallaxy/server/internal/auth"
	"github.com/VincNT21/kallaxy/server/internal/database"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

type User struct {
	ID        pgtype.UUID      `json:"id"`
	CreatedAt pgtype.Timestamp `json:"created_at"`
	UpdatedAt pgtype.Timestamp `json:"updated_at"`
	Username  string           `json:"username"`
	Email     string           `json:"email"`
}

func (cfg *apiConfig) CreateUser(w http.ResponseWriter, r *http.Request) {
	// Parameters struct match what we'll get from request
	type parameters struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	// Response struct match what we'll use for response
	type response struct {
		User
	}

	// Get body from request
	decoder := json.NewDecoder(r.Body)
	params := parameters{}

	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 500, "couldn't decode body from request", err)
		return
	}

	// Hash password
	hash, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, 500, "couldn't hash password", err)
		return
	}

	// Call query function
	user, err := cfg.db.CreateUser(context.Background(), database.CreateUserParams{
		Username:       params.Username,
		HashedPassword: hash,
		Email:          params.Email,
	})
	if err != nil {
		// Check if error comes from a database unique constraint violation
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			respondWithError(w, 409, "username or email is already used by another user", err)
		} else {
			respondWithError(w, 500, "couldn't create user in DB", err)
		}
		return
	}

	// Log
	log.Printf("New user '%s' created in DB", user.Username)

	// Respond with new user's data
	respondWithJson(w, 201, response{
		User: User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Username:  user.Username,
			Email:     user.Email,
		},
	})
}
