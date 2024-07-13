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
const NotionPageCreationEndpoint = "https://api.notion.com/v1/pages"

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

type SelectProp struct {
	Select Select `json:"select"`
}

type Select struct {
	Name string `json:"name"`
}

type DateProp struct {
	Date Date `json:"date"`
}

type Date struct {
	Start string `json:"start"`
}

type RelationProp struct {
	Pages []Page `json:"relation"`
}

type Page struct {
	ID string `json:"id"`
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

type Icon struct {
	Type  string `json:"type"`
	Emoji string `json:"emoji"`
}

type NameWriteProp struct {
	Title []TokenWrite `json:"title"`
}

type TokenWrite struct {
	Text TextWrite `json:"text"`
}

type TextWrite struct {
	Content string `json:"content"`
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

// ⬇️ Tasks

type Task struct {
	Parent     ParentDatabase `json:"parent"`
	Icon       Icon           `json:"icon"`
	Properties TaskProperties `json:"properties"`
}

type ParentDatabase struct {
	Type       string `json:"type"`
	DatabaseID string `json:"database_id"`
}

type TaskProperties struct {
	Name     NameWriteProp `json:"Name"`
	Priority SelectProp    `json:"Priority"`
	DueDate  DateProp      `json:"Due Date"`
	Course   RelationProp  `json:"Course"`
	Project  RelationProp  `json:"Project"`
}

func CreateTask(task Task) (string, error) {
	jsonData, err := json.Marshal(task)

	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", NotionPageCreationEndpoint, bytes.NewBuffer(jsonData))

	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Notion-Version", "2022-06-28")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("NOTION")))

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

	// Unmarshall body
	var page Page
	err = json.Unmarshal(body, &page)

	if err != nil {
		return "", fmt.Errorf("Failed to create page: " + err.Error() + "\n\nNotion Response:\n" + string(body))
	}

	return page.ID, nil
}
