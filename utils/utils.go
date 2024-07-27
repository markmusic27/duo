package process

import (
	"fmt"
	"strconv"
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

func RemoveEmptyStrings(slice []string) []string {
	// Create a new slice to hold non-empty strings
	result := []string{}

	// Iterate over the original slice
	for _, str := range slice {
		if str != "" {
			// Append non-empty strings to the result slice
			result = append(result, str)
		}
	}

	return result
}

func RemoveEmptyLines(lines []string) []string {
	var filtered []string
	for _, line := range lines {
		// Remove all spaces, tabs, etc.
		trimmedLine := strings.ReplaceAll(line, " ", "")
		trimmedLine = strings.ReplaceAll(trimmedLine, "\t", "")

		// Filter out empty lines
		if trimmedLine != "" {
			filtered = append(filtered, line)
		}
	}
	return filtered
}

func IndexOf(slice []string, item string) int {
	for i, v := range slice {
		if v == item {
			return i
		}
	}
	return -1 // Item not found
}

func IsNumber(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}
