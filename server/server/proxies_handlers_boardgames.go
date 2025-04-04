package server

import (
	"io"
	"log"
	"net/http"

	"github.com/clbanning/mxj/v2"
)

// GET /external_api/boardgame/search <query parameter : ?query=xxx>
func (cfg *apiConfig) handlerBoardgameSearch(w http.ResponseWriter, r *http.Request) {
	const bggSearchApiUrl = "https://boardgamegeek.com/xmlapi2/search"

	// Get query parameters
	requestQuery := "?" + r.URL.RawQuery

	// Create request
	apiURL := bggSearchApiUrl + requestQuery
	log.Printf("--DEBUG-- Making external request to %s", apiURL)
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		respondWithError(w, 500, "couldn't create Get request for BGG API", err)
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

	// Read the XML body
	xmlData, err := io.ReadAll(resp.Body)
	if err != nil {
		respondWithError(w, 500, "failed to read response body", err)
		return
	}

	// Convert XML response to JSON
	mxj.PrependAttrWithHyphen(false)

	mv, err := mxj.NewMapXml(xmlData)
	if err != nil {
		respondWithError(w, 500, "couldn't convert XML response to JSON", err)
		return
	}

	jsonData, err := mv.Json()
	if err != nil {
		respondWithError(w, 500, "couldn't convert XML response to JSON", err)
	}

	// Send the response
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

// GET /external_api/boardgame <query parameter : ?id=xxx>
func (cfg *apiConfig) handlerBoardgameDetails(w http.ResponseWriter, r *http.Request) {
	const bggDetailsApiUrl = "https://boardgamegeek.com/xmlapi2/things"

	// Get query parameters
	requestQuery := "?" + r.URL.RawQuery

	// Create request
	apiURL := bggDetailsApiUrl + requestQuery
	log.Printf("--DEBUG-- Making external request to %s", apiURL)
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		respondWithError(w, 500, "couldn't create Get request for BGG API", err)
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

	// Read the XML body
	xmlData, err := io.ReadAll(resp.Body)
	if err != nil {
		respondWithError(w, 500, "failed to read response body", err)
		return
	}

	// Convert XML response to JSON
	mxj.PrependAttrWithHyphen(false)

	mv, err := mxj.NewMapXml(xmlData)
	if err != nil {
		respondWithError(w, 500, "couldn't convert XML response to JSON", err)
		return
	}

	jsonData, err := mv.Json()
	if err != nil {
		respondWithError(w, 500, "couldn't convert XML response to JSON", err)
	}

	// Send the response
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}
