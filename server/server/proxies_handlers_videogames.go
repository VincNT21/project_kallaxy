package server

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

// GET /external_api/videogame/search
func (cfg *apiConfig) handlerVideoGameSearch(w http.ResponseWriter, r *http.Request) {
	const RAWGSearchUrl = "https://api.rawg.io/api/games"

	// Get query parameters
	requestQuery := r.URL.RawQuery

	// Create request
	apiURL := RAWGSearchUrl + "?key=" + cfg.rawgKey + "&" + requestQuery + "&search_precise=true"
	log.Printf("--DEBUG-- Making external request to %s", apiURL)
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		respondWithError(w, 500, "couldn't create Get request for RAWG API", err)
		return
	}

	// Make request to external API
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		respondWithError(w, 500, "failed to fetch data", err)
		return
	}
	defer resp.Body.Close()

	// Pass though the response
	w.Header().Set("Content-Type", "application/json")
	io.Copy(w, resp.Body)
}

type parametersVideogameDetails struct {
	GameID string `json:"game_id"`
}

// GET /external_api/videogame
func (cfg *apiConfig) handlerVideoGameDetails(w http.ResponseWriter, r *http.Request) {
	const RAWGUrl = "https://api.rawg.io/api/games"
	// Parse data from request body
	var params parametersVideogameDetails
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		respondWithError(w, 500, "couldn't decode body from request", err)
		return
	}

	// Create request
	apiURL := RAWGUrl + "/" + params.GameID + "?key=" + cfg.rawgKey
	log.Printf("--DEBUG-- Making external request to %s", apiURL)
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		respondWithError(w, 500, "couldn't create Get request for RAWG API", err)
		return
	}

	// Make request to external API
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		respondWithError(w, 500, "failed to fetch data", err)
		return
	}
	defer resp.Body.Close()

	// Pass though the response
	w.Header().Set("Content-Type", "application/json")
	io.Copy(w, resp.Body)
}
