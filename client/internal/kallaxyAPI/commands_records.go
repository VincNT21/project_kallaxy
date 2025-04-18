package kallaxyapi

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/VincNT21/kallaxy/client/models"
)

func (c *RecordsClient) CreateRecord(mediumID, startDate, endDate, comments string) (models.Record, error) {
	type parametersCreateUserMediumRecord struct {
		MediumID  string `json:"medium_id"`
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
		Comments  string `json:"comments"`
	}

	// Convert input data to match server's requirement
	if startDate != "" {
		parsedStartDate, err := c.apiClient.Helpers.FormatDateToServerFormat(startDate)
		if err != nil {
			return models.Record{}, err
		}
		startDate = parsedStartDate
	}

	if endDate != "" {
		parsedEndDate, err := c.apiClient.Helpers.FormatDateToServerFormat(endDate)
		if err != nil {
			return models.Record{}, err
		}
		endDate = parsedEndDate
	}

	// Parameters for request
	params := parametersCreateUserMediumRecord{
		MediumID:  mediumID,
		StartDate: startDate,
		EndDate:   endDate,
		Comments:  comments,
	}

	// Make request
	r, err := c.apiClient.makeHttpRequest(c.apiClient.Config.Endpoints.Records.CreateRecord, params)
	if err != nil {
		log.Printf("--ERROR-- with CreateRecord(): %v\n", err)
		return models.Record{}, err
	}
	defer r.Body.Close()

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

func (c *RecordsClient) UpdateRecord(recordID, startDate, endDate, comments string) (models.Record, error) {
	type parametersUpdateRecord struct {
		RecordID  string `json:"record_id"`
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
		Comments  string `json:"comments"`
	}

	// Convert input data to match server's requirement
	if startDate != "" {
		parsedStartDate, err := c.apiClient.Helpers.FormatDateToServerFormat(startDate)
		if err != nil {
			return models.Record{}, err
		}
		startDate = parsedStartDate
	}

	if endDate != "" {
		parsedEndDate, err := c.apiClient.Helpers.FormatDateToServerFormat(endDate)
		if err != nil {
			return models.Record{}, err
		}
		endDate = parsedEndDate
	}

	params := parametersUpdateRecord{
		RecordID:  recordID,
		StartDate: startDate,
		EndDate:   endDate,
		Comments:  comments,
	}

	// Make request
	r, err := c.apiClient.makeHttpRequest(c.apiClient.Config.Endpoints.Records.UpdateRecord, params)
	if err != nil {
		log.Printf("--ERROR-- with UpdateRecord(): %v\n", err)
		return models.Record{}, err
	}
	defer r.Body.Close()

	// Decode response
	var record models.Record
	err = json.NewDecoder(r.Body).Decode(&record)
	if err != nil {
		log.Printf("--ERROR-- with UpdateRecord(): %v\n", err)
		return models.Record{}, err
	}

	// Return data
	log.Println("--DEBUG-- UpdateRecord() OK")
	return record, nil
}

func (c *RecordsClient) DeleteRecord(mediumID string) error {
	type parametersDeleteRecord struct {
		MediumID string `json:"medium_id"`
	}

	params := parametersDeleteRecord{
		MediumID: mediumID,
	}

	// Make request
	r, err := c.apiClient.makeHttpRequest(c.apiClient.Config.Endpoints.Records.DeleteRecord, params)
	if err != nil {
		log.Printf("--ERROR-- with DeleteRecord(): %v\n", err)
		return err
	}
	defer r.Body.Close()

	// Check response's status code
	if r.StatusCode != 200 {
		log.Printf("--ERROR-- with DeleteRecord(). Response status code: %v\n", r.StatusCode)
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
	log.Println("--DEBUG-- DeleteRecord() OK")
	return nil
}
