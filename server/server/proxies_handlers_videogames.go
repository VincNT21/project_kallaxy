package server

import (
	"io"
	"log"
	"net/http"
)

// GET /external_api/videogame/search (query parameter: ?search=<title>&platforms=<platformsID>)
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

// GET /external_api/videogame (query parameter: ?id=<gameID>)
func (cfg *apiConfig) handlerVideoGameDetails(w http.ResponseWriter, r *http.Request) {
	const RAWGUrl = "https://api.rawg.io/api/games"

	// Get game ID from request query parameter
	gameID := r.URL.Query().Get("id")

	// Create request
	apiURL := RAWGUrl + "/" + gameID + "?key=" + cfg.rawgKey
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
