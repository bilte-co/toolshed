package csv

import (
	"slices"
)

func ValidateHeaders(expectedHeaders []string, headers []string) bool {
	for _, expected := range expectedHeaders {
		if !slices.Contains(headers, expected) {
			return false
		}
	}
	return true
}
