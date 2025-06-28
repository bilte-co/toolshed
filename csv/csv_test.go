package csv_test

import (
	"testing"

	"github.com/bilte-co/toolshed/csv"
)

func TestValidateHeaders(t *testing.T) {
	expectedHeaders := []string{
		"FlightDate", "Reporting_Airline", "DOT_ID_Reporting_Airline", "IATA_CODE_Reporting_Airline", "Tail_Number", "Flight_Number_Reporting_Airline",
	}

	headers := []string{
		"Year", "Quarter", "Month", "DayofMonth", "DayOfWeek", "FlightDate", "Reporting_Airline", "DOT_ID_Reporting_Airline", "IATA_CODE_Reporting_Airline",
		"Tail_Number", "Flight_Number_Reporting_Airline",
	}

	if !csv.ValidateHeaders(expectedHeaders, headers) {
		t.Error("validateHeaders should return true for valid BTS headers")
	}

	invalidHeaders := []string{"Col1", "Col2", "Col3"}
	if csv.ValidateHeaders(expectedHeaders, invalidHeaders) {
		t.Error("validateHeaders should return false for invalid headers")
	}
}
