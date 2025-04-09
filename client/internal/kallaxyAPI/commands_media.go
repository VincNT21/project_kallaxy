package kallaxyapi

import (
	"encoding/json"
	"log"
	"strconv"

	"github.com/VincNT21/kallaxy/client/models"
)

func (c *MediaClient) CreateMediumAndRecord(title, mediaType, creator, releaseYear, imageUrl, startDate, endDate string) (models.Medium, models.Record, error) {

	// Make request for Medium creation in a goroutine
	medium, err := c.apiClient.Media.CreateMedium(title, mediaType, creator, releaseYear, imageUrl)
	if err != nil {
		log.Printf("--ERROR-- with CreateMediumAndRecord(): %v\n", err)
		return models.Medium{}, models.Record{}, err
	}

	// Make request for Record creation
	record, err := c.apiClient.Records.CreateRecord(medium.ID, startDate, endDate)
	if err != nil {
		log.Printf("--ERROR-- with CreateMediumAndRecord(): %v\n", err)
		return models.Medium{}, models.Record{}, err
	}

	// Return data
	return medium, record, nil
}

func (c *MediaClient) CreateMedium(title, mediaType, creator, releaseYear, imageUrl string) (models.Medium, error) {
	type parametersCreateMedium struct {
		Title       string          `json:"title"`
		MediaType   string          `json:"media_type"`
		Creator     string          `json:"creator"`
		ReleaseYear int32           `json:"release_year"`
		ImageUrl    string          `json:"image_url"`
		Metadata    json.RawMessage `json:"metadata"`
	}

	// Convert input data to match server's requirement
	releaseYearInt, err := strconv.Atoi(releaseYear)
	if err != nil {
		return models.Medium{}, nil
	}

	// Parameters for Create Medium request
	params := parametersCreateMedium{
		Title:       title,
		MediaType:   mediaType,
		Creator:     creator,
		ReleaseYear: int32(releaseYearInt),
		ImageUrl:    imageUrl,
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
