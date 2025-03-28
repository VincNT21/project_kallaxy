package server

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/VincNT21/kallaxy/server/internal/database"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

// POST /api/media
func (cfg *apiConfig) handlerCreateMedium(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Title       string          `json:"title"`
		MediaType   string          `json:"media_type"`
		Creator     string          `json:"creator"`
		ReleaseYear int32           `json:"release_year"`
		ImageUrl    string          `json:"image_url"`
		Metadata    json.RawMessage `json:"metadata"`
	}

	type response struct {
		Medium
	}

	// Parse data from request body
	var params parameters
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		respondWithError(w, 500, "couldn't decode body from request", err)
		return
	}

	// Call query function
	medium, err := cfg.db.CreateMedium(r.Context(), database.CreateMediumParams{
		MediaType:   params.MediaType,
		Title:       params.Title,
		Creator:     params.Creator,
		ReleaseYear: params.ReleaseYear,
		ImageUrl:    params.ImageUrl,
		Metadata:    params.Metadata,
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			// This is a unique constraint violation
			respondWithError(w, 400, "A medium with same title already exists in database", err)
			return
		}
		respondWithError(w, 500, "couldn't create new medium in database", err)
		return
	}

	// Respond
	respondWithJson(w, 201, response{
		Medium: Medium{
			ID:          medium.ID,
			MediaType:   medium.MediaType,
			CreatedAt:   medium.CreatedAt,
			UpdatedAt:   medium.UpdatedAt,
			Title:       medium.Title,
			Creator:     medium.Creator,
			ReleaseYear: medium.ReleaseYear,
			ImageUrl:    medium.ImageUrl,
			Metadata:    medium.Metadata,
		}})
}

// GET /api/media (query parameters: "?title=xxx" / "?type=xxx"
func (cfg *apiConfig) handlerGetMedia(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Title     string
		MediaType string
	}

	type response struct {
		Media []Medium `json:"media"`
	}

	// Get parameters from request query parameters
	p := parameters{}
	p.Title = r.URL.Query().Get("title")
	p.MediaType = r.URL.Query().Get("type")

	// Handle different cases
	if p.Title != "" {
		cfg.getMediumByTitle(w, r, p.Title)
		return
	} else if p.MediaType != "" {
		cfg.getMediaByType(w, r, p.MediaType)
		return
	} else {
		// We could get all media here, but for now, it respond with error
		respondWithError(w, 400, "must provide either 'title' or 'type' query parameter", nil)
		return
	}

}

// Sub-function for handlerGetMedia
func (cfg *apiConfig) getMediumByTitle(w http.ResponseWriter, r *http.Request, title string) {
	type response struct {
		Medium
	}

	// Call query function
	medium, err := cfg.db.GetMediumByTitle(r.Context(), title)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, 404, fmt.Sprintf("No medium with title %s in database", title), err)
			return
		}
		respondWithError(w, 500, "couldn't get medium by title", err)
		return
	}

	// Respond
	respondWithJson(w, 200, response{
		Medium: Medium{
			ID:          medium.ID,
			MediaType:   medium.MediaType,
			CreatedAt:   medium.CreatedAt,
			UpdatedAt:   medium.UpdatedAt,
			Title:       medium.Title,
			Creator:     medium.Creator,
			ReleaseYear: medium.ReleaseYear,
			ImageUrl:    medium.ImageUrl,
			Metadata:    medium.Metadata,
		}})

}

// Sub-function for handlerGetMedia
func (cfg *apiConfig) getMediaByType(w http.ResponseWriter, r *http.Request, mediaType string) {
	type response struct {
		Media []Medium `json:"media"`
	}

	// Call query function
	media, err := cfg.db.GetMediaByType(r.Context(), mediaType)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, 404, "No media of given type in database", err)
			return
		}
		respondWithError(w, 500, "couldn't get media by given type", err)
		return
	}

	responseMedia := []Medium{}
	for _, medium := range media {
		responseMedia = append(responseMedia, Medium{
			ID:          medium.ID,
			MediaType:   medium.MediaType,
			CreatedAt:   medium.CreatedAt,
			UpdatedAt:   medium.UpdatedAt,
			Title:       medium.Title,
			Creator:     medium.Creator,
			ReleaseYear: medium.ReleaseYear,
			ImageUrl:    medium.ImageUrl,
			Metadata:    medium.Metadata,
		})
	}

	// Respond
	respondWithJson(w, 200, response{
		Media: responseMedia,
	})

}

// PUT /api/media
func (cfg *apiConfig) handlerUpdateMedium(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		ID          pgtype.UUID     `json:"id"`
		Title       string          `json:"title"`
		MediaType   string          `json:"media_type"`
		Creator     string          `json:"creator"`
		ReleaseYear int32           `json:"release_year"`
		ImageUrl    string          `json:"image_url"`
		Metadata    json.RawMessage `json:"metadata"`
	}

	type response struct {
		Medium
	}

	// Parse data from request body
	var params parameters
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		respondWithError(w, 500, "couldn't decode body from request", err)
		return
	}

	// Call query function
	medium, err := cfg.db.UpdateMedium(r.Context(), database.UpdateMediumParams{
		ID:          params.ID,
		Title:       params.Title,
		Creator:     params.Creator,
		ReleaseYear: params.ReleaseYear,
		ImageUrl:    params.ImageUrl,
		Metadata:    params.Metadata,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, 404, "No medium with given ID in database", err)
			return
		}
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			// This is a unique constraint violation
			respondWithError(w, 400, "A medium with same title already exists in database", err)
			return
		}
		respondWithError(w, 500, "couldn't udpate medium by given ID", err)
		return
	}

	// Respond
	respondWithJson(w, 200, response{
		Medium: Medium{
			ID:          medium.ID,
			MediaType:   medium.MediaType,
			CreatedAt:   medium.CreatedAt,
			UpdatedAt:   medium.UpdatedAt,
			Title:       medium.Title,
			Creator:     medium.Creator,
			ReleaseYear: medium.ReleaseYear,
			ImageUrl:    medium.ImageUrl,
			Metadata:    medium.Metadata,
		}})
}

// DELETE /api/media
func (cfg *apiConfig) handlerDeleteMedium(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		ID pgtype.UUID `json:"id"`
	}

	// Parse data from request body
	var params parameters
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		respondWithError(w, 500, "couldn't decode body from request", err)
		return
	}

	// Call query function
	count, err := cfg.db.DeleteMedium(r.Context(), params.ID)
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
