package context

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	kallaxyapi "github.com/VincNT21/kallaxy/client/internal/kallaxyAPI"
	"github.com/VincNT21/kallaxy/client/models"
)

type PageManager interface {
	GetBackWindow()
	GetLoginWindow()
	GetCreateUserWindow()
	GetHomeWindow()
}

type AppContext struct {
	Cache       string // To be implemented
	APIClient   *kallaxyapi.APIClient
	PageManager PageManager
}

// Create and configured the shared AppContext
func NewAppContext(baseURL string) *AppContext {

	cache := "cache.NewCache()"
	apiClient := kallaxyapi.NewApiClient(baseURL)

	return &AppContext{
		Cache:     cache,
		APIClient: apiClient,
	}
}

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
	f, err := os.Open("/home/vincnt/workspace/project_kallaxy/client/config/appstate.json")
	if err != nil {
		log.Printf("couldn't open appstate.json: %v", err)
		return
	}
	defer f.Close()

	// Read data from file
	var user models.ClientUser
	err = json.NewDecoder(f).Decode(&user)
	if err != nil {
		log.Printf("couldn't decode data from appstate.json: %v", err)
		return
	}

	// Load data into APIClient
	c.APIClient.LastUser = user
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
	f, err := os.Create("/home/vincnt/workspace/project_kallaxy/client/config/appstate.json")
	if err != nil {
		log.Printf("couldn't create appstate.json: %v", err)
		return
	}
	defer f.Close()

	// Get data from APIClient
	data, err := json.Marshal(c.APIClient.LastUser)
	if err != nil {
		log.Printf("couldn't json.Marshal data for appstate.json: %v", err)
		return
	}

	// Write to file
	_, err = f.Write(data)
	if err != nil {
		log.Printf("couldn't write data in appstate.json: %v", err)
		return
	}
}
