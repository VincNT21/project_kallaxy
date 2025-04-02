package context

import (
	"net/http"
	"time"
)

type AppContext struct {
	HttpClient *http.Client
	Cache      string // To be implemented
	APIClient  string // To be implemented
	UserToken  string
	Config     APIConfig
}

// Initializes everything the GUI and internal API need
func NewAppContext() *AppContext {
	// Create and configured the shared AppContext
	httpClient := &http.Client{
		Timeout: time.Second * 10,
	}
	cache := "cache.NewCache()"
	apiClient := "kallaxyapi.NewClient(httpClient)"

	return &AppContext{
		HttpClient: httpClient,
		Cache:      cache,
		APIClient:  apiClient,
	}
}
