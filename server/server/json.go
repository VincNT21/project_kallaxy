package server

import (
	"encoding/json"
	"log"
	"net/http"
)

// Write a proper response JSON formatted
func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("--ERROR-- Couldn't marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(code)
	w.Write(data)
}

// Format a proper error response payload and call respond with JSON
func respondWithError(w http.ResponseWriter, code int, msg string, err error) {
	if err != nil {
		log.Printf("--ERROR-- %v", err)
		log.Printf("--INFO-- Responding with %d error status code: %s", code, msg)
	}

	type errorResponse struct {
		Error string `json:"error"`
	}

	respondWithJson(w, code, errorResponse{
		Error: msg,
	})
}
