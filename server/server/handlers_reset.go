package server

import "net/http"

// This handler is only used for integration tests
// No endpoint for it exists in production server
func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	// Reset table users
	err := cfg.db.ResetUsers(r.Context())
	if err != nil {
		respondWithError(w, 500, "couldn't reset table users", err)
		return
	}

	// Reset table media
	err = cfg.db.ResetMedia(r.Context())
	if err != nil {
		respondWithError(w, 500, "couldn't reset table media", err)
		return
	}
	w.WriteHeader(200)
}
