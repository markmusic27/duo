package process

import (
	"fmt"
	"strings"
	"time"
)

func ClampString(input string, clamp int) string {
	if len(input) > clamp {
		return input[:clamp]
	}
	return input
}

func RemoveNewline(input string) string {
	// Replace \n with an empty string
	result := strings.ReplaceAll(input, "\n", "")
	// Replace \r with an empty string
	result = strings.ReplaceAll(result, "\r", "")
	return result
}

func RemoveIndex(slice []string, s int) []string {
	return append(slice[:s], slice[s+1:]...)
}

func ExtractTimezone(dateStr string) (string, error) {
	// Parse the datetime part
	t, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		return "", fmt.Errorf("error parsing date: %v", err)
	}

	// Get the timezone offset in seconds
	_, offsetSeconds := t.Zone()

	// Convert the offset to hours and minutes
	offsetHours := offsetSeconds / 3600
	offsetMinutes := (offsetSeconds % 3600) / 60

	// Format the timezone
	timezone := fmt.Sprintf("%+03d:%02d", offsetHours, offsetMinutes)

	// Determine if the timezone is behind or ahead of UTC
	var direction string
	if offsetHours == 0 && offsetMinutes == 0 {
		direction = "UTC"
		timezone = "Z"
	} else if offsetHours < 0 {
		direction = fmt.Sprintf("%d hours behind UTC", -offsetHours)
	} else {
		direction = fmt.Sprintf("%d hours ahead of UTC", offsetHours)
	}

	// Combine the formatted timezone and the direction
	timezoneWithDirection := fmt.Sprintf("%s (i.e., %s)", timezone, direction)

	return timezoneWithDirection, nil
}
