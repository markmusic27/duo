package process

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"
)

var Databases = []string{"task", "note"}

// ⬇️ Get Type

type GetTypeResponseBody struct {
	Type string `json:"type"`
}

func GetType(message string) (string, error) {
	db := ""

	for i := 0; i < len(Databases); i++ {
		var filler string
		if i == 0 {
			filler = ""
		} else {
			filler = "\n"
		}

		db = db + filler + "- " + Databases[i]
	}

	template := strings.ReplaceAll(TypeTemplate, "*DB*", db)

	raw, err := Prompt(message, template)

	if err != nil {
		return "", err
	}

	var res GetTypeResponseBody

	err = json.Unmarshal([]byte(CleanCode(raw)), &res)

	if err != nil {
		return "", err
	}

	return res.Type, nil
}

// ⬇️ Tasks

func IngestTask(task string) error {
	template := TaskTemplate

	// Add date information
	currentTime := time.Now()
	template = strings.ReplaceAll(template, "*DATE*", currentTime.Format(time.RFC3339))
	template = strings.ReplaceAll(template, "*WEEKDAY*", currentTime.Weekday().String())

	// Add Course context
	filter := &Filter{
		Property: "Status",
		Status: &StatusFilter{
			Equals: "In progress",
		},
	}
	courses, err := FetchCourses(filter)

	if err != nil {
		return err
	}

	courseContext := ""
	for _, course := range courses {
		courseContext = courseContext + fmt.Sprintf("\n\t - %s: %s", course.Properties.Name.Title[0].Text, course.ID)
	}

	template = strings.ReplaceAll(template, "*COURSES*", courseContext)

	// Add Project context
	projects, err := FetchProjects(filter)

	if err != nil {
		return err
	}

	projectContext := ""
	for _, project := range projects {
		projectContext = projectContext + fmt.Sprintf("\n\t - %s: %s", project.Properties.Name.Title[0].Text, project.ID)
	}

	template = strings.ReplaceAll(template, "*PROJECTS*", projectContext)

	log.Println(template)

	return nil
}
