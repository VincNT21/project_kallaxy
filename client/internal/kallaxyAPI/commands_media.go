package kallaxyapi

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/VincNT21/kallaxy/client/models"
)

func (c *MediaClient) CreateMediumAndRecord(title, mediaType, creator, pubDate, imageUrl, startDate, endDate, comments string, metadata map[string]interface{}) (models.Medium, models.Record, error) {

	// Make request for Medium creation
	medium, err := c.apiClient.Media.CreateMedium(title, mediaType, creator, pubDate, imageUrl, metadata)
	if err != nil {
		log.Printf("--ERROR-- with CreateMediumAndRecord(): %v\n", err)
		return models.Medium{}, models.Record{}, err
	}

	// Make request for Record creation
	record, err := c.apiClient.Records.CreateRecord(medium.ID, startDate, endDate, comments)
	if err != nil {
		log.Printf("--ERROR-- with CreateMediumAndRecord(): %v\n", err)
		return models.Medium{}, models.Record{}, err
	}

	// Get image from temp Cache and put it in local cache
	if imageUrl != "" {
		dataTemp, exists := c.apiClient.Cache.GetFromTemp(imageUrl)
		if exists {
			c.apiClient.Cache.Add(imageUrl, dataTemp)
		}
	}

	// Return data
	return medium, record, nil
}

func (c *MediaClient) CreateMedium(title, mediaType, creator, pubDate, imageUrl string, metadata map[string]interface{}) (models.Medium, error) {
	type parametersCreateMedium struct {
		Title     string                 `json:"title"`
		MediaType string                 `json:"media_type"`
		Creator   string                 `json:"creator"`
		PubDate   string                 `json:"pub_date"`
		ImageUrl  string                 `json:"image_url"`
		Metadata  map[string]interface{} `json:"metadata"`
	}

	// Parameters for Create Medium request
	params := parametersCreateMedium{
		Title:     title,
		MediaType: mediaType,
		Creator:   creator,
		PubDate:   pubDate,
		ImageUrl:  imageUrl,
		Metadata:  metadata,
	}

	// Make request for Medium creation
	r, err := c.apiClient.makeHttpRequest(c.apiClient.Config.Endpoints.Media.CreateMedia, params)
	if err != nil {
		log.Printf("--ERROR-- with CreateMedium(): %v\n", err)
		return models.Medium{}, err
	}
	defer r.Body.Close()

	// Decode response
	var medium models.Medium
	err = json.NewDecoder(r.Body).Decode(&medium)
	if err != nil {
		log.Printf("--ERROR-- with CreateMedium(): %v\n", err)
		return models.Medium{}, err
	}

	// Return data
	log.Println("--DEBUG-- CreateMedium() OK")
	return medium, nil
}

func (c *MediaClient) GetMediaWithRecords() (models.MediaWithRecords, error) {

	// Make request
	r, err := c.apiClient.makeHttpRequest(c.apiClient.Config.Endpoints.Media.GetMediaWithRecords, nil)
	if err != nil {
		log.Printf("--ERROR-- with GetMediaWithRecords(): %v\n", err)
		return models.MediaWithRecords{}, err
	}
	defer r.Body.Close()

	// Decode response
	var mediaRecords models.MediaWithRecords
	err = json.NewDecoder(r.Body).Decode(&mediaRecords)
	if err != nil {
		log.Printf("--ERROR-- with GetMediaWithRecords(): %v\n", err)
		return models.MediaWithRecords{}, err
	}

	// Return data
	log.Println("--DEBUG-- GetMediaWithRecords() OK")
	return mediaRecords, nil
}

func (c *MediaClient) UpdateMedium(mediumId, title, creator, pubDate, imageUrl string, metadata map[string]interface{}) (models.ClientMedium, error) {
	type parametersUpdateMedium struct {
		MediumID string                 `json:"medium_id"`
		Title    string                 `json:"title"`
		Creator  string                 `json:"creator"`
		PubDate  string                 `json:"pub_date"`
		ImageUrl string                 `json:"image_url"`
		Metadata map[string]interface{} `json:"metadata"`
	}

	params := parametersUpdateMedium{
		MediumID: mediumId,
		Title:    title,
		Creator:  creator,
		PubDate:  pubDate,
		ImageUrl: imageUrl,
		Metadata: metadata,
	}

	// Make request
	r, err := c.apiClient.makeHttpRequest(c.apiClient.Config.Endpoints.Media.UpdateMedia, params)
	if err != nil {
		log.Printf("--ERROR-- with UpdateMedium(): %v\n", err)
		return models.ClientMedium{}, err
	}
	defer r.Body.Close()

	// Check response's status code
	if r.StatusCode != 200 {
		log.Printf("--ERROR-- with UpdateMedium(). Response status code: %v\n", r.StatusCode)
		switch r.StatusCode {
		case 400:
			return models.ClientMedium{}, models.ErrBadRequest
		case 401:
			return models.ClientMedium{}, models.ErrUnauthorized
		case 404:
			return models.ClientMedium{}, models.ErrNotFound
		case 409:
			return models.ClientMedium{}, models.ErrConflict
		case 500:
			return models.ClientMedium{}, models.ErrServerIssue
		default:
			return models.ClientMedium{}, fmt.Errorf("unknown error status code: %v", r.StatusCode)
		}
	}

	// Decode response
	var updatedMedium models.ClientMedium
	err = json.NewDecoder(r.Body).Decode(&updatedMedium)
	if err != nil {
		log.Printf("--ERROR-- with UpdateMedium(): %v\n", err)
		return models.ClientMedium{}, err
	}

	// Return data
	log.Println("--DEBUG-- UpdateMedium() OK")
	return updatedMedium, nil
}

func (c *MediaClient) DeleteMedium(mediumID string) error {
	type parametersDeleteMedium struct {
		MediumID string `json:"medium_id"`
	}

	params := parametersDeleteMedium{
		MediumID: mediumID,
	}

	// Make request
	r, err := c.apiClient.makeHttpRequest(c.apiClient.Config.Endpoints.Media.DeleteMedia, params)
	if err != nil {
		log.Printf("--ERROR-- with DeleteMedium(): %v\n", err)
		return err
	}
	defer r.Body.Close()

	// Check response's status code
	if r.StatusCode != 200 {
		log.Printf("--ERROR-- with DeleteMedium(). Response status code: %v\n", r.StatusCode)
		switch r.StatusCode {
		case 400:
			return models.ErrBadRequest
		case 401:
			return models.ErrUnauthorized
		case 404:
			return models.ErrNotFound
		case 409:
			return models.ErrConflict
		case 500:
			return models.ErrServerIssue
		default:
			return fmt.Errorf("unknown error status code: %v", r.StatusCode)
		}
	}

	// Return
	log.Println("--DEBUG-- DeleteMedium() OK")
	return nil
}
