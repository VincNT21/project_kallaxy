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

func (cfg *apiConfig) handlerCreateMedium(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Title       string          `json:"title"`
		Type        string          `json:"type"`
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
		Type:        params.Type,
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
			Type:        medium.Type,
			CreatedAt:   medium.CreatedAt,
			UpdatedAt:   medium.UpdatedAt,
			Title:       medium.Title,
			Creator:     medium.Creator,
			ReleaseYear: medium.ReleaseYear,
			ImageUrl:    medium.ImageUrl,
			Metadata:    medium.Metadata,
		}})
}

func (cfg *apiConfig) handlerUpdateMedium(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		ID          pgtype.UUID     `json:"id"`
		Title       string          `json:"title"`
		Type        string          `json:"type"`
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
		respondWithError(w, 500, "couldn't udpate medium by given ID", err)
		return
	}

	// Respond
	respondWithJson(w, 200, response{
		Medium: Medium{
			ID:          medium.ID,
			Type:        medium.Type,
			CreatedAt:   medium.CreatedAt,
			UpdatedAt:   medium.UpdatedAt,
			Title:       medium.Title,
			Creator:     medium.Creator,
			ReleaseYear: medium.ReleaseYear,
			ImageUrl:    medium.ImageUrl,
			Metadata:    medium.Metadata,
		}})
}

func (cfg *apiConfig) handlerGetMediumByTitle(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Title string `json:"title"`
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
	medium, err := cfg.db.GetMediumByTitle(r.Context(), params.Title)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, 404, fmt.Sprintf("No medium with title %s in database", params.Title), err)
			return
		}
		respondWithError(w, 500, "couldn't get medium by title", err)
		return
	}

	// Respond
	respondWithJson(w, 200, response{
		Medium: Medium{
			ID:          medium.ID,
			Type:        medium.Type,
			CreatedAt:   medium.CreatedAt,
			UpdatedAt:   medium.UpdatedAt,
			Title:       medium.Title,
			Creator:     medium.Creator,
			ReleaseYear: medium.ReleaseYear,
			ImageUrl:    medium.ImageUrl,
			Metadata:    medium.Metadata,
		}})

}

func (cfg *apiConfig) handlerGetMediaByType(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Type string `json:"type"`
	}

	type response struct {
		Media []Medium `json:"media"`
	}

	// Parse data from request body
	var params parameters
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		respondWithError(w, 500, "couldn't decode body from request", err)
		return
	}

	// Call query function
	media, err := cfg.db.GetMediaByType(r.Context(), params.Type)
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
			Type:        medium.Type,
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
