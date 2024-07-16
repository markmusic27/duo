package process

func ClampString(input string, clamp int) string {
	if len(input) > clamp {
		return input[:clamp]
	}
	return input
}
