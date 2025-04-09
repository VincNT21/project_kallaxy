package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

// GET /external_api/movie_tv/search_movie
func (cfg *apiConfig) handlerMovieSearch(w http.ResponseWriter, r *http.Request) {
	const movieDbSearchUrl = "https://api.themoviedb.org/3/search/movie"

	// Get query parameters
	requestQuery := "?" + r.URL.RawQuery

	// Create request
	apiURL := movieDbSearchUrl + requestQuery
	log.Printf("--DEBUG-- Making external request to %s", apiURL)
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		respondWithError(w, 500, "couldn't create Get request for The Movie DB API", err)
		return
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", cfg.moviedbKey))

	// Make request to external API
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		respondWithError(w, 500, "failed to fetch data", err)
		return
	}
	defer resp.Body.Close()

	// Pass through the response
	w.Header().Set("Content-Type", "application/json")
	io.Copy(w, resp.Body)
}

// GET /external_api/movie_tv/search_tv
func (cfg *apiConfig) handlerTVSearch(w http.ResponseWriter, r *http.Request) {
	const movieDbSearchUrl = "https://api.themoviedb.org/3/search/tv"

	// Get query parameters
	requestQuery := "?" + r.URL.RawQuery

	// Create request
	apiURL := movieDbSearchUrl + requestQuery
	log.Printf("--DEBUG-- Making external request to %s", apiURL)
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		respondWithError(w, 500, "couldn't create Get request for The Movie DB API", err)
		return
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", cfg.moviedbKey))

	// Make request to external API
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		respondWithError(w, 500, "failed to fetch data", err)
		return
	}
	defer resp.Body.Close()

	// Pass through the response
	w.Header().Set("Content-Type", "application/json")
	io.Copy(w, resp.Body)
}

// GET /external_api/movie/search
func (cfg *apiConfig) handlerMultiSearch(w http.ResponseWriter, r *http.Request) {
	const movieDbSearchUrl = "https://api.themoviedb.org/3/search/multi"

	// Get query parameters
	requestQuery := "?" + r.URL.RawQuery

	// Create request
	apiURL := movieDbSearchUrl + requestQuery
	log.Printf("--DEBUG-- Making external request to %s", apiURL)
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		respondWithError(w, 500, "couldn't create Get request for The Movie DB API", err)
		return
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", cfg.moviedbKey))

	// Make request to external API
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		respondWithError(w, 500, "failed to fetch data", err)
		return
	}
	defer resp.Body.Close()

	// Pass through the response
	w.Header().Set("Content-Type", "application/json")
	io.Copy(w, resp.Body)
}

type parametersMovieDetails struct {
	MovieID  string `json:"movie_id"`
	TvID     string `json:"tv_id"`
	Language string `json:"language"`
}

// GET /external_api/movie_tv
func (cfg *apiConfig) handlerMovieTvDetails(w http.ResponseWriter, r *http.Request) {
	const movieDbMovieDetailsUrl = "https://api.themoviedb.org/3/movie"
	const movieDbTvDetailsUrl = "https://api.themoviedb.org/3/tv"

	// Parse data from request body
	var params parametersMovieDetails
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		respondWithError(w, 500, "couldn't decode body from request", err)
		return
	}

	// Create request
	apiURL := ""
	if params.MovieID != "" {
		apiURL = movieDbMovieDetailsUrl + "/" + params.MovieID + "?language" + params.Language
	} else if params.TvID != "" {
		apiURL = movieDbTvDetailsUrl + "/" + params.TvID + "?language" + params.Language
	} else {
		respondWithError(w, 400, "no movie id or tv id in request body", errors.New("both field 'movie_id' and 'tv_id' are empty string"))
		return
	}
	log.Printf("--DEBUG-- Making external request to %s", apiURL)
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		respondWithError(w, 500, "couldn't create Get request for The Movie DB API", err)
		return
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", cfg.moviedbKey))

	// Make request to external API
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		respondWithError(w, 500, "failed to fetch data", err)
		return
	}
	defer resp.Body.Close()

	// Pass through the response
	w.Header().Set("Content-Type", "application/json")
	io.Copy(w, resp.Body)
}

// GET /external_api/movie_tv/movie_credits (query parameters: "?movie_id=xxxx")
func (cfg *apiConfig) handlerMovieCreditsDetails(w http.ResponseWriter, r *http.Request) {
	const movieDbMovieDetailsUrl = "https://api.themoviedb.org/3/movie"
	// const movieDbTvCreditsDetailsUrl = " https://api.themoviedb.org/3/tv/{series_id}/season/{season_number}/credits"

	// Get movie ID from request query parameters
	movieID := r.URL.Query().Get("movie_id")

	// Create request
	creditsUrl := movieDbMovieDetailsUrl + fmt.Sprintf("/%v/credits", url.QueryEscape(movieID))
	req, err := http.NewRequest("GET", creditsUrl, nil)
	if err != nil {
		respondWithError(w, 500, "couldn't create Get request for The Movie DB API", err)
		return
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", cfg.moviedbKey))

	// Make request to external API
	log.Printf("--DEBUG-- Making external request to %s", creditsUrl)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		respondWithError(w, 500, "failed to fetch data", err)
		return
	}
	defer resp.Body.Close()

	// Pass through the response
	w.Header().Set("Content-Type", "application/json")
	io.Copy(w, resp.Body)

}
