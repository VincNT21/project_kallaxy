package context

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"runtime"

	kallaxyapi "github.com/VincNT21/kallaxy/client/internal/kallaxyAPI"
	"github.com/VincNT21/kallaxy/client/models"
)

func getLocalConfigStoragePath() string {
	outputDir := ""

	// Check if OS used is Windows or else
	if runtime.GOOS == "windows" {
		// For windows builds, Get the path to the currently running executable
		execPath, err := os.Executable()
		if err != nil {
			log.Fatalf("Error with getLocalConfigStoragePath(): %v", err)
		} else {
			outputDir = filepath.Dir(execPath)
		}
	} else {
		// For Linux/Mac, use working directory + client directory
		workingDir, err := os.Getwd()
		if err != nil {
			log.Fatalf("Error with getLocalConfigStoragePath(): %v", err)
		} else {
			outputDir = filepath.Join(workingDir, "client")
		}
	}

	// If we couldn't determine a directory, fall back to current directory
	if outputDir == "" {
		outputDir = "."
	}

	// Build the path to the "config/" folder
	return filepath.Join(outputDir, "config")
}

type AppState struct {
	LastUser      models.ClientUser `json:"last_user"`
	ClientVersion string            `json:"client_version"`
	ServerVersion string            `json:"server_version"`
}

// Check if appstate data exists and loads it
func (c *AppContext) LoadAppstate() {
	localStoragePath := getLocalConfigStoragePath()
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
	localStoragePath := getLocalConfigStoragePath()
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
	localStoragePath := getLocalConfigStoragePath()
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
	localStoragePath := getLocalConfigStoragePath()
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
	localStoragePath := getLocalConfigStoragePath()
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
	localStoragePath := getLocalConfigStoragePath()
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
