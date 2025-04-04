package kallaxyapi

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/VincNT21/kallaxy/client/models"
)

func (c *RecordsClient) CreateRecord(mediumID, startDate, endDate string) (models.Record, error) {
	type parametersCreateUserMediumRecord struct {
		MediumID  string `json:"medium_id"`
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
	}

	// Convert input data to match server's requirement
	if startDate != "" {
		parsedStartDate, err := time.Parse("2006/01/02", startDate)
		if err != nil {
			return models.Record{}, err
		}
		startDate = parsedStartDate.Format(time.RFC3339)
	}

	if endDate != "" {
		parsedEndDate, err := time.Parse("2006/01/02", endDate)
		if err != nil {
			return models.Record{}, err
		}
		endDate = parsedEndDate.Format(time.RFC3339)
	}

	// Parameters for request
	params := parametersCreateUserMediumRecord{
		MediumID:  mediumID,
		StartDate: startDate,
		EndDate:   endDate,
	}

	// Make request
	r, err := c.apiClient.makeHttpRequest(c.apiClient.Config.Endpoints.Records.CreateRecord, params)
	if err != nil {
		log.Printf("--ERROR-- with CreateRecord(): %v\n", err)
		return models.Record{}, err
	}
	defer r.Body.Close()

	// Check response's status code
	if r.StatusCode != 201 {
		log.Printf("--ERROR-- with CreateRecord(). Response status code: %v\n", r.StatusCode)
		switch r.StatusCode {
		case 400:
			return models.Record{}, models.ErrBadRequest
		case 401:
			return models.Record{}, models.ErrUnauthorized
		case 404:
			return models.Record{}, models.ErrNotFound
		case 409:
			return models.Record{}, models.ErrConflict
		case 500:
			return models.Record{}, models.ErrServerIssue
		default:
			return models.Record{}, fmt.Errorf("unknown error status code: %v", r.StatusCode)
		}
	}

	// Decode response
	var record models.Record
	err = json.NewDecoder(r.Body).Decode(&record)
	if err != nil {
		log.Printf("--ERROR-- with CreateRecord(): %v\n", err)
		return models.Record{}, err
	}

	// Return data
	log.Println("--DEBUG-- CreateRecord() OK")
	return record, nil
}
