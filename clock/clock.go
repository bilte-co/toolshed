// Package clock provides time parsing utilities for common time formats.
// It focuses on parsing condensed time formats like "1040" (10:40 AM) or "940" (9:40 AM)
// and combining them with a base date to create full timestamps.
//
// Example usage:
//
//	baseDate := time.Date(2023, 12, 25, 0, 0, 0, 0, time.UTC)
//	
//	// Parse "1040" as 10:40 AM on December 25, 2023
//	parsedTime := clock.ParseTimeFromHHMM(baseDate, "1040")
//	if parsedTime != nil {
//		fmt.Println(parsedTime.Format("2006-01-02 15:04:05"))
//	}
//
//	// Parse "940" as 9:40 AM (automatically pads to 0940)
//	parsedTime = clock.ParseTimeFromHHMM(baseDate, "940")
//	if parsedTime != nil {
//		fmt.Println(parsedTime.Format("2006-01-02 15:04:05"))
//	}
package clock

import (
	"strconv"
	"strings"
	"time"
)

// ParseTimeFromHHMM parses a time string in HHMM or HMM format and combines it with a base date.
// The function accepts 3 or 4 digit time strings (e.g., "940" or "1040") and returns
// a pointer to the resulting time in UTC, or nil if parsing fails.
//
// The time string is automatically padded to 4 digits if it's 3 digits long.
// Hours must be 00-23 and minutes must be 00-59.
func ParseTimeFromHHMM(baseDate time.Time, timeStr string) *time.Time {
	timeStr = strings.TrimSpace(timeStr)
	if timeStr == "" {
		return nil
	}

	// Pad to 4 digits if needed
	if len(timeStr) == 3 {
		timeStr = "0" + timeStr
	}
	if len(timeStr) != 4 {
		return nil
	}

	hour, err := strconv.Atoi(timeStr[:2])
	if err != nil {
		return nil
	}
	minute, err := strconv.Atoi(timeStr[2:])
	if err != nil {
		return nil
	}

	if hour > 23 || minute > 59 {
		return nil
	}

	result := time.Date(baseDate.Year(), baseDate.Month(), baseDate.Day(), hour, minute, 0, 0, time.UTC)
	return &result
}
