package server

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/VincNT21/kallaxy/server/internal/auth"
	"github.com/VincNT21/kallaxy/server/internal/database"
	"github.com/jackc/pgx/v5/pgtype"
)

// POST /auth/login
func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type response struct {
		User
		Tokens
	}

	// Decode request body
	params := parametersLogin{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 500, "couldn't decode body from request", err)
		return
	}

	// Get user info from database
	user, err := cfg.db.GetUserByUsername(r.Context(), params.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, 401, "username doesn't exist", err)
			return
		}
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

// POST /auth/logout
func (cfg *apiConfig) handlerLogout(w http.ResponseWriter, r *http.Request) {
	// Get user ID
	userID := r.Context().Value(userIDKey).(pgtype.UUID)

	// Call query function
	count, err := cfg.db.RevokeAllRefreshTokensByUserID(r.Context(), userID)
	if err != nil {
		respondWithError(w, 500, "couldn't revoke all refresh token by user ID", err)
		return
	}
	if count == 0 {
		respondWithError(w, 404, "no refresh token found for user ID", err)
		return
	}

	// Respond
	w.WriteHeader(204)
}

// POST /auth/refresh
func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Tokens
	}

	// Get the Refresh Token bearer from header
	refreshTokenOldStr, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "Missing of malformed refresh token in Authorization header", err)
		return
	}

	// Check if refresh token exists in db
	refreshTokenOld, err := cfg.db.GetRefreshToken(r.Context(), refreshTokenOldStr)
	if err != nil {
		respondWithError(w, 401, "Refresh token doesn't exist in db", err)
		return
	}

	// Check if refresh token is revoked
	if refreshTokenOld.RevokedAt.Valid {
		respondWithError(w, 401, "Refresh token has been revoked", errors.New("refresh token has been revoked"))
		return
	}

	// Check refresh token expiration
	if time.Now().UTC().Compare(refreshTokenOld.ExpiresAt.Time) >= 0 {
		respondWithError(w, 401, "Refresh token has expired", errors.New("refresh token has expired"))
		return
	}

	// Create a new JWT access token (validity: 1 hour)
	accessToken, err := auth.MakeJWT(refreshTokenOld.UserID, cfg.jwtsecret, time.Hour)
	if err != nil {
		respondWithError(w, 500, "couldn't create a new access token", err)
		return
	}

	// Create a new refresh token and store it in db
	refreshTokenNewStr, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, 500, "couldn't create a new refresh token", err)
		return
	}
	refreshTokenNew, err := cfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:  refreshTokenNewStr,
		UserID: refreshTokenOld.UserID,
		ExpiresAt: pgtype.Timestamp{
			Time:  time.Now().UTC().AddDate(0, 0, 60),
			Valid: true,
		},
	})
	if err != nil {
		respondWithError(w, 500, "couldn't insert refresh token into db", err)
	}

	// Revoked old refresh token
	_, err = cfg.db.RevokeRefreshToken(r.Context(), refreshTokenOldStr)
	if err != nil {
		respondWithError(w, 500, "couldn't delete old refresh token from db", err)
		return
	}

	// Respond with both tokens
	respondWithJson(w, 201, response{
		Tokens: Tokens{
			AccessToken:  accessToken,
			RefreshToken: refreshTokenNew.Token,
		},
	})
}

// POST /auth/revoke
func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	// Get the Refresh Token bearer from header
	refreshTokenStr, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "Missing of malformed refresh token in Authorization header", err)
		return
	}

	// Revoke Refresh token in DB
	count, err := cfg.db.RevokeRefreshToken(r.Context(), refreshTokenStr)
	if err != nil {
		respondWithError(w, 500, "couldn't revoke refresh token in db", err)
		return
	}
	if count == 0 {
		// Refresh token given is not in db
		respondWithError(w, 404, "refresh token not found", err)
		return
	}

	// If ok, respond
	w.WriteHeader(204)
}

// GET /auth/login
func (cfg *apiConfig) handlerConfirmPassword(w http.ResponseWriter, r *http.Request) {

	// Parse data from request body
	var params parametersConfirmPassword
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		respondWithError(w, 500, "couldn't decode body from request", err)
		return
	}

	// Get user ID
	userID := r.Context().Value(userIDKey).(pgtype.UUID)

	// Get user Info
	user, err := cfg.db.GetUserByID(r.Context(), userID)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		w.WriteHeader(500)
		return
	}

	// Check password validity
	err = auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		respondWithError(w, 401, "invalid password", err)
		return
	}

	w.WriteHeader(200)

}
