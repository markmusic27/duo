package process

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
)

// ⬇️ OpenAI Types
type OAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OAIRequestBody struct {
	Model    string       `json:"model"`
	Messages []OAIMessage `json:"messages"`
}

type ResponseChoice struct {
	Message OAIMessage `json:"message"`
}

type OAIResponseBody struct {
	Choices []ResponseChoice `json:"choices"`
}

// ⬇️ OpenAI Calls

const OpenAIEndpoint = "https://api.openai.com/v1/chat/completions"

func Prompt(user string, system string) (string, error) {
	requestBody := OAIRequestBody{
		Model: "gpt-4o",
		Messages: []OAIMessage{
			{Role: "system", Content: system},
			{Role: "user", Content: user},
		},
	}

	jsonData, err := json.Marshal(requestBody)

	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", OpenAIEndpoint, bytes.NewBuffer(jsonData))

	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("OPENAI")))

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Unmarshal the response body
	var responseBody OAIResponseBody
	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		return "", err
	}

	// Return the content of the first choice
	if len(responseBody.Choices) > 0 {
		return responseBody.Choices[0].Message.Content, nil
	}

	return "", fmt.Errorf("no response from OpenAI")

}

func CleanCode(input string) string {

	pattern := regexp.MustCompile("(?i)```[a-z]*\n|```")

	cleaned := pattern.ReplaceAllString(input, "")
	return cleaned
}
