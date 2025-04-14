package kallaxyapi

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/VincNT21/kallaxy/client/models"
)

func (c *AdminClient) CheckForClientUpdate(clientVersion string) (isNewVersion bool, versionTage, versionDescription, dlUrl string, err error) {

	url := "https://api.github.com/repos/VincNT21/project_kallaxy/releases/latest"

	// Make GET request
	r, err := http.Get(url)
	if err != nil {
		log.Printf("--ERROR-- with CheckForClientUpdate(): %v\n", err)
		return false, "", "", "", err
	}
	defer r.Body.Close()

	// Decode response
	var response models.ResponseGitHubApiRelease
	err = json.NewDecoder(r.Body).Decode(&response)
	if err != nil {
		log.Printf("--ERROR-- with CheckForClientUpdate(): %v\n", err)
		return false, "", "", "", err
	}

	// Check if current version is outdated
	if response.TagName != c.apiClient.ClientVersion {
		log.Println("--DEBUG-- CheckForClientUpdate() OK, found a new version")
		if len(response.Assets) < 1 {
			return false, "", "", "", errors.New("problem with GitHubApi response format")
		}
		return true, response.TagName, response.Body, response.Assets[0].BrowserDownloadURL, nil
	}

	// Else
	log.Println("--DEBUG-- CheckForClientUpdate() OK, no new version found")
	return false, "", "", "", nil
}

func (c *AdminClient) CheckForServerVersion() (string, error) {
	type responseServerVersion struct {
		ServerVersion string `json:"server_version"`
	}

	url := c.apiClient.Config.BaseURL + "/server/version"

	r, err := http.Get(url)
	if err != nil {
		log.Printf("--ERROR-- with CheckForServerVersion(): %v\n", err)
		return "", err
	}
	defer r.Body.Close()

	var response responseServerVersion
	err = json.NewDecoder(r.Body).Decode(&response)
	if err != nil {
		log.Printf("--ERROR-- with CheckForServerVersion(): %v\n", err)
		return "", err
	}

	return response.ServerVersion, nil
}
