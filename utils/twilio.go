package process

import (
	"bytes"
	"errors"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func Message(phone string, message string) error {
	// Create the data for the POST request
	values := url.Values{}
	values.Set("From", os.Getenv("TWILIO_PHONE"))
	values.Set("Body", message)
	values.Set("To", phone)

	// Encode the data to be sent
	encodedValues := values.Encode()

	// Prepare the request
	req, err := http.NewRequest("POST", "https://api.twilio.com/2010-04-01/Accounts/"+os.Getenv("TWILIO_ACCOUNT_SID")+"/Messages.json", strings.NewReader(encodedValues))
	if err != nil {
		return err
	}

	// Set headers
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Set basic authentication
	req.SetBasicAuth(os.Getenv("TWILIO_ACCOUNT_SID"), os.Getenv("TWILIO_AUTH_TOKEN"))

	// Create an HTTP client and send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if !(resp.StatusCode >= 200 && resp.StatusCode < 300) {
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		return errors.New("non-successful response status " + resp.Status + " - response body: " + buf.String())
	}

	return nil
}

func TruncateString(input string) string {
	maxLength := 100
	if len(input) > maxLength {
		return input[:maxLength]
	}
	return input
}
