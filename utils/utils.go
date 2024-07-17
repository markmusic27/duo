package process

import "strings"

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
