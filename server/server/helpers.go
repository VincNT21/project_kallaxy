package server

import (
	"fmt"

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
		return interval, fmt.Errorf("End date: %v is after Start date: %v", endDate.Time, startDate.Time)
	}

	// Calculate the duration in days
	duration := endDate.Time.Sub(startDate.Time)
	days := int32(duration.Hours() / 24)

	// Set up the interval
	interval.Valid = true
	interval.Days = days
	interval.Months = 0
	interval.Microseconds = 0

	return interval, nil
}
