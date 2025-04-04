package kallaxyapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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

	// Check response's status code
	if r.StatusCode != 201 {
		log.Printf("--ERROR-- with CreateMedium(). Response status code: %v\n", r.StatusCode)
		switch r.StatusCode {
		case 400:
			return models.Medium{}, models.ErrBadRequest
		case 401:
			return models.Medium{}, models.ErrUnauthorized
		case 409:
			return models.Medium{}, models.ErrConflict
		case 500:
			return models.Medium{}, models.ErrServerIssue
		default:
			return models.Medium{}, fmt.Errorf("unknown error status code: %v", r.StatusCode)
		}
	}

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

func (c *MediaClient) GetImageUrl(mediumType, mediumTitleOrId string) (string, error) {
	switch mediumType {
	case "book":
		return fmt.Sprintf("https://covers.openlibrary.org/b/isbn/%s-M.jpg", mediumTitleOrId), nil
	case "movie":
		movie, err := c.apiClient.External.SearchForMovie(mediumTitleOrId)
		if err != nil {
			return "", err
		}
		if len(movie.Results) == 0 {
			return "", err
		}
		return fmt.Sprintf("https://image.tmdb.org/t/p/w200%s", movie.Results[0].PosterPath), nil
	case "series":
		series, err := c.apiClient.External.SearchForSeries(mediumTitleOrId)
		if err != nil {
			return "", err
		}
		if len(series.Results) == 0 {
			return "", err
		}
		return fmt.Sprintf("https://image.tmdb.org/t/p/w200%s", series.Results[0].PosterPath), nil
	case "videogame":
		videogame, err := c.apiClient.External.SearchForVideogameOnSteam(mediumTitleOrId)
		if err != nil {
			return "", err
		}
		if len(videogame.Results) == 0 {
			return "", err
		}
		return videogame.Results[0].BackgroundImage, err
	case "boardgame":
		boardgame, err := c.apiClient.External.SearchForBoardgame(mediumTitleOrId)
		if err != nil {
			return "", err
		}
		return boardgame.Items.Item.Image, nil
	default:
		return "", nil
	}
}

func (c *MediaClient) FetchImage(imageUrl string) (*bytes.Buffer, error) {

	// Make request
	r, err := http.Get(imageUrl)
	if err != nil {
		log.Printf("--ERROR-- with FetchImage(): %v\n", err)
		return nil, err
	}
	defer r.Body.Close()

	// Check response's status code
	if r.StatusCode != 200 {
		log.Printf("--ERROR-- with FetchImage(). Response status code: %v\n", r.StatusCode)
		return nil, fmt.Errorf("problem with FetchImage() request, status code: %v", r.StatusCode)
	}

	// Load the response body into a buffer
	var buf bytes.Buffer
	_, err = buf.ReadFrom(r.Body)
	if err != nil {
		log.Printf("--ERROR-- with FetchImage(): %v\n", err)
		return nil, err
	}

	// Return data
	log.Println("--DEBUG-- FetchImage() OK")
	return &buf, nil
}

func (c *MediaClient) GetMediaWithRecords() (models.MediaWithRecords, error) {

	// Make request
	r, err := c.apiClient.makeHttpRequest(c.apiClient.Config.Endpoints.Media.GetMediaWithRecords, nil)
	if err != nil {
		log.Printf("--ERROR-- with GetMediaWithRecords(): %v\n", err)
		return models.MediaWithRecords{}, err
	}
	defer r.Body.Close()

	// Check response's status code
	if r.StatusCode != 200 {
		log.Printf("--ERROR-- with GetMediaWithRecords(). Response status code: %v\n", r.StatusCode)
		switch r.StatusCode {
		case 400:
			return models.MediaWithRecords{}, models.ErrBadRequest
		case 401:
			return models.MediaWithRecords{}, models.ErrUnauthorized
		case 404:
			return models.MediaWithRecords{}, models.ErrNotFound
		case 409:
			return models.MediaWithRecords{}, models.ErrConflict
		case 500:
			return models.MediaWithRecords{}, models.ErrServerIssue
		default:
			return models.MediaWithRecords{}, fmt.Errorf("unknown error status code: %v", r.StatusCode)
		}
	}

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

func (c *MediaClient) GetMediaTypes(mediaRecords models.MediaWithRecords) map[string]bool {

	mediaTypes := make(map[string]bool)

	// Iterate over MediaRecords to check for different media types
	for mediaType := range mediaRecords.MediaRecords {
		mediaTypes[mediaType] = true
	}

	// Return data
	log.Println("--DEBUG-- GetMediaTypes() OK")
	return mediaTypes
}
