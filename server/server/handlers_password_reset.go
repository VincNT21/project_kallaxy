package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/VincNT21/kallaxy/server/internal/auth"
	"github.com/VincNT21/kallaxy/server/internal/database"
	"github.com/jackc/pgx/v5/pgtype"
)

type parametersPasswordResetRequest struct {
	Email string `json:"email"`
}

type responsePasswordResetRequest struct {
	Message    string `json:"message"`
	ResetLink  string `json:"reset_link"`
	ResetToken string `json:"reset_token"`
	Username   string `json:"username"`
}

// POST /auth/password_reset
func (cfg *apiConfig) handlerPasswordResetRequest(w http.ResponseWriter, r *http.Request) {

	// Parse email from request
	var params parametersPasswordResetRequest
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		respondWithError(w, 500, "couldn't decode body from request", err)
		return
	}

	// Find user by email
	user, err := cfg.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		// Client won't know if email exists, for security
		respondWithJson(w, 200, responsePasswordResetRequest{
			Message: "If your email exists in our system, you'll receive reset instruction",
		})
		return
	}

	// Generate reset token (valid for 6 hour)
	token := auth.GenerateResetToken()
	if token == "" {
		respondWithError(w, 500, "couldn't generate a reset token", errors.New("error with GenerateResetToken(): returning an empty string"))
		return
	}
	expiry := pgtype.Timestamp{
		Time:  time.Now().Add(6 * time.Hour),
		Valid: true,
	}

	// Store token in database
	cfg.db.StorePasswordToken(r.Context(), database.StorePasswordTokenParams{
		Token:     token,
		UserID:    user.ID,
		UserEmail: user.Email,
		ExpiresAt: expiry,
	})

	// For local development, just creating the path portion
	// instead of full URL
	resetLink := fmt.Sprintf("/auth/password_reset?token=%s", token)

	// In dev mode, return link in response
	respondWithJson(w, 200, responsePasswordResetRequest{
		Message:    "Password reset initiated",
		ResetLink:  resetLink,
		ResetToken: token,
		Username:   user.Username,
	})

	// In production, this would send an email instead
	// sendResestEmail(user.Email, resetLink)
	// respondWithJson(w, 200, response{ Message: "If your email exists in our system, you'll receive reset instruction"})
	// return

}

type responseVerifyResetToken struct {
	Valid bool   `json:"valid"`
	Email string `json:"email"`
}

// GET /auth/password_reset?token=xxxxxxx
func (cfg *apiConfig) handlerVerifyResetToken(w http.ResponseWriter, r *http.Request) {

	// Get token from URL query parameters
	token := r.URL.Query().Get("token")
	if token == "" {
		respondWithError(w, 400, "Missing reset token", nil)
		return
	}

	// Verify if token exists and is valid
	resetToken, err := cfg.db.GetPasswordResetToken(r.Context(), token)
	if err != nil || time.Now().After(resetToken.ExpiresAt.Time) || resetToken.UsedAt.Valid {
		respondWithError(w, 400, "Invalid or expired reset token", err)
		return
	}

	// If valid, respond
	respondWithJson(w, 200, responseVerifyResetToken{
		Valid: true,
		Email: resetToken.UserEmail,
	})

}

type parametersResetPassword struct {
	Token       string `json:"token"`
	NewPassword string `json:"new_password"`
}

// PUT /auth/password_reset
func (cfg *apiConfig) handlerResetPassword(w http.ResponseWriter, r *http.Request) {

	type response struct {
		User
	}

	// Parse data from request body
	var params parametersResetPassword
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		respondWithError(w, 500, "couldn't decode body from request", err)
		return
	}

	// Get Reset token from DB
	resetToken, err := cfg.db.GetPasswordResetToken(r.Context(), params.Token)
	if err != nil || time.Now().After(resetToken.ExpiresAt.Time) || resetToken.UsedAt.Valid {
		respondWithError(w, 400, "Invalid or expired reset token", err)
		return
	}

	// Validate password strength (to add)

	// Those steps would be better with a transaction system ()

	// Change password in db
	hash, err := auth.HashPassword(params.NewPassword)
	if err != nil {
		respondWithError(w, 500, "couldn't hash password", err)
		return
	}
	user, err := cfg.db.UpdatePassword(r.Context(), database.UpdatePasswordParams{
		ID:             resetToken.UserID,
		HashedPassword: hash,
	})
	if err != nil {
		respondWithError(w, 500, "couldn't update password in database", err)
		return
	}

	// Revoke all refresh tokens from user
	_, err = cfg.db.RevokeAllRefreshTokensByUserID(r.Context(), resetToken.UserID)
	if err != nil {
		respondWithError(w, 500, "couldn't revoke all user's refresh token", err)
		return
	}

	// Invalidate all Reset tokens from user
	err = cfg.db.InvalidateResetTokensByUserId(r.Context(), user.ID)
	if err != nil {
		respondWithError(w, 500, "couldn't invalidate all user's reset tokens", err)
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
