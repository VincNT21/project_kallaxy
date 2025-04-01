package server

import (
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func calculateDuration(startDate, endDate pgtype.Timestamp) (pgtype.Interval, error) {
	var interval pgtype.Interval

	// Check if both dates are present, if not return a null interval
	if !startDate.Valid || !endDate.Valid {
		interval.Valid = false
		return interval, nil
	}

	// Check if startDate is before endDate
	if !startDate.Time.Before(endDate.Time) {
		return interval, fmt.Errorf("end date: %v is after Start date: %v", endDate.Time, startDate.Time)
	}

	// Calculate the duration in days
	duration := endDate.Time.Sub(startDate.Time)
	days := int32(duration.Hours() / 24)

	// Set up the interval
	interval.Valid = true
	interval.Days = days

	return interval, nil
}

func convertIdToPgtype(stringID string) (pgtype.UUID, error) {
	var id pgtype.UUID
	err := id.Scan(stringID)
	if err != nil {
		return id, errors.New("invalid string id format")
	}
	return id, nil
}

func convertDateToPgtype(stringdate string) (pgtype.Timestamp, error) {
	var date pgtype.Timestamp

	if stringdate == "" {
		return date, nil
	}

	formatDate, err := time.Parse(time.RFC3339, stringdate)
	if err != nil {
		return date, errors.New("invalid string date format")
	}

	err = date.Scan(formatDate)
	if err != nil {
		return date, errors.New("invalid string date format")
	}
	return date, nil
}
