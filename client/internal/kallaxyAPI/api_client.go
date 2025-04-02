package kallaxyapi

import (
	"net/http"
	"time"

	"github.com/VincNT21/kallaxy/client/models"
)

type APIClient struct {
	HttpClient *http.Client
	Config     *APIConfig
	LastUser   models.ClientUser

	Users    *UsersClient
	Media    *MediaClient
	Records  *RecordsClient
	Auth     *AuthClient
	External *ExternalAPIClient
	Admin    *AdminClient
}

type UsersClient struct {
	apiClient *APIClient // Reference back to the parent
}

type MediaClient struct {
	apiClient *APIClient // Reference back to the parent
}

type RecordsClient struct {
	apiClient *APIClient // Reference back to the parent
}

type AuthClient struct {
	apiClient *APIClient // Reference back to the parent
}

type ExternalAPIClient struct {
	apiClient *APIClient // Reference back to the parent
}

type AdminClient struct {
	apiClient *APIClient // Reference back to the parent
}

// Constructs the APIClient using api_config
func NewApiClient(baseURL string) *APIClient {
	// Create a http.Client
	httpClient := &http.Client{
		Timeout: time.Second * 10,
	}
	// Get proper initialized api Config
	apiCfg := initApiConfig(baseURL)

	// Initialize the APIClient itself
	apiClient := &APIClient{
		HttpClient: httpClient,
		Config:     apiCfg,
	}

	// Initialize the subclients and give them access to the parent
	apiClient.Users = &UsersClient{apiClient: apiClient}
	apiClient.Media = &MediaClient{apiClient: apiClient}
	apiClient.Records = &RecordsClient{apiClient: apiClient}
	apiClient.Auth = &AuthClient{apiClient: apiClient}
	apiClient.External = &ExternalAPIClient{apiClient: apiClient}
	apiClient.Admin = &AdminClient{apiClient: apiClient}

	return apiClient
}
