package server

import (
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func TestCalculateDuration(t *testing.T) {
	// Create some Timestamp date
	var startDate pgtype.Timestamp
	startTime := time.Now().UTC()
	startDate.Valid = true
	startDate.Time = startTime

	var endDateValid pgtype.Timestamp
	endTimeValid := time.Now().UTC().AddDate(0, 0, 21)
	endDateValid.Valid = true
	endDateValid.Time = endTimeValid

	var endDateMonths pgtype.Timestamp
	endTimeMonths := time.Now().UTC().AddDate(0, 2, 0)
	endDateMonths.Valid = true
	endDateMonths.Time = endTimeMonths

	var endDatePast pgtype.Timestamp
	endTimePast := time.Now().UTC().AddDate(0, 0, -10)
	endDatePast.Valid = true
	endDatePast.Time = endTimePast

	var endDateInvalid pgtype.Timestamp
	endDateInvalid.Valid = false

	// Create tests table
	tests := []struct {
		name           string
		startDate      pgtype.Timestamp
		endDate        pgtype.Timestamp
		wantErr        bool
		wantValid      bool
		wantDaysResult int32
	}{
		{
			name:           "Valid",
			startDate:      startDate,
			endDate:        endDateValid,
			wantErr:        false,
			wantValid:      true,
			wantDaysResult: 21,
		},
		{
			name:           "Valid, with month",
			startDate:      startDate,
			endDate:        endDateMonths,
			wantErr:        false,
			wantValid:      true,
			wantDaysResult: 61,
		},
		{
			name:      "End date not set",
			startDate: startDate,
			endDate:   endDateInvalid,
			wantErr:   false,
			wantValid: false,
		},
		{
			name:      "End date in past",
			startDate: startDate,
			endDate:   endDatePast,
			wantErr:   true,
			wantValid: false,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			interval, err := calculateDuration(tc.startDate, tc.endDate)
			if (err != nil) != tc.wantErr {
				t.Errorf("calculateDuration() error = %v, wantErr = %v", err, tc.wantErr)
			}
			if tc.wantValid {
				if !interval.Valid {
					t.Error("interval should be valid")
				}
				if interval.Days != tc.wantDaysResult {
					t.Errorf("calculateDuration() interval.Days = %v, wantDaysResult = %v", interval.Days, tc.wantDaysResult)
				}
			}
		})
	}
}
