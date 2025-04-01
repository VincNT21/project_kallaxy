package server

import (
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

	// Pass though the response
	w.Header().Set("Content-Type", "application/json")
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

	// Pass though the response
	w.Header().Set("Content-Type", "application/json")
	io.Copy(w, resp.Body)
}

// GET /external_api/book/author (query parameters: "?author=xxxx")
func (cfg *apiConfig) handlerBookAuthor(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Author string
	}
	const openLibraryIsbnUrl = "https://openlibrary.org/author/"

	// Get query parameter
	var p parameters
	p.Author = r.URL.Query().Get("author")

	// Create request
	apiURL := openLibraryIsbnUrl + p.Author + ".json"
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

	// Pass though the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	io.Copy(w, resp.Body)
}
