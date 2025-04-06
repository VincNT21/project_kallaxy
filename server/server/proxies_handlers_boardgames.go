package server

import (
	"encoding/json"
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

	// Problem to handle : if the search returns only one item,
	// `item` field will be a struct of one item and not a list

	// First try to unmarshal as if it has multiple items
	var multiResponse responseBoardgameSearch

	jsonData, err := mv.Json()
	if err != nil {
		respondWithError(w, 500, "couldn't convert XML response to JSON", err)
		return
	}

	err = json.Unmarshal(jsonData, &multiResponse)

	// Check if Items.Item is empty (which means it might be a single item response)
	if err != nil || len(multiResponse.Items.Item) == 0 {
		// Try the alternative structure with a single item
		var singleResponse responseBoardgameSearchAlternative
		err = json.Unmarshal(jsonData, &singleResponse)
		if err != nil {
			respondWithError(w, 500, "failed to parse response", err)
			return
		}

		// Convert single item to multi-item format
		multiResponse.Items.Total = singleResponse.Items.Total
		multiResponse.Items.Item = []struct {
			ID   string `json:"id"`
			Name struct {
				Value string `json:"value"`
			} `json:"name"`
			Type          string `json:"type"`
			Yearpublished struct {
				Value string `json:"value"`
			} `json:"yearpublished"`
		}{
			{
				ID:            singleResponse.Items.Item.ID,
				Name:          singleResponse.Items.Item.Name,
				Type:          singleResponse.Items.Item.Type,
				Yearpublished: singleResponse.Items.Item.Yearpublished,
			},
		}
	}

	// Now multiResponse always has the array format, regardless of original XML
	responseJSON, err := json.Marshal(multiResponse)
	if err != nil {
		respondWithError(w, 500, "failed to serialize response", err)
		return
	}

	// Send the response
	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJSON)
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
