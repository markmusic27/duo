package process

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

// ⬇️ Notion Endpoints

const NotionDatabaseQueryEndpoint = "https://api.notion.com/v1/databases/*ID*/query"

// ⬇️ Database Query Body
type DatabaseQueryBody struct {
	Filter *Filter `json:"filter,omitempty"`
}

type Filter struct {
	Property string          `json:"property"`
	Checkbox *CheckboxFilter `json:"checkbox,omitempty"`
	Status   *StatusFilter   `json:"status"`
}

type CheckboxFilter struct {
	Equals bool `json:"equals"`
}

type StatusFilter struct {
	Equals string `json:"equals"`
}

// ⬇️ General Props
type EmailProp struct {
	Email string `json:"email"`
}

type WebsiteProp struct {
	URL string `json:"url"`
}

type TextProp struct {
	Tokens []Token `json:"rich_text"`
}

type Token struct {
	Text string `json:"plain_text"`
}

type NameProp struct {
	Title []Token `json:"title"`
}

// ⬇️ Course

type CourseResponse struct {
	Courses []Course `json:"results"`
}

type Course struct {
	ID         string           `json:"id"`
	Properties CourseProperties `json:"properties"`
	URL        string           `json:"url"`
}

type CourseProperties struct {
	Name           NameProp    `json:"Name"`
	Description    TextProp    `json:"Description"`
	Professor      TextProp    `json:"Professor"`
	ProfessorEmail EmailProp   `json:"Professor Email"`
	Location       TextProp    `json:"Location"`
	Website        WebsiteProp `json:"Website"`
}

func FetchCourses(filter *Filter) ([]Course, error) {
	requestBody := DatabaseQueryBody{
		Filter: filter,
	}

	jsonData, err := json.Marshal(requestBody)

	if err != nil {
		return nil, err
	}

	endpoint := strings.ReplaceAll(NotionDatabaseQueryEndpoint, "*ID*", os.Getenv("COURSEID"))
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonData))

	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Notion-Version", "2022-06-28")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("NOTION")))

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Unmarshal the response body
	var responseBody CourseResponse
	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		return nil, err
	}

	return responseBody.Courses, nil
}

// ⬇️ Projects

type Project struct {
	ID         string            `json:"id"`
	Properties ProjectProperties `json:"properties"`
	URL        string            `json:"url"`
}

type ProjectProperties struct {
	Name NameProp `json:"name"`
}

type ProjectResponse struct {
	Projects []Project `json:"results"`
}

func FetchProjects(filter *Filter) ([]Project, error) {
	requestBody := DatabaseQueryBody{
		Filter: filter,
	}

	jsonData, err := json.Marshal(requestBody)

	if err != nil {
		return nil, err
	}

	endpoint := strings.ReplaceAll(NotionDatabaseQueryEndpoint, "*ID*", os.Getenv("PROJECTID"))
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonData))

	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Notion-Version", "2022-06-28")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("NOTION")))

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Unmarshall body
	var response ProjectResponse
	err = json.Unmarshal(body, &response)

	if err != nil {
		return nil, err
	}

	return response.Projects, nil
}
