package process

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type RequestBody struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type ResponseChoice struct {
	Message Message `json:"message"`
}

type ResponseBody struct {
	Choices []ResponseChoice `json:"choices"`
}

const OpenAIEndpoint = "https://api.openai.com/v1/chat/completions"

func Prompt(user string, system string) (string, error) {
	requestBody := RequestBody{
		Model: "gpt-4o",
		Messages: []Message{
			{Role: "system", Content: system},
			{Role: "user", Content: user},
		},
	}

	jsonData, err := json.Marshal(requestBody)

	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", OpenAIEndpoint, bytes.NewBuffer(jsonData))
}
