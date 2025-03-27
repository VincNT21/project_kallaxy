package server

import (
	"errors"
	"net/http"
	"time"

	"github.com/VincNT21/kallaxy/server/internal/auth"
	"github.com/VincNT21/kallaxy/server/internal/database"
	"github.com/jackc/pgx/v5/pgtype"
)

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
