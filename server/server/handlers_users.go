package server

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/VincNT21/kallaxy/server/internal/auth"
	"github.com/VincNT21/kallaxy/server/internal/database"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

// POST /api/users
func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
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

	// Validate required fields
	if params.Username == "" {
		respondWithError(w, 400, "Username is required", errors.New("no username provided in request body"))
		return
	}
	if params.Password == "" {
		respondWithError(w, 400, "Password is required", errors.New("no password provided in request body"))
		return
	}
	if params.Email == "" {
		respondWithError(w, 400, "email is required", errors.New("no email provided in request body"))
		return
	}

	// Hash password
	hash, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, 500, "couldn't hash password", err)
		return
	}

	// Call query function
	user, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
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

// GET /api/users
func (cfg *apiConfig) handlerGetUserByID(w http.ResponseWriter, r *http.Request) {
	type response struct {
		User
	}

	// Get user's ID from context
	userID := r.Context().Value(userIDKey).(pgtype.UUID)

	// Call query function
	user, err := cfg.db.GetUserByID(r.Context(), userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Means that user doesn't exist which shouldn't happen if the JWT is valid
			respondWithError(w, 500, "server error: user record inconsistency", err)
			return
		}
		respondWithError(w, 500, "couldn't get user by ID in DB", err)
		return
	}

	// Respond
	respondWithJson(w, 200, response{
		User: User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Username:  user.Username,
			Email:     user.Email,
		},
	})

}

// PUT /api/users
func (cfg *apiConfig) handlerUpdateUser(w http.ResponseWriter, r *http.Request) {
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

	// Validate required fields
	if params.Username == "" {
		respondWithError(w, 400, "Username is required", errors.New("no username provided in request body"))
		return
	}
	if params.Password == "" {
		respondWithError(w, 400, "Password is required", errors.New("no password provided in request body"))
		return
	}
	if params.Email == "" {
		respondWithError(w, 400, "email is required", errors.New("no email provided in request body"))
		return
	}

	// Hash password
	hash, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, 500, "couldn't hash password", err)
		return
	}

	// Get user ID
	userID := r.Context().Value(userIDKey).(pgtype.UUID)

	// Call query function
	user, err := cfg.db.UpdateUser(r.Context(), database.UpdateUserParams{
		ID:             userID,
		Username:       params.Username,
		HashedPassword: hash,
		Email:          params.Email,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Means that user doesn't exist which shouldn't happen if the JWT is valid
			respondWithError(w, 500, "server error: user record inconsistency", err)
			return
		}
		respondWithError(w, 500, "couldn't update user info in DB", err)
		return
	}

	// Logout user by revoking all their refresh token
	count, err := cfg.db.RevokeAllRefreshTokensByUserID(r.Context(), user.ID)
	if err != nil {
		respondWithError(w, 500, "couldn't logout user", err)
		return
	}
	if count == 0 {
		respondWithError(w, 500, "error with logout user: no refresh token found", err)
		return
	}

	// Respond
	respondWithJson(w, 200, response{
		User: User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Username:  user.Username,
			Email:     user.Email,
		},
	})
}

// DELETE /api/users
func (cfg *apiConfig) handlerDeleteUser(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(userIDKey).(pgtype.UUID)
	count, err := cfg.db.DeleteUser(r.Context(), userID)
	if err != nil {
		respondWithError(w, 500, "couldn't delete user with given ID", err)
		return
	}
	if count == 0 {
		// Means that user doesn't exist which shouldn't happen if the JWT is valid
		respondWithError(w, 500, "server error: user record inconsistency", err)
		return
	}

	// Respond
	w.WriteHeader(200)
}
