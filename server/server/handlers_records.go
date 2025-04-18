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

// POST /api/records
func (cfg *apiConfig) handlerCreateUserMediumRecord(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Record
	}

	// Parse data from request body
	var params parametersCreateUserMediumRecord
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		respondWithError(w, 500, "couldn't decode body from request", err)
		return
	}

	// Check if all required fields are provided
	if params.MediumID == "" {
		respondWithError(w, 400, "Some required field is missing in request body", errors.New("a imperative field is missing in request's body"))
		return
	}

	// Convert MediumID to pgtype.UUID
	mediumID, err := convertIdToPgtype(params.MediumID)
	if err != nil {
		respondWithError(w, 400, "medium_id not in good format", err)
		return
	}

	// Convert dates to pgtype.Timestamp
	startDate, err := convertDateToPgtype(params.StartDate)
	if err != nil {
		respondWithError(w, 400, "start_date not in good format", err)
		return
	}
	endDate, err := convertDateToPgtype(params.EndDate)
	if err != nil {
		respondWithError(w, 400, "end_date not in good format", err)
		return
	}

	// Get user ID
	userID := r.Context().Value(userIDKey).(pgtype.UUID)

	// Set is_finished
	isFinished := pgtype.Bool{Valid: true}
	if startDate.Valid && endDate.Valid {
		isFinished.Bool = true
	} else {
		isFinished.Bool = false
	}

	// Calculate interval duration
	interval, err := calculateDuration(startDate, endDate)
	if err != nil {
		respondWithError(w, 400, "start date is after end date", err)
		return
	}

	// Call query function
	record, err := cfg.db.CreateUserMediumRecord(r.Context(), database.CreateUserMediumRecordParams{
		UserID:     userID,
		MediaID:    mediumID,
		IsFinished: isFinished,
		StartDate:  startDate,
		EndDate:    endDate,
		Duration:   interval,
		Comments:   params.Comments,
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23503" {
			// This is a foreign key constraint violation
			respondWithError(w, 404, "given user id or media id didn't exist in database", err)
			return
		} else if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			// This is a unique constraint violation (about couple user/media id)
			respondWithError(w, 409, "there is already a record in database with same user-medium couple", err)
			return
		}
		respondWithError(w, 500, "couldn't create new record in database", err)
		return
	}

	// Respond
	respondWithJson(w, 201, response{
		Record: Record{
			ID:         record.ID,
			CreatedAt:  record.CreatedAt,
			UpdatedAt:  record.UpdatedAt,
			UserID:     record.UserID,
			MediaID:    record.MediaID,
			IsFinished: record.IsFinished,
			StartDate:  record.StartDate,
			EndDate:    record.EndDate,
			Duration:   record.Duration.Days,
			Comments:   record.Comments,
		},
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
		respondWithError(w, 404, "no record found for given user id", errors.New("records list returning from query was empty"))
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
			Duration:   record.Duration.Days,
			Comments:   record.Comments,
		})
	}

	// Respond
	respondWithJson(w, 200, response)
}

// PUT /api/records
func (cfg *apiConfig) handlerUpdateRecord(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Record
	}

	// Parse data from request body
	var params parametersUpdateRecord
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		respondWithError(w, 500, "couldn't decode body from request", err)
		return
	}

	// Convert RecordID to pgtype.UUID
	recordID, err := convertIdToPgtype(params.RecordID)
	if err != nil {
		respondWithError(w, 400, "record_id not in good format", err)
		return
	}

	// Convert dates to pgtype.Timestamp
	paramStartDate, err := convertDateToPgtype(params.StartDate)
	if err != nil {
		respondWithError(w, 400, "start_date not in good format", err)
		return
	}
	paramEndDate, err := convertDateToPgtype(params.EndDate)
	if err != nil {
		respondWithError(w, 400, "end_date not in good format", err)
		return
	}

	// Get already previous info from record in database
	previousDates, err := cfg.db.GetDatesFromRecord(r.Context(), recordID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, 404, "No record found with given ID", err)
			return
		}
		respondWithError(w, 500, "couldn't get record's dates from db", err)
		return
	}

	// Check if dates has been modified
	startDate := pgtype.Timestamp{}
	if !paramStartDate.Valid || paramStartDate == previousDates.StartDate {
		startDate = previousDates.StartDate
	} else {
		startDate = paramStartDate
	}
	endDate := pgtype.Timestamp{}
	if !paramEndDate.Valid || paramEndDate == previousDates.EndDate {
		endDate = previousDates.EndDate
	} else {
		endDate = paramEndDate
	}

	// Set is_finished
	isFinished := pgtype.Bool{Valid: true}
	if startDate.Valid && endDate.Valid {
		isFinished.Bool = true
	} else {
		isFinished.Bool = false
	}

	// Calculate interval duration
	interval, err := calculateDuration(startDate, endDate)
	if err != nil {
		respondWithError(w, 400, "start date is after end date", err)
		return
	}

	// Call query function
	record, err := cfg.db.UpdateRecord(r.Context(), database.UpdateRecordParams{
		ID:         recordID,
		IsFinished: isFinished,
		StartDate:  startDate,
		EndDate:    endDate,
		Duration:   interval,
		Comments:   params.Comments,
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
	respondWithJson(w, 200, response{
		Record: Record{
			ID:         record.ID,
			CreatedAt:  record.CreatedAt,
			UpdatedAt:  record.UpdatedAt,
			UserID:     record.UserID,
			MediaID:    record.MediaID,
			IsFinished: record.IsFinished,
			StartDate:  record.StartDate,
			EndDate:    record.EndDate,
			Duration:   record.Duration.Days,
			Comments:   record.Comments,
		},
	})
}

type parametersDeleteRecord struct {
	MediumID string `json:"medium_id"`
}

// DELETE /api/records
func (cfg *apiConfig) handlerDeleteRecord(w http.ResponseWriter, r *http.Request) {

	// Get userID from access token
	userID := r.Context().Value(userIDKey).(pgtype.UUID)

	// Parse data from request body
	var params parametersDeleteRecord
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		respondWithError(w, 500, "couldn't decode body from request", err)
		return
	}

	// Convert MediumID to pgtype.UUID
	mediumID, err := convertIdToPgtype(params.MediumID)
	if err != nil {
		respondWithError(w, 400, "record_id not in good format", err)
		return
	}

	// Call query function
	count, err := cfg.db.DeleteRecord(r.Context(), database.DeleteRecordParams{
		MediaID: mediumID,
		UserID:  userID,
	})
	if err != nil {
		respondWithError(w, 500, "couldn't delete record in database", err)
		return
	}
	if count == 0 {
		respondWithError(w, 404, "No record with given IDs in database", nil)
		return
	}

	// Respond
	w.WriteHeader(200)

}
