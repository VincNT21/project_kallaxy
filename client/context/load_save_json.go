package context

import (
	"encoding/json"
	"log"
	"os"

	kallaxyapi "github.com/VincNT21/kallaxy/client/internal/kallaxyAPI"
	"github.com/VincNT21/kallaxy/client/models"
)

// local storage is a temporary approach
const localStorageAppStatePath = "/home/vincnt/workspace/project_kallaxy/client/config/appstate.json"
const localStorageMetadataFieldsSpecsPath = "/home/vincnt/workspace/project_kallaxy/client/config/metadata_fields_specs.json"
const localStorageMetadataFieldsMapPath = "/home/vincnt/workspace/project_kallaxy/client/config/metadata_fields_map.json"

// Check if appstate data exists and loads it
func (c *AppContext) LoadAppstate() {
	f, err := os.Open(localStorageAppStatePath)
	if err != nil {
		log.Printf("--ERROR-- with LoadAppstate(), couldn't open appstate.json: %v\n", err)
		return
	}
	defer f.Close()

	// Read data from file
	var user models.ClientUser
	err = json.NewDecoder(f).Decode(&user)
	if err != nil {
		log.Printf("--ERROR-- with LoadAppstate(), couldn't decode data from appstate.json: %v\n", err)
		return
	}

	// Load data into APIClient
	log.Println("--DEBUG-- LoadAppstate() OK")
	c.APIClient.CurrentUser = user
}

// Store appstate data in local file
func (c *AppContext) SaveAppstate() {
	// Create/erase local appstate file
	f, err := os.Create(localStorageAppStatePath)
	if err != nil {
		log.Printf("--ERROR-- with SaveAppstate(), couldn't create appstate.json: %v\n", err)
		return
	}
	defer f.Close()

	// Get data from APIClient
	data, err := json.Marshal(c.APIClient.CurrentUser)
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

func (c *AppContext) LoadMetadataFieldsSpecs() {
	var fieldSpecs map[string]kallaxyapi.FieldSpec

	data, err := os.ReadFile(localStorageMetadataFieldsSpecsPath)
	if err != nil {
		log.Printf("--ERROR-- with LoadMetadataFieldSpecs(), couldn't json.Marshal data for metadata_fields_specs.json: %v\n", err)
		return
	}

	err = json.Unmarshal(data, &fieldSpecs)
	if err != nil {
		log.Printf("--ERROR-- with LoadMetadataFieldSpecs(), couldn't write data in metadata_fields_specs.json: %v\n", err)
		return
	}

	c.MetadataFieldsSpecs = fieldSpecs
}

func (c *AppContext) SaveMetadataFieldsSpecs() {

	jsonData, err := json.MarshalIndent(c.MetadataFieldsSpecs, "", "  ")
	if err != nil {
		log.Printf("--ERROR-- with SaveMetadataFieldSpecs(), couldn't write data in metadata_fields_specs.json: %v\n", err)
		return
	}
	err = os.WriteFile(localStorageMetadataFieldsSpecsPath, jsonData, 0644)
	if err != nil {
		log.Printf("--ERROR-- with SaveMetadataFieldSpecs(), couldn't write data in metadata_fields_specs.json: %v\n", err)
		return
	}

}

func (c *AppContext) LoadMetadataFieldsMap() {
	var fieldsMap map[string][]string

	data, err := os.ReadFile(localStorageMetadataFieldsMapPath)
	if err != nil {
		log.Printf("--ERROR-- with LoadMetadataFieldsMap(), couldn't json.Marshal data for metadata_fields_specs.json: %v\n", err)
		return
	}

	err = json.Unmarshal(data, &fieldsMap)
	if err != nil {
		log.Printf("--ERROR-- with LoadMetadataFieldsMap(), couldn't write data in metadata_fields_specs.json: %v\n", err)
		return
	}

	c.MetadataFieldsMap = fieldsMap
}

func (c *AppContext) SaveMedataFieldsMap() {

	jsonData, err := json.MarshalIndent(c.MetadataFieldsMap, "", "  ")
	if err != nil {
		log.Printf("--ERROR-- with SaveMedataFieldsMap(), couldn't write data in metadata_fields_specs.json: %v\n", err)
		return
	}
	err = os.WriteFile(localStorageMetadataFieldsMapPath, jsonData, 0644)
	if err != nil {
		log.Printf("--ERROR-- with SaveMedataFieldsMap(), couldn't write data in metadata_fields_specs.json: %v\n", err)
		return
	}

}
