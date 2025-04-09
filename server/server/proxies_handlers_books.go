package server

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

// GET /external_api/book/search (query parameters: "?title=xxx" / "?author=xxx" ...)
func (cfg *apiConfig) handlerBookSearch(w http.ResponseWriter, r *http.Request) {
	const openLibrarySearchUrl = "https://openlibrary.org/search.json"

	// Get query parameters
	requestQuery := "?" + r.URL.RawQuery

	// Create request
	apiURL := openLibrarySearchUrl + requestQuery
	log.Printf("--DEBUG-- Making external request to %s", apiURL)
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		respondWithError(w, 500, "couldn't create Get request for Open Library API", err)
		return
	}
	req.Header.Set("User-Agent", cfg.openlibraryUA)

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
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

// GET /external_api/book/isbn (query parameters: "?isbn=xxxx")
func (cfg *apiConfig) handlerBookByISBN(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		ISBN string
	}
	const openLibraryIsbnUrl = "https://openlibrary.org/isbn/"

	// Get query parameter
	var p parameters
	p.ISBN = r.URL.Query().Get("isbn")

	// Create request
	apiURL := openLibraryIsbnUrl + p.ISBN + ".json"
	log.Printf("--DEBUG-- Making external request to %s", apiURL)
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		respondWithError(w, 500, "couldn't create Get request for Open Library API", err)
		return
	}
	req.Header.Set("User-Agent", cfg.openlibraryUA)

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
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

// GET /external_api/book/author (query parameters: "?author=xxxx")
func (cfg *apiConfig) handlerBookAuthor(w http.ResponseWriter, r *http.Request) {
	const openLibraryIsbnUrl = "https://openlibrary.org/author/"

	// Get query parameter
	author := r.URL.Query().Get("author")

	// Create request
	apiURL := openLibraryIsbnUrl + author + ".json"
	log.Printf("--DEBUG-- Making external request to %s", apiURL)
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		respondWithError(w, 500, "couldn't create Get request for Open Library API", err)
		return
	}
	req.Header.Set("User-Agent", cfg.openlibraryUA)

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
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

// GET /external_api/book/search_isbn (query parameters: "?key=xxxx")
func (cfg *apiConfig) handlerGetBookISBN(w http.ResponseWriter, r *http.Request) {
	type responseISBN struct {
		ISBN10 string `json:"isbn10"`
		ISBN13 string `json:"isbn13"`
	}

	type responseFromApiExtract struct {
		Entries []struct {
			Isbn10 []string `json:"isbn_10"`
			Isbn13 []string `json:"isbn_13"`
		} `json:"entries"`
	}

	const openLibraryWorksUrl = "https://openlibrary.org/works/"

	// Get query parameter
	key := r.URL.Query().Get("key")

	// Create request
	apiURL := openLibraryWorksUrl + key + "/editions.json"
	log.Printf("--DEBUG-- Making external request to %s", apiURL)
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		respondWithError(w, 500, "couldn't create Get request for Open Library API", err)
		return
	}
	req.Header.Set("User-Agent", cfg.openlibraryUA)

	// Make request to external API
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		respondWithError(w, 500, "failed to fetch data", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		respondWithError(w, resp.StatusCode, fmt.Sprintf("request to %s return status code: %v", apiURL, resp.StatusCode), fmt.Errorf("request to %s return status code: %v", apiURL, resp.StatusCode))
		return
	}

	// Parse response from API
	var bookInfo responseFromApiExtract
	err = json.NewDecoder(resp.Body).Decode(&bookInfo)
	if err != nil || len(bookInfo.Entries) == 0 {
		respondWithError(w, 500, "failed to decode Open Library API response", err)
		return
	}

	// Create response for client
	var response responseISBN
	if len(bookInfo.Entries[0].Isbn10) == 0 {
		response.ISBN10 = ""
	} else {
		response.ISBN10 = bookInfo.Entries[0].Isbn10[0]
	}
	if len(bookInfo.Entries[0].Isbn13) == 0 {
		response.ISBN13 = ""
	} else {
		response.ISBN13 = bookInfo.Entries[0].Isbn13[0]
	}

	// Send the response
	respondWithJson(w, 200, response)

}
