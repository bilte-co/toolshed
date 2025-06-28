package clock_test

import (
	"testing"
	"time"

	"github.com/bilte-co/toolshed/clock"
)

func TestParseTimeFromHHMM(t *testing.T) {
	baseDate := time.Date(2025, 3, 7, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		input    string
		expected *time.Time
	}{
		{"1040", &time.Time{}},
		{"940", &time.Time{}},
		{"2359", &time.Time{}},
		{"0000", &time.Time{}},
		{"", nil},
		{"invalid", nil},
		{"2400", nil}, // Invalid hour
	}

	for _, tt := range tests {
		result := clock.ParseTimeFromHHMM(baseDate, tt.input)
		if tt.expected == nil && result != nil {
			t.Errorf("parseTimeFromHHMM(%q) = %v, want nil", tt.input, result)
		}
		if tt.expected != nil && result == nil {
			t.Errorf("parseTimeFromHHMM(%q) = nil, want non-nil", tt.input)
		}
		if tt.expected != nil && result != nil {
			// Check specific cases
			if tt.input == "1040" {
				expected := time.Date(2025, 3, 7, 10, 40, 0, 0, time.UTC)
				if !result.Equal(expected) {
					t.Errorf("parseTimeFromHHMM(%q) = %v, want %v", tt.input, result, expected)
				}
			}
			if tt.input == "940" {
				expected := time.Date(2025, 3, 7, 9, 40, 0, 0, time.UTC)
				if !result.Equal(expected) {
					t.Errorf("parseTimeFromHHMM(%q) = %v, want %v", tt.input, result, expected)
				}
			}
		}
	}
}
