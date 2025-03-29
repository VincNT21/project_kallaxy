package server

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/VincNT21/kallaxy/server/internal/database"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

type parametersCreateUserMediumRecord struct {
	UserID     pgtype.UUID      `json:"user_id"`
	MediaID    pgtype.UUID      `json:"media_id"`
	IsFinished pgtype.Bool      `json:"is_finished"`
	StartDate  pgtype.Timestamp `json:"start_date"`
	EndDate    pgtype.Timestamp `json:"end_date"`
}

// POST /api/records
func (cfg *apiConfig) handlerCreateUserMediumRecord(w http.ResponseWriter, r *http.Request) {

	// Parse data from request body
	var params parametersCreateUserMediumRecord
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		respondWithError(w, 500, "couldn't decode body from request", err)
		return
	}

	// Check if all required fields are provide
	emptyID := pgtype.UUID{Valid: false}
	if params.UserID == emptyID || params.MediaID == emptyID {
		respondWithError(w, 400, "Some required field is missing in request body", errors.New("a imperative field is missing in request's body"))
		return
	}

	// Calculate interval duration
	interval, err := calculateDuration(params.StartDate, params.EndDate)
	if err != nil {
		respondWithError(w, 500, "couldn't calculate interval between start date and end date", err)
		return
	}

	// Call query function
	record, err := cfg.db.CreateUserMediumRecord(r.Context(), database.CreateUserMediumRecordParams{
		UserID:     params.UserID,
		MediaID:    params.MediaID,
		IsFinished: params.IsFinished,
		StartDate:  params.StartDate,
		EndDate:    params.EndDate,
		Duration:   interval,
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23503" {
			// This is a foreign key constraint violation
			respondWithError(w, 404, "given user id or media id didn't exist in database", err)
			return
		} else if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			// This is a unique constraint violation (about couple user/media id)
			respondWithError(w, 409, "there is already a record in database with same couple user/media id", err)
			return
		}
		respondWithError(w, 500, "couldn't create new record in database", err)
		return
	}

	// Respond
	respondWithJson(w, 201, Record{
		ID:         record.ID,
		CreatedAt:  record.CreatedAt,
		UpdatedAt:  record.UpdatedAt,
		UserID:     record.UserID,
		MediaID:    record.MediaID,
		IsFinished: record.IsFinished,
		StartDate:  record.StartDate,
		EndDate:    record.EndDate,
		Duration:   record.Duration,
	})
}

type responseGetRecordsByUserID struct {
	Records []Record `json:"records"`
}

// GET /api/records
func (cfg *apiConfig) handlerGetRecordsByUserID(w http.ResponseWriter, r *http.Request) {

	// Get userID from access token
	userID := r.Context().Value(userIDKey).(pgtype.UUID)

	// Call query function
	records, err := cfg.db.GetRecordsByUserID(r.Context(), userID)
	if err != nil {
		respondWithError(w, 500, "couldn't get records by user ID in database", err)
		return
	}
	if len(records) == 0 {
		respondWithError(w, 404, "no record found for giver user id", errors.New("records list returning from query was empty"))
		return
	}

	response := responseGetRecordsByUserID{}
	for _, record := range records {
		response.Records = append(response.Records, Record{
			ID:         record.ID,
			CreatedAt:  record.CreatedAt,
			UpdatedAt:  record.UpdatedAt,
			UserID:     record.UserID,
			MediaID:    record.MediaID,
			IsFinished: record.IsFinished,
			StartDate:  record.StartDate,
			EndDate:    record.EndDate,
			Duration:   record.Duration,
		})
	}

	// Respond
	respondWithJson(w, 200, response)
}

type parametersUpdateRecord struct {
	RecordID   pgtype.UUID      `json:"record_id"`
	IsFinished pgtype.Bool      `json:"is_finished"`
	StartDate  pgtype.Timestamp `json:"start_date"`
	EndDate    pgtype.Timestamp `json:"end_date"`
}

// PUT /api/records
func (cfg *apiConfig) handlerUpdateRecord(w http.ResponseWriter, r *http.Request) {

	// Parse data from request body
	var params parametersUpdateRecord
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		respondWithError(w, 500, "couldn't decode body from request", err)
		return
	}

	// Calculate interval duration
	interval, err := calculateDuration(params.StartDate, params.EndDate)
	if err != nil {
		respondWithError(w, 500, "couldn't calculate interval between start date and end date", err)
		return
	}

	// Call query function
	record, err := cfg.db.UpdateRecord(r.Context(), database.UpdateRecordParams{
		ID:         params.RecordID,
		IsFinished: params.IsFinished,
		StartDate:  params.StartDate,
		EndDate:    params.EndDate,
		Duration:   interval,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, 404, "No record found with given ID", err)
			return
		}
		respondWithError(w, 500, "couldn't update record in database", err)
		return
	}

	// Respond
	respondWithJson(w, 200, Record{
		ID:         record.ID,
		CreatedAt:  record.CreatedAt,
		UpdatedAt:  record.UpdatedAt,
		UserID:     record.UserID,
		MediaID:    record.MediaID,
		IsFinished: record.IsFinished,
		StartDate:  record.StartDate,
		EndDate:    record.EndDate,
		Duration:   record.Duration,
	})
}

type parametersDeleteRecord struct {
	RecordID pgtype.UUID `json:"record_id"`
}

// DELETE /api/records
func (cfg *apiConfig) handlerDeleteRecord(w http.ResponseWriter, r *http.Request) {

	type response struct {
	}

	// Parse data from request body
	var params parametersDeleteRecord
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		respondWithError(w, 500, "couldn't decode body from request", err)
		return
	}

	// Call query function
	count, err := cfg.db.DeleteRecord(r.Context(), params.RecordID)
	if err != nil {
		respondWithError(w, 500, "couldn't delete record in database", err)
		return
	}
	if count == 0 {
		respondWithError(w, 404, "No record with given ID in database", nil)
		return
	}

	// Respond
	w.WriteHeader(200)

}
