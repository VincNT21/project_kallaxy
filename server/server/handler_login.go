package server

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/VincNT21/kallaxy/server/internal/auth"
	"github.com/VincNT21/kallaxy/server/internal/database"
	"github.com/jackc/pgx/v5/pgtype"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	type response struct {
		User
		Tokens
	}

	// Decode request body
	params := parameters{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 500, "couldn't decode body from request", err)
		return
	}

	// Get user info from database
	user, err := cfg.db.GetUserByUsername(r.Context(), params.Username)
	if err != nil {
		respondWithError(w, 500, "couldn't get user info from DB", err)
		return
	}

	// Check password validity
	err = auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		respondWithError(w, 401, "invalid password", err)
		return
	}

	// Create a JWT
	accessToken, err := auth.MakeJWT(user.ID, cfg.jwtsecret, time.Hour)
	if err != nil {
		respondWithError(w, 500, "couldn't create a JWT", err)
		return
	}

	// Create a Refresh token
	refreshTokenString, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, 500, "couldn't create a refresh token", err)
		return
	}

	// Store Refresh token in database
	refreshToken, err := cfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:  refreshTokenString,
		UserID: user.ID,
		ExpiresAt: pgtype.Timestamp{
			Time:  time.Now().UTC().AddDate(0, 0, 60),
			Valid: true,
		},
	})
	if err != nil {
		respondWithError(w, 500, "couldn't insert refresh token into db", err)
		return
	}

	// Respond
	respondWithJson(w, 201, response{
		User: User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Username:  user.Username,
			Email:     user.Email,
		},
		Tokens: Tokens{
			AccessToken:  accessToken,
			RefreshToken: refreshToken.Token,
		},
	})
}
