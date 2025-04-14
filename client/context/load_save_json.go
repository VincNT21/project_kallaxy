package context

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	kallaxyapi "github.com/VincNT21/kallaxy/client/internal/kallaxyAPI"
	"github.com/VincNT21/kallaxy/client/models"
)

func getLocalStoragePath() string {
	// Get the path to the currently running executable
	execPath, err := os.Executable()
	if err != nil {
		log.Fatalf("Error with getLocalStoragePath(): %v", err)
	}

	// Determine the directory where the executable is located
	execDir := filepath.Dir(execPath)

	// Build the path to the "config/" folder
	fmt.Println(filepath.Join(execDir, "config"))
	return filepath.Join(execDir, "config")
}

type AppState struct {
	LastUser      models.ClientUser `json:"last_user"`
	ClientVersion string            `json:"client_version"`
	ServerVersion string            `json:"server_version"`
}

// Check if appstate data exists and loads it
func (c *AppContext) LoadAppstate() {
	localStoragePath := getLocalStoragePath()
	appStateFilePath := filepath.Join(localStoragePath, "appstate.json")
	f, err := os.Open(appStateFilePath)
	if err != nil {
		log.Printf("--ERROR-- with LoadAppstate(), couldn't open appstate.json: %v\n", err)
		return
	}
	defer f.Close()

	// Read data from file
	var appState AppState
	err = json.NewDecoder(f).Decode(&appState)
	if err != nil {
		log.Printf("--ERROR-- with LoadAppstate(), couldn't decode data from appstate.json: %v\n", err)
		return
	}

	// Load data into APIClient
	log.Println("--DEBUG-- LoadAppstate() OK")
	c.APIClient.CurrentUser = appState.LastUser
	c.APIClient.ClientVersion = appState.ClientVersion
	c.APIClient.ServerVersion = appState.ServerVersion
}

// Store appstate data in local file
func (c *AppContext) SaveAppstate() {
	localStoragePath := getLocalStoragePath()
	appStateFilePath := filepath.Join(localStoragePath, "appstate.json")
	// Create/erase local appstate file
	f, err := os.Create(appStateFilePath)
	if err != nil {
		log.Printf("--ERROR-- with SaveAppstate(), couldn't create appstate.json: %v\n", err)
		return
	}
	defer f.Close()

	// Get data from APIClient
	appState := AppState{
		LastUser:      c.APIClient.CurrentUser,
		ClientVersion: c.APIClient.ClientVersion,
		ServerVersion: c.APIClient.ServerVersion,
	}
	data, err := json.Marshal(appState)
	if err != nil {
		log.Printf("--ERROR-- with SaveAppstate(), couldn't json.Marshal data for appstate.json: %v\n", err)
		return
	}

	// Write to file
	_, err = f.Write(data)
	if err != nil {
		log.Printf("--ERROR-- with SaveAppstate(), couldn't write data in appstate.json: %v\n", err)
		return
	}
}

func (c *AppContext) LoadMetadataFieldsSpecs() error {
	localStoragePath := getLocalStoragePath()
	metadataFieldsSpecsFilePath := filepath.Join(localStoragePath, "metadata_fields_specs.json")

	var fieldSpecs map[string]kallaxyapi.FieldSpec

	data, err := os.ReadFile(metadataFieldsSpecsFilePath)
	if err != nil {
		log.Printf("--ERROR-- with LoadMetadataFieldSpecs(), couldn't json.Marshal data for metadata_fields_specs.json: %v\n", err)
		return err
	}

	err = json.Unmarshal(data, &fieldSpecs)
	if err != nil {
		log.Printf("--ERROR-- with LoadMetadataFieldSpecs(), couldn't write data in metadata_fields_specs.json: %v\n", err)
		return err
	}

	c.MetadataFieldsSpecs = fieldSpecs

	return nil
}

func (c *AppContext) SaveMetadataFieldsSpecs() {
	localStoragePath := getLocalStoragePath()
	metadataFieldsSpecsFilePath := filepath.Join(localStoragePath, "metadata_fields_specs.json")

	jsonData, err := json.MarshalIndent(c.MetadataFieldsSpecs, "", "  ")
	if err != nil {
		log.Printf("--ERROR-- with SaveMetadataFieldSpecs(), couldn't write data in metadata_fields_specs.json: %v\n", err)
		return
	}
	err = os.WriteFile(metadataFieldsSpecsFilePath, jsonData, 0644)
	if err != nil {
		log.Printf("--ERROR-- with SaveMetadataFieldSpecs(), couldn't write data in metadata_fields_specs.json: %v\n", err)
		return
	}

}

func (c *AppContext) LoadMetadataFieldsMap() error {
	localStoragePath := getLocalStoragePath()
	metadataFieldsMapFile := filepath.Join(localStoragePath, "metadata_fields_map.json")
	var fieldsMap map[string][]string

	data, err := os.ReadFile(metadataFieldsMapFile)
	if err != nil {
		log.Printf("--ERROR-- with LoadMetadataFieldsMap(), couldn't json.Marshal data for metadata_fields_specs.json: %v\n", err)
		return err
	}

	err = json.Unmarshal(data, &fieldsMap)
	if err != nil {
		log.Printf("--ERROR-- with LoadMetadataFieldsMap(), couldn't write data in metadata_fields_specs.json: %v\n", err)
		return err
	}

	c.MetadataFieldsMap = fieldsMap

	return nil
}

func (c *AppContext) SaveMedataFieldsMap() {
	localStoragePath := getLocalStoragePath()
	metadataFieldsMapFile := filepath.Join(localStoragePath, "metadata_fields_map.json")

	jsonData, err := json.MarshalIndent(c.MetadataFieldsMap, "", "  ")
	if err != nil {
		log.Printf("--ERROR-- with SaveMedataFieldsMap(), couldn't write data in metadata_fields_specs.json: %v\n", err)
		return
	}
	err = os.WriteFile(metadataFieldsMapFile, jsonData, 0644)
	if err != nil {
		log.Printf("--ERROR-- with SaveMedataFieldsMap(), couldn't write data in metadata_fields_specs.json: %v\n", err)
		return
	}

}
