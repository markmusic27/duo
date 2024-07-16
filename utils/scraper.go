package process

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
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

func ContextualizeLink(url string) (string, error) {
	body, err := FetchPageContent(url)

	if err != nil {
		return "", err
	}

	const limit = 40000
	var capped string

	if len(body) > limit {
		capped = body[:limit]
	} else {
		capped = body
	}

	summary, err := Prompt(capped, SummarizationTemplate, "gpt-3.5-turbo")

	if err != nil {
		return "", err
	}

	return summary, nil
}

// ⬇️ Youtube
func ExtractYoutubeID(video string) (string, error) {
	u, err := url.Parse(video)
	if err != nil {
		return "", err
	}

	// Check if it's a valid YouTube URL
	if u.Host != "www.youtube.com" && u.Host != "youtube.com" && u.Host != "youtu.be" {
		return "", fmt.Errorf("invalid YouTube URL")
	}

	var videoID string
	if u.Host == "youtu.be" {
		// Handle short URL format (e.g., https://youtu.be/dQw4w9WgXcQ)
		videoID = strings.TrimPrefix(u.Path, "/")
	} else {
		// Handle long URL format (e.g., https://www.youtube.com/watch?v=dQw4w9WgXcQ)
		queryParams := u.Query()
		videoID = queryParams.Get("v")
	}

	if videoID == "" {
		return "", fmt.Errorf("no video ID found in URL")
	}

	return videoID, nil
}

func FetchYoutubeData(video string) (string, error) {
	ctx := context.Background()
	service, err := youtube.NewService(ctx, option.WithAPIKey(os.Getenv("GCP_KEY")))
	if err != nil {
		return "", err
	}

	id, err := ExtractYoutubeID(video)
	if err != nil {
		return "", err
	}

	call := service.Videos.List([]string{"snippet"}).Id(id)
	res, err := call.Do()
	if err != nil {
		return "", err
	}

	if len(res.Items) == 0 {
		return "", fmt.Errorf("no video found with the given ID")
	}

	data := fmt.Sprintf(`
YouTube Video Context:
	- Title: %s
	- Channel: %s
	- Description: %s
	`, res.Items[0].Snippet.Title, res.Items[0].Snippet.ChannelTitle, RemoveNewline(ClampString(res.Items[0].Snippet.Description, 250)))

	return data, nil
}
