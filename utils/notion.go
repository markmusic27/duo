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

type PageText struct {
	Tokens []PageToken `json:"rich_text"`
}

type PageToken struct {
	Type string    `json:"type"`
	Text TextWrite `json:"text"`
}

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

type RichText struct {
	Components []RichTextComponent `json:"rich_text"`
}

type TextContent struct {
	Content string `json:"content"`
}

type RichTextComponent struct {
	Text TextContent `json:"text"`
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

// ⬇️ Areas / Interests
type AreaInterestResponse struct {
	AreasInterests []AreaInterest `json:"results"`
}

type AreaInterest struct {
	ID         string            `json:"id"`
	Properties AreaInterestProps `json:"properties"`
	URL        string            `json:"url"`
}

type AreaInterestProps struct {
	Name NameProp `json:"Name"`
}

type Note struct {
	Parent     ParentDatabase `json:"parent"`
	Icon       Icon           `json:"icon"`
	Properties NoteProperties `json:"properties"`
	Children   []Block        `json:"children"`
}

type NoteProperties struct {
	Name         NameWriteProp `json:"Name"`
	Type         SelectProp    `json:"Type"`
	AreaInterest RelationProp  `json:"Area / Interest"`
	Description  RichText      `json:"Description"`
	Course       RelationProp  `json:"Course"`
	Project      RelationProp  `json:"Project"`
}

func FetchAreasInterests(filter *Filter) ([]AreaInterest, error) {
	requestBody := DatabaseQueryBody{
		Filter: filter,
	}

	jsonData, err := json.Marshal(requestBody)

	if err != nil {
		return nil, err
	}

	endpoint := strings.ReplaceAll(NotionDatabaseQueryEndpoint, "*ID*", os.Getenv("AREAINTERESTID"))
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
	var response AreaInterestResponse
	err = json.Unmarshal(body, &response)

	if err != nil {
		return nil, err
	}

	areas := response.AreasInterests

	// Removes area of Capture from areas returned (prevents AI from adding them accidentally)
	for i, area := range response.AreasInterests {
		if area.Properties.Name.Title[0].Text == "Capture" {
			areas = append(areas[:i], areas[i+1:]...)
		}
	}

	return areas, nil
}

func CreateNote(note Note) (string, error) {
	jsonData, err := json.Marshal(note)

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

func FetchNoteTypes() ([]string, error) {
	endpoint := fmt.Sprintf("https://api.notion.com/v1/databases/%s", os.Getenv("NOTESID"))
	req, err := http.NewRequest("GET", endpoint, nil)

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
	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	return ExtractNoteTypes(data), nil
}

func ExtractNoteTypes(data map[string]interface{}) []string {
	var types []string

	properties, ok := data["properties"].(map[string]interface{})
	if !ok {
		return types
	}

	typeProperty, ok := properties["Type"].(map[string]interface{})
	if !ok {
		return types
	}

	selectProperty, ok := typeProperty["select"].(map[string]interface{})
	if !ok {
		return types
	}

	options, ok := selectProperty["options"].([]interface{})
	if !ok {
		return types
	}

	for _, option := range options {
		optionMap, ok := option.(map[string]interface{})
		if !ok {
			continue
		}
		if name, ok := optionMap["name"].(string); ok {
			types = append(types, name)
		}
	}

	return types
}

// ⬇️ Tasks

type Task struct {
	Parent     ParentDatabase `json:"parent"`
	Icon       Icon           `json:"icon"`
	Properties TaskProperties `json:"properties"`
	Children   []Block        `json:"children"`
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
