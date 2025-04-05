package context

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	kallaxyapi "github.com/VincNT21/kallaxy/client/internal/kallaxyAPI"
	"github.com/VincNT21/kallaxy/client/models"
)

type PageManager interface {
	GetLoginWindow()
	GetBackWindow()
	GetCreateUserWindow(func())
	GetHomeWindow()
	GetUserParametersWindow()
	GetCreateMediaWindow()
	ShowImageWindow(fyne.Window, string, func(string))
	GetShelfWindow()
	BuildMediaContainers(models.MediaWithRecords) (*container.Scroll, error)
}

type AppContext struct {
	APIClient   *kallaxyapi.APIClient
	PageManager PageManager
}

// local storage of appState is a temporary approach
const localStorageAppStatePath = "/home/vincnt/workspace/project_kallaxy/client/config/appstate.json"

// Create and configured the shared AppContext
func NewAppContext(baseURL string) *AppContext {

	// Init apiClient
	apiClient := kallaxyapi.NewApiClient(baseURL)

	return &AppContext{
		APIClient: apiClient,
	}
}

// Unused for now. Use a fixed path temporarily
func getClientConfigDir() string {
	execPath, err := os.UserConfigDir()
	if err == nil {
		return filepath.Join(filepath.Dir(execPath), "config")
	}

	return "./config"
}

// Check if appstate data exists and loads it
func (c *AppContext) LoadsAppstate() {
	/*
		configDir := getClientConfigDir()

		filepath := filepath.Join(configDir, "appstate.json")

	*/
	f, err := os.Open(localStorageAppStatePath)
	if err != nil {
		log.Printf("--ERROR-- with LoadsAppstate(), couldn't open appstate.json: %v\n", err)
		return
	}
	defer f.Close()

	// Read data from file
	var user models.ClientUser
	err = json.NewDecoder(f).Decode(&user)
	if err != nil {
		log.Printf("--ERROR-- with LoadsAppstate(), couldn't decode data from appstate.json: %v\n", err)
		return
	}

	// Load data into APIClient
	c.APIClient.CurrentUser = user
}

// Store appstate data in local file
func (c *AppContext) DumpAppstate() {
	/*
		configDir := getClientConfigDir()

		// Ensure directory exists
		if err := os.MkdirAll(configDir, 0755); err != nil {
			log.Printf("--ERROR-- couldn't create client config directory: %v", err)
			return
		}

		filepath := filepath.Join(configDir, "appstate.json")
	*/

	// Create/erase local appstate file
	f, err := os.Create(localStorageAppStatePath)
	if err != nil {
		log.Printf("--ERROR-- with DumpAppstate(), couldn't create appstate.json: %v\n", err)
		return
	}
	defer f.Close()

	// Get data from APIClient
	data, err := json.Marshal(c.APIClient.CurrentUser)
	if err != nil {
		log.Printf("--ERROR-- with DumpAppstate(), couldn't json.Marshal data for appstate.json: %v\n", err)
		return
	}

	// Write to file
	_, err = f.Write(data)
	if err != nil {
		log.Printf("--ERROR-- with DumpAppstate(), couldn't write data in appstate.json: %v\n", err)
		return
	}
}
