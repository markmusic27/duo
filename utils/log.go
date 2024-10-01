package process

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type LogPayload struct {
	Project     string `json:"project"`
	Channel     string `json:"channel"`
	Event       string `json:"event"`
	Description string `json:"description,omitempty"`
	UserID      string `json:"user_id"`
	Icon        string `json:"icon,omitempty"`
	Notfiy      bool   `json:"notify"`
}

func Log(emoji, message, description string) error {

	apiToken := os.Getenv("LOGSNAG")

	url := "https://api.logsnag.com/v1/log"
	payload := LogPayload{
		Project:     "duo",
		Channel:     "logs",
		Event:       message,
		Description: description,
		Icon:        emoji,
		UserID:      "mark",
		Notfiy:      true,
	}

	// Convert the payload to JSON
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %v", err)
	}

	// Create the HTTP request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiToken)

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}

	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %v", err)
	}
	log.Println(string(bodyBytes))

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received non-OK response: %s", resp.Status)
	}

	fmt.Println("Log sent successfully!")
	return nil
}
