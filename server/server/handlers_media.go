package server

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/VincNT21/kallaxy/server/internal/database"
	"github.com/jackc/pgx/v5/pgconn"
)

type parametersCreateMedium struct {
	Title     string                 `json:"title"`
	MediaType string                 `json:"media_type"`
	Creator   string                 `json:"creator"`
	PubDate   string                 `json:"pub_date"`
	ImageUrl  string                 `json:"image_url"`
	Metadata  map[string]interface{} `json:"metadata"`
}

// POST /api/media
func (cfg *apiConfig) handlerCreateMedium(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Medium
	}

	// Parse data from request body
	var params parametersCreateMedium
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		respondWithError(w, 500, "couldn't decode body from request", err)
		return
	}

	// Check if all required fields are provided
	if params.Title == "" || params.MediaType == "" || params.Creator == "" || params.PubDate == "" {
		respondWithError(w, 400, "Some required field is missing in request body", errors.New("a imperative field is missing in request's body"))
		return
	}

	// Convert metadata map to []byte
	metadataBytes, err := mapToBytes(params.Metadata)
	if err != nil {
		respondWithError(w, 500, "couldn't convert metadata map for database", err)
		return
	}

	// Call query function
	medium, err := cfg.db.CreateMedium(r.Context(), database.CreateMediumParams{
		MediaType: params.MediaType,
		Title:     params.Title,
		Creator:   params.Creator,
		PubDate:   params.PubDate,
		ImageUrl:  params.ImageUrl,
		Metadata:  metadataBytes,
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			// This is a unique constraint violation
			respondWithError(w, 409, "A medium with same title already exists in database", err)
			return
		}
		respondWithError(w, 500, "couldn't create new medium in database", err)
		return
	}

	// Convert metadata back to map
	metadataMap, err := bytesToMap(medium.Metadata)
	if err != nil {
		respondWithError(w, 500, "couldn't convert metadata map from database", err)
		return
	}

	// Respond
	respondWithJson(w, 201, response{
		Medium: Medium{
			ID:        medium.ID,
			MediaType: medium.MediaType,
			CreatedAt: medium.CreatedAt,
			UpdatedAt: medium.UpdatedAt,
			Title:     medium.Title,
			Creator:   medium.Creator,
			PubDate:   medium.PubDate,
			ImageUrl:  medium.ImageUrl,
			Metadata:  metadataMap,
		},
	})
}

type parametersGetMediumByTitleAndType struct {
	Title     string `json:"title"`
	MediaType string `json:"media_type"`
}

// GET /api/media
func (cfg *apiConfig) handlerGetMediumByTitleAndType(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Medium
	}

	// Parse data from request body
	var params parametersGetMediumByTitleAndType
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		respondWithError(w, 500, "couldn't decode body from request", err)
		return
	}

	// Call query function
	medium, err := cfg.db.GetMediumByTitleAndType(r.Context(), database.GetMediumByTitleAndTypeParams{
		Lower:   params.Title,
		Lower_2: params.MediaType,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, 404, fmt.Sprintf("No %s medium with title %s in database", params.MediaType, params.Title), err)
			return
		}
		respondWithError(w, 500, "couldn't get medium by title", err)
		return
	}

	// Convert metadata back to map
	metadataMap, err := bytesToMap(medium.Metadata)
	if err != nil {
		respondWithError(w, 500, "couldn't convert metadata map from database", err)
		return
	}

	// Respond
	respondWithJson(w, 200, response{
		Medium: Medium{
			ID:        medium.ID,
			MediaType: medium.MediaType,
			CreatedAt: medium.CreatedAt,
			UpdatedAt: medium.UpdatedAt,
			Title:     medium.Title,
			Creator:   medium.Creator,
			PubDate:   medium.PubDate,
			ImageUrl:  medium.ImageUrl,
			Metadata:  metadataMap,
		},
	})

}

type parametersGetMediaByType struct {
	MediaType string `json:"media_type"`
}

type responseGetMediaByType struct {
	Media []Medium `json:"media"`
}

// GET /api/media/type
func (cfg *apiConfig) handlerGetMediaByType(w http.ResponseWriter, r *http.Request) {

	// Parse data from request body
	var params parametersGetMediaByType
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		respondWithError(w, 500, "couldn't decode body from request", err)
		return
	}

	// Call query function
	media, err := cfg.db.GetMediaByType(r.Context(), params.MediaType)
	if err != nil {
		respondWithError(w, 500, "couldn't get media by given type", err)
		return
	}
	if len(media) == 0 {
		respondWithError(w, 404, "No media of given type in database", err)
		return
	}

	response := responseGetMediaByType{}
	for _, medium := range media {
		// Convert metadata back to map
		metadataMap, err := bytesToMap(medium.Metadata)
		if err != nil {
			respondWithError(w, 500, "couldn't convert metadata map from database", err)
			return
		}

		response.Media = append(response.Media, Medium{
			ID:        medium.ID,
			MediaType: medium.MediaType,
			CreatedAt: medium.CreatedAt,
			UpdatedAt: medium.UpdatedAt,
			Title:     medium.Title,
			Creator:   medium.Creator,
			PubDate:   medium.PubDate,
			ImageUrl:  medium.ImageUrl,
			Metadata:  metadataMap,
		})
	}

	// Respond
	respondWithJson(w, 200, response)

}

type parametersUpdateMedium struct {
	MediumID string                 `json:"medium_id"`
	Title    string                 `json:"title"`
	Creator  string                 `json:"creator"`
	PubDate  string                 `json:"pub_date"`
	ImageUrl string                 `json:"image_url"`
	Metadata map[string]interface{} `json:"metadata"`
}

// PUT /api/media
func (cfg *apiConfig) handlerUpdateMedium(w http.ResponseWriter, r *http.Request) {

	// Parse data from request body
	var params parametersUpdateMedium
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		respondWithError(w, 500, "couldn't decode body from request", err)
		return
	}

	// Convert MediumID to pgtype.UUID
	mediumID, err := convertIdToPgtype(params.MediumID)
	if err != nil {
		respondWithError(w, 400, "medium_id not in good format", err)
		return
	}

	// Convert metadata map to []byte
	metadataBytes, err := mapToBytes(params.Metadata)
	if err != nil {
		respondWithError(w, 500, "couldn't convert metadata map for database", err)
		return
	}

	// Call query function
	medium, err := cfg.db.UpdateMedium(r.Context(), database.UpdateMediumParams{
		ID:       mediumID,
		Title:    params.Title,
		Creator:  params.Creator,
		PubDate:  params.PubDate,
		ImageUrl: params.ImageUrl,
		Metadata: metadataBytes,
	})

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, 404, "No medium with given ID in database", err)
			return
		}
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			// This is a unique constraint violation
			respondWithError(w, 409, "A medium with same title already exists in database", err)
			return
		}
		respondWithError(w, 500, "couldn't udpate medium by given ID", err)
		return
	}

	// Convert metadata back to map
	metadataMap, err := bytesToMap(medium.Metadata)
	if err != nil {
		respondWithError(w, 500, "couldn't convert metadata map from database", err)
		return
	}

	// Respond
	respondWithJson(w, 200, Medium{
		ID:        medium.ID,
		MediaType: medium.MediaType,
		CreatedAt: medium.CreatedAt,
		UpdatedAt: medium.UpdatedAt,
		Title:     medium.Title,
		Creator:   medium.Creator,
		PubDate:   medium.PubDate,
		ImageUrl:  medium.ImageUrl,
		Metadata:  metadataMap,
	})
}

type parametersDeleteMedium struct {
	MediumID string `json:"medium_id"`
}

// DELETE /api/media
func (cfg *apiConfig) handlerDeleteMedium(w http.ResponseWriter, r *http.Request) {

	// Parse data from request body
	var params parametersDeleteMedium
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		respondWithError(w, 500, "couldn't decode body from request", err)
		return
	}

	// Convert MediumID to pgtype.UUID
	mediumID, err := convertIdToPgtype(params.MediumID)
	if err != nil {
		respondWithError(w, 400, "medium_id not in good format", err)
		return
	}

	// Call query function
	count, err := cfg.db.DeleteMedium(r.Context(), mediumID)
	if err != nil {
		respondWithError(w, 500, "couldn't delete medium with id on database", err)
		return
	}
	if count == 0 {
		respondWithError(w, 404, "No medium with given ID in database", nil)
		return
	}

	// Respond
	w.WriteHeader(200)

}
