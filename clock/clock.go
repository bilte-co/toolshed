package clock

import (
	"strconv"
	"strings"
	"time"
)

func ParseTimeFromHHMM(baseDate time.Time, timeStr string) *time.Time {
	// Handle 3-4 digit time format (e.g., "1040" = 10:40, "940" = 9:40)
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
