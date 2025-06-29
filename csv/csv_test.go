package csv_test

import (
	"testing"

	"github.com/bilte-co/toolshed/csv"
	"github.com/stretchr/testify/require"
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

func TestValidateHeaders_EmptyInputs(t *testing.T) {
	// Empty expected headers should return true (no requirements to check)
	result := csv.ValidateHeaders([]string{}, []string{"Col1", "Col2"})
	require.True(t, result)

	// Empty actual headers with expected headers should return false
	result = csv.ValidateHeaders([]string{"Col1"}, []string{})
	require.False(t, result)

	// Both empty should return true
	result = csv.ValidateHeaders([]string{}, []string{})
	require.True(t, result)
}

func TestValidateHeaders_NilInputs(t *testing.T) {
	// Nil expected headers should return true (no requirements to check)
	result := csv.ValidateHeaders(nil, []string{"Col1", "Col2"})
	require.True(t, result)

	// Nil actual headers with expected headers should return false
	result = csv.ValidateHeaders([]string{"Col1"}, nil)
	require.False(t, result)

	// Both nil should return true
	result = csv.ValidateHeaders(nil, nil)
	require.True(t, result)
}

func TestValidateHeaders_ExactMatch(t *testing.T) {
	headers := []string{"Name", "Email", "Age"}
	
	// Exact match in same order
	result := csv.ValidateHeaders(headers, headers)
	require.True(t, result)

	// Exact match in different order
	result = csv.ValidateHeaders([]string{"Name", "Email", "Age"}, []string{"Age", "Name", "Email"})
	require.True(t, result)
}

func TestValidateHeaders_SingleHeader(t *testing.T) {
	// Single expected header present
	result := csv.ValidateHeaders([]string{"Name"}, []string{"Name", "Email", "Age"})
	require.True(t, result)

	// Single expected header missing
	result = csv.ValidateHeaders([]string{"Phone"}, []string{"Name", "Email", "Age"})
	require.False(t, result)

	// Single expected and single actual, match
	result = csv.ValidateHeaders([]string{"Name"}, []string{"Name"})
	require.True(t, result)

	// Single expected and single actual, no match
	result = csv.ValidateHeaders([]string{"Name"}, []string{"Email"})
	require.False(t, result)
}

func TestValidateHeaders_CaseSensitive(t *testing.T) {
	// Case sensitivity test
	result := csv.ValidateHeaders([]string{"Name"}, []string{"name"})
	require.False(t, result)

	result = csv.ValidateHeaders([]string{"NAME"}, []string{"Name"})
	require.False(t, result)
}

func TestValidateHeaders_DuplicateHeaders(t *testing.T) {
	// Duplicate headers in expected should still work
	result := csv.ValidateHeaders([]string{"Name", "Name"}, []string{"Name", "Email"})
	require.True(t, result)

	// Duplicate headers in actual should still work
	result = csv.ValidateHeaders([]string{"Name"}, []string{"Name", "Name", "Email"})
	require.True(t, result)

	// Missing header despite duplicates in actual
	result = csv.ValidateHeaders([]string{"Name", "Phone"}, []string{"Name", "Name", "Email"})
	require.False(t, result)
}

func TestValidateHeaders_WhitespaceHeaders(t *testing.T) {
	// Headers with whitespace should be treated as different
	result := csv.ValidateHeaders([]string{"Name"}, []string{" Name"})
	require.False(t, result)

	result = csv.ValidateHeaders([]string{"Name "}, []string{"Name"})
	require.False(t, result)

	// Empty string header
	result = csv.ValidateHeaders([]string{""}, []string{"", "Name"})
	require.True(t, result)

	result = csv.ValidateHeaders([]string{""}, []string{"Name"})
	require.False(t, result)
}

func TestValidateHeaders_PartialMatch(t *testing.T) {
	expected := []string{"Name", "Email", "Age", "Phone"}
	actual := []string{"Name", "Email", "City"}

	// Should return false when some expected headers are missing
	result := csv.ValidateHeaders(expected, actual)
	require.False(t, result)
}

func TestValidateHeaders_SupersetHeaders(t *testing.T) {
	expected := []string{"Name", "Email"}
	actual := []string{"Name", "Email", "Age", "Phone", "City", "Country"}

	// Should return true when actual headers contain all expected headers plus more
	result := csv.ValidateHeaders(expected, actual)
	require.True(t, result)
}
