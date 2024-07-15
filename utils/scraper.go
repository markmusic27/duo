package process

import (
	"fmt"
	"io"
	"net/http"

	md "github.com/JohannesKaufmann/html-to-markdown"
)

func FetchPageContent(url string) (string, error) {
	// Make the GET request
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to make GET request: %v", err)
	}
	defer resp.Body.Close()

	// Check for non-200 status codes
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("non-200 HTTP status: %d", resp.StatusCode)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	// Convert HTML to Markdown
	converter := md.NewConverter("", true, nil)
	markdown, err := converter.ConvertString(string(body))
	if err != nil {
		return "", fmt.Errorf("failed to convert HTML to Markdown: %v", err)
	}

	return markdown, nil
}
