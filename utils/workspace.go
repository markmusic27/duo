package process

import (
	"encoding/json"
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

func IngestTask(task string) (string, error) {
	template := TaskTemplate

	// Add date information
	currentTime := time.Now()
	template = strings.ReplaceAll(template, "*DATE*", currentTime.Format(time.RFC3339))
	template = strings.ReplaceAll(template, "*WEEKDAY*", currentTime.Weekday().String())

	// Add Notion context
	FetchCourses()

	return template, nil
}
