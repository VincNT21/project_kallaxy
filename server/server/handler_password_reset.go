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

func (cfg *apiConfig) handlerPasswordResetRequest(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string `json:"email"`
	}

	type response struct {
		Message   string `json:"message"`
		ResetLink string `json:"reset_link"`
		Token     string `json:"token"`
	}

	// Parse email from request
	var params parameters
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		respondWithError(w, 500, "couldn't decode body from request", err)
		return
	}

	// Find user by email
	user, err := cfg.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		// Client won't know if email exists, for security
		respondWithJson(w, 200, response{
			Message: "If your email exists in our system, you'll receive reset instruction",
		})
	}

	// Generate reset token (valid for 6 hour)
	token := auth.GenerateResetToken()
	if token == "" {
		respondWithError(w, 500, "couldn't generate a reset token", errors.New("Error with GenerateResetToken(): returning an empty string"))
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
		ExpiresAt: expiry,
	})

	// For loca development, just creating the path portion
	// instead of full URL
	resetLink := fmt.Sprintf("/reset-password?token=%s", token)

	// In dev mode, return link in response
	respondWithJson(w, 200, response{
		Message:   "Password reset initiated",
		ResetLink: resetLink,
		Token:     token,
	})
	return

	// In production, this would send an email instead
	// sendResestEmail(user.Email, resetLink)
	// respondWithJson(w, 200, response{ Message: "If your email exists in our system, you'll receive reset instruction"})
	// return

}
