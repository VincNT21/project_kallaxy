package context

import (
	"fyne.io/fyne/v2"
	kallaxyapi "github.com/VincNT21/kallaxy/client/internal/kallaxyAPI"
)

type AppContext struct {
	APIClient           *kallaxyapi.APIClient
	MainWindow          fyne.Window
	PageManager         PageManager
	MetadataFieldsMap   map[string][]string
	MetadataFieldsSpecs map[string]kallaxyapi.FieldSpec
}

type PageManager interface {
	ShowLoginPage()
	ShowWelcomeBackPage()
	ShowHomePage()
	ShowCreateUserPage()
	ShowCreateMediaPage(mediaType string)
	ShowShelfPage()
	ShowParametersPage()
}

// Create and configured the shared AppContext
func NewAppContext(baseURL string) *AppContext {

	// Init apiClient
	apiClient := kallaxyapi.NewApiClient(baseURL)

	// Load Metadata models from local field

	return &AppContext{
		APIClient: apiClient,
	}
}
