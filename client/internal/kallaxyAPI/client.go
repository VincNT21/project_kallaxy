package kallaxyapi

import "net/http"

type APIClient struct {
	HttpClient *http.Client
	Config     APIConfig

	Users    *UsersClient
	Media    *MediaClient
	Records  *RecordsClient
	Auth     *AuthClient
	External *ExternalAPIClient
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

// Constructs the APIClient using api_config
// func NewApiClient()

// Initialize the APIClient itself
/*
apiClient := &APIClient{
	BaseURL:    baseURL,
	HTTPClient: httpClient,
}
*/

// Initialize the subclients and give them access to the parent
// apiClient.Users = &UsersClient{apiClient: apiClient}
