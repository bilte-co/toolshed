package clock_test

import (
	"testing"
	"time"

	"github.com/bilte-co/toolshed/clock"
	"github.com/stretchr/testify/require"
)

func TestParseTimeFromHHMM_ValidInputs(t *testing.T) {
	baseDate := time.Date(2025, 3, 7, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name     string
		input    string
		expected time.Time
	}{
		{
			name:     "4-digit time 10:40",
			input:    "1040",
			expected: time.Date(2025, 3, 7, 10, 40, 0, 0, time.UTC),
		},
		{
			name:     "3-digit time 9:40 (auto-padded)",
			input:    "940",
			expected: time.Date(2025, 3, 7, 9, 40, 0, 0, time.UTC),
		},
		{
			name:     "3-digit time 1:23 (auto-padded)",
			input:    "123",
			expected: time.Date(2025, 3, 7, 1, 23, 0, 0, time.UTC),
		},
		{
			name:     "midnight 00:00",
			input:    "0000",
			expected: time.Date(2025, 3, 7, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "3-digit midnight 0:00 (auto-padded)",
			input:    "000",
			expected: time.Date(2025, 3, 7, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "end of day 23:59",
			input:    "2359",
			expected: time.Date(2025, 3, 7, 23, 59, 0, 0, time.UTC),
		},
		{
			name:     "noon 12:00",
			input:    "1200",
			expected: time.Date(2025, 3, 7, 12, 0, 0, 0, time.UTC),
		},
		{
			name:     "input with leading/trailing spaces",
			input:    "  1545  ",
			expected: time.Date(2025, 3, 7, 15, 45, 0, 0, time.UTC),
		},
		{
			name:     "boundary hour 23",
			input:    "2300",
			expected: time.Date(2025, 3, 7, 23, 0, 0, 0, time.UTC),
		},
		{
			name:     "boundary minute 59",
			input:    "1259",
			expected: time.Date(2025, 3, 7, 12, 59, 0, 0, time.UTC),
		},
		{
			name:     "3-digit time 2:40 AM (auto-padded)",
			input:    "240",
			expected: time.Date(2025, 3, 7, 2, 40, 0, 0, time.UTC),
		},
		{
			name:     "negative hour normalization -1:23 -> 23:23 previous day",
			input:    "-123",
			expected: time.Date(2025, 3, 6, 23, 23, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := clock.ParseTimeFromHHMM(baseDate, tt.input)
			require.NotNil(t, result, "expected non-nil result for input %q", tt.input)
			require.True(t, result.Equal(tt.expected), "expected %v, got %v", tt.expected, *result)
		})
	}
}

func TestParseTimeFromHHMM_InvalidInputs(t *testing.T) {
	baseDate := time.Date(2025, 3, 7, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "empty string",
			input: "",
		},
		{
			name:  "only spaces",
			input: "   ",
		},
		{
			name:  "too short (1 digit)",
			input: "1",
		},
		{
			name:  "too short (2 digits)",
			input: "12",
		},
		{
			name:  "too long (5 digits)",
			input: "12345",
		},
		{
			name:  "too long (6 digits)",
			input: "123456",
		},
		{
			name:  "non-numeric input",
			input: "abcd",
		},
		{
			name:  "mixed alphanumeric",
			input: "12ab",
		},
		{
			name:  "mixed alphanumeric 3-digit",
			input: "1ab",
		},
		{
			name:  "special characters",
			input: "12:34",
		},

		{
			name:  "4-digit negative number",
			input: "-999",
		},
		{
			name:  "invalid hour 24",
			input: "2400",
		},
		{
			name:  "invalid hour 25",
			input: "2500",
		},
		{
			name:  "invalid hour 99",
			input: "9900",
		},
		{
			name:  "invalid minute 60",
			input: "1260",
		},
		{
			name:  "invalid minute 99",
			input: "1299",
		},
		{
			name:  "both hour and minute invalid",
			input: "2460",
		},

		{
			name:  "3-digit invalid minute",
			input: "160",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := clock.ParseTimeFromHHMM(baseDate, tt.input)
			require.Nil(t, result, "expected nil result for invalid input %q", tt.input)
		})
	}
}

func TestParseTimeFromHHMM_DifferentBaseDates(t *testing.T) {
	tests := []struct {
		name     string
		baseDate time.Time
		input    string
		expected time.Time
	}{
		{
			name:     "different year",
			baseDate: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			input:    "1530",
			expected: time.Date(2020, 1, 1, 15, 30, 0, 0, time.UTC),
		},
		{
			name:     "leap year date",
			baseDate: time.Date(2020, 2, 29, 0, 0, 0, 0, time.UTC),
			input:    "0815",
			expected: time.Date(2020, 2, 29, 8, 15, 0, 0, time.UTC),
		},
		{
			name:     "end of year",
			baseDate: time.Date(2023, 12, 31, 0, 0, 0, 0, time.UTC),
			input:    "2359",
			expected: time.Date(2023, 12, 31, 23, 59, 0, 0, time.UTC),
		},
		{
			name:     "base date with existing time components",
			baseDate: time.Date(2025, 6, 15, 14, 30, 45, 123456789, time.UTC),
			input:    "0900",
			expected: time.Date(2025, 6, 15, 9, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := clock.ParseTimeFromHHMM(tt.baseDate, tt.input)
			require.NotNil(t, result, "expected non-nil result")
			require.True(t, result.Equal(tt.expected), "expected %v, got %v", tt.expected, *result)
		})
	}
}

func TestParseTimeFromHHMM_EdgeCases(t *testing.T) {
	baseDate := time.Date(2025, 3, 7, 0, 0, 0, 0, time.UTC)

	t.Run("all zeros padded", func(t *testing.T) {
		result := clock.ParseTimeFromHHMM(baseDate, "000")
		require.NotNil(t, result)
		expected := time.Date(2025, 3, 7, 0, 0, 0, 0, time.UTC)
		require.True(t, result.Equal(expected))
	})

	t.Run("single digit hour", func(t *testing.T) {
		result := clock.ParseTimeFromHHMM(baseDate, "100")
		require.NotNil(t, result)
		expected := time.Date(2025, 3, 7, 1, 0, 0, 0, time.UTC)
		require.True(t, result.Equal(expected))
	})

	t.Run("4-digit valid time 10:01", func(t *testing.T) {
		result := clock.ParseTimeFromHHMM(baseDate, "1001")
		require.NotNil(t, result)
		expected := time.Date(2025, 3, 7, 10, 1, 0, 0, time.UTC)
		require.True(t, result.Equal(expected))
	})
}
