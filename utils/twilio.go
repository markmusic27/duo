package process

func TruncateString(input string) string {
	maxLength := 100
	if len(input) > maxLength {
		return input[:maxLength]
	}
	return input
}
