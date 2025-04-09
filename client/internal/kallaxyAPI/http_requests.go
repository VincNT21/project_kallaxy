package kallaxyapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/VincNT21/kallaxy/client/models"
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

	// Check response's status code
	if resp.StatusCode != 200 && resp.StatusCode != 201 && resp.StatusCode != 203 {
		log.Printf("--ERROR-- with %s request to %s. Response status code: %v\n", endpoint.Method, endpoint.Path, resp.StatusCode)
		switch resp.StatusCode {
		case 400:
			return nil, models.ErrBadRequest
		case 401:
			return nil, models.ErrUnauthorized
		case 404:
			return nil, models.ErrNotFound
		case 409:
			return nil, models.ErrConflict
		case 500:
			return nil, models.ErrServerIssue
		default:
			return nil, fmt.Errorf("unknown error status code: %v", resp.StatusCode)
		}
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

	// Check response's status code
	if resp.StatusCode != 200 && resp.StatusCode != 201 && resp.StatusCode != 203 {
		log.Printf("--ERROR-- with %s request to %s. Response status code: %v\n", endpoint.Method, endpoint.Path, resp.StatusCode)
		switch resp.StatusCode {
		case 400:
			return nil, models.ErrBadRequest
		case 401:
			return nil, models.ErrUnauthorized
		case 404:
			return nil, models.ErrNotFound
		case 409:
			return nil, models.ErrConflict
		case 500:
			return nil, models.ErrServerIssue
		default:
			return nil, fmt.Errorf("unknown error status code: %v", resp.StatusCode)
		}
	}

	// Return response
	return resp, nil
}

func (c *APIClient) makeHttpRequestWithQueryParameters(endpoint Endpoint, queryParameter string) (*http.Response, error) {
	// Parse URL
	url := fmt.Sprintf("%s%s?%s", c.Config.BaseURL, endpoint.Path, queryParameter)

	// Create the request with correct headers
	req, err := http.NewRequest(endpoint.Method, url, nil)
	if err != nil {
		return &http.Response{}, fmt.Errorf("couldn't create http.NewRequest: %v", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Config.AuthToken))

	// Make request
	log.Printf("--DEBUG-- Making request to %s\n", url)
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return &http.Response{}, fmt.Errorf("couldn't Do Request: %v", err)
	}
	// Check response's status code
	if resp.StatusCode != 200 && resp.StatusCode != 201 && resp.StatusCode != 203 {
		log.Printf("--ERROR-- with %s request to %s. Response status code: %v\n", endpoint.Method, endpoint.Path, resp.StatusCode)
		switch resp.StatusCode {
		case 400:
			return nil, models.ErrBadRequest
		case 401:
			return nil, models.ErrUnauthorized
		case 404:
			return nil, models.ErrNotFound
		case 409:
			return nil, models.ErrConflict
		case 500:
			return nil, models.ErrServerIssue
		default:
			return nil, fmt.Errorf("unknown error status code: %v", resp.StatusCode)
		}
	}

	// Return response
	return resp, nil
}
