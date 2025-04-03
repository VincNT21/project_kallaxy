package kallaxyapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func (c *APIClient) makeHttpRequest(endpoint Endpoint, params interface{}) (*http.Response, error) {
	// Parse URL
	url := fmt.Sprintf("%s%s", c.Config.BaseURL, endpoint.Path)

	// Marshal the request body
	reqBody, err := json.Marshal(params)
	if err != nil {
		return &http.Response{}, fmt.Errorf("couldn't json.Marshal given body: %v", err)
	}
	// Create the request with correct headers
	req, err := http.NewRequest(endpoint.Method, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return &http.Response{}, fmt.Errorf("couldn't create http.NewRequest: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Config.AuthToken))

	// Make request
	log.Printf("--DEBUG-- Making request to %s\n", url)
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return &http.Response{}, fmt.Errorf("couldn't Do Request: %v", err)
	}

	// Return response
	return resp, nil
}

func (c *APIClient) makeHttpRequestWithResfreshToken(endpoint Endpoint) (*http.Response, error) {
	// Parse URL
	url := fmt.Sprintf("%s%s", c.Config.BaseURL, endpoint.Path)

	// Create the request with correct headers
	req, err := http.NewRequest(endpoint.Method, url, nil)
	if err != nil {
		return &http.Response{}, fmt.Errorf("couldn't create http.NewRequest: %v", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.CurrentUser.RefreshToken))

	// Make request
	log.Printf("--DEBUG-- Making request to %s\n", url)
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return &http.Response{}, fmt.Errorf("couldn't Do Request: %v", err)
	}

	// Return response
	return resp, nil
}
