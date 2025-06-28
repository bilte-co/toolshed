// Package csv provides utilities for CSV file processing and validation.
// It includes functions for validating CSV headers against expected schemas
// and other common CSV manipulation tasks.
//
// Example usage:
//
//	expectedHeaders := []string{"Name", "Email", "Age"}
//	actualHeaders := []string{"Name", "Email", "Age", "Phone"}
//
//	// Check if all expected headers are present
//	if csv.ValidateHeaders(expectedHeaders, actualHeaders) {
//		fmt.Println("All required headers are present")
//	}
package csv

import (
	"slices"
)

// ValidateHeaders checks if all expected headers are present in the actual headers slice.
// It returns true if every header in expectedHeaders is found in headers, false otherwise.
// The order of headers does not matter, and headers can contain additional columns
// beyond those specified in expectedHeaders.
func ValidateHeaders(expectedHeaders []string, headers []string) bool {
	for _, expected := range expectedHeaders {
		if !slices.Contains(headers, expected) {
			return false
		}
	}
	return true
}
