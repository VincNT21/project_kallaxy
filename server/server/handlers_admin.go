package server

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/jackc/pgx/v5/pgtype"
)

// This handler is only used for integration tests
// No endpoint for it exists in production server
// POST /admin/reset
func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	// Reset users table
	err := cfg.db.ResetUsers(r.Context())
	if err != nil {
		respondWithError(w, 500, "couldn't reset table users", err)
		return
	}

	// Reset media table
	err = cfg.db.ResetMedia(r.Context())
	if err != nil {
		respondWithError(w, 500, "couldn't reset table media", err)
		return
	}

	// Reset records table
	err = cfg.db.ResetRecords(r.Context())
	if err != nil {
		respondWithError(w, 500, "couldn't reset table records", err)
		return
	}

	w.WriteHeader(200)
}

// This handler is only used for integration tests
// No endpoint for it exists in production server
// GET /admin/user
func (cfg *apiConfig) handlerCheckUserExists(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		UserID pgtype.UUID `json:"user_id"`
	}

	// Parse data from request body
	var params parameters
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		respondWithError(w, 500, "couldn't decode body from request", err)
		return
	}

	// Try to get User by ID
	_, err = cfg.db.GetUserByID(r.Context(), params.UserID)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		w.WriteHeader(404)
		return
	}
	w.WriteHeader(200)

}

// This handler is only used for integration tests
// No endpoint for it exists in production server
type parametersCheckMediumExists struct {
	MediumID pgtype.UUID `json:"medium_id"`
}

// GET /admin/medium
func (cfg *apiConfig) handlerCheckMediumExists(w http.ResponseWriter, r *http.Request) {

	// Parse data from request body
	var params parametersCheckMediumExists
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		respondWithError(w, 500, "couldn't decode body from request", err)
		return
	}

	// Try to get User by ID
	_, err = cfg.db.GetMediumByID(r.Context(), params.MediumID)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		w.WriteHeader(404)
		return
	}
	w.WriteHeader(200)

}

// This handler is only used for integration tests
// No endpoint for it exists in production server
type parametersCheckRecordExists struct {
	RecordID pgtype.UUID `json:"record_id"`
}

// GET /admin/medium
func (cfg *apiConfig) handlerCheckRecordExists(w http.ResponseWriter, r *http.Request) {

	// Parse data from request body
	var params parametersCheckRecordExists
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		respondWithError(w, 500, "couldn't decode body from request", err)
		return
	}

	// Try to get User by ID
	_, err = cfg.db.GetRecordByID(r.Context(), params.RecordID)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		w.WriteHeader(404)
		return
	}
	w.WriteHeader(200)

}
