package process

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type AITimezoneResponse struct {
	Timezone *string `json:"timezone,omitempty"`
	Error    *string `json:"error,omitempty"`
}

func SetTimezoneFromLocation(location string) (string, error) {
	res, err := Prompt(location, TimezoneTemplate)
	if err != nil {
		return "", err
	}

	var gen AITimezoneResponse
	err = json.Unmarshal([]byte(CleanCode(res)), &gen)
	if err != nil {
		return "", err
	}

	iana, err := ValidateIANATimezone(*gen.Timezone)
	if err != nil {
		return "", err
	}

	err = os.Setenv("LOCATION", iana)
	if err != nil {
		return "", err
	}

	return iana, nil
}

func ValidateIANATimezone(iana string) (string, error) {
	_, err := time.LoadLocation(iana)
	if err != nil {
		return "", fmt.Errorf("invalid IANA timezone string: %s", iana)
	}
	return iana, nil
}
