package server

import "net/http"

// This handler is only used for integration tests
// No endpoint for it exists in production server
func (cfg *apiConfig) handlerResetUsers(w http.ResponseWriter, r *http.Request) {
	err := cfg.db.ResetUsers(r.Context())
	if err != nil {
		respondWithError(w, 500, "couldn't reset table users", err)
		return
	}
	w.WriteHeader(200)
}

// This handler is only used for integration tests
// No endpoint for it exists in production server
func (cfg *apiConfig) handlerResetMedia(w http.ResponseWriter, r *http.Request) {
	err := cfg.db.ResetMedia(r.Context())
	if err != nil {
		respondWithError(w, 500, "couldn't reset table users", err)
		return
	}
	w.WriteHeader(200)
}
