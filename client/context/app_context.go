package context

import (
	"fyne.io/fyne/v2"
	kallaxyapi "github.com/VincNT21/kallaxy/client/internal/kallaxyAPI"
	"github.com/VincNT21/kallaxy/client/models"
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
	ShowCompartmentTreePage(mediaType string, mediaList []models.MediumWithRecord)
	ShowUpdateMediaPage(mediaType, mediumID string, mediaList []models.MediumWithRecord)
	ShowUpdateRecordPage(mediaType, mediumID string, mediaList []models.MediumWithRecord)
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
