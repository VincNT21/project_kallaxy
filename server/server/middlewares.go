package server

import (
	"context"
	"net/http"

	"github.com/VincNT21/kallaxy/server/internal/auth"
)

type contextKey string

const userIDKey contextKey = "user_id"

func (cfg *apiConfig) authMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get access token
		tokenString, err := auth.GetBearerToken(r.Header)
		if err != nil {
			respondWithError(w, 401, "Missing or malformed access token in Authorization header", err)
			return
		}

		// Validate JWT and get user ID
		userID, err := auth.ValidateJWT(tokenString, cfg.jwtsecret)
		if err != nil {
			respondWithError(w, 401, "Invalid or expired access token", err)
			return
		}

		// Add user ID to request context
		ctx := context.WithValue(r.Context(), userIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
