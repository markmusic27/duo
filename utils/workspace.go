package process

import (
	"encoding/json"
	"fmt"
	"os"
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

type GeneratedTask struct {
	Emoji    string   `json:"emoji"`
	Task     string   `json:"task"`
	Deadline string   `json:"deadline"`
	Priority int64    `json:"priority"`
	Body     string   `json:"body"`
	Course   []string `json:"course"`
	Project  []string `json:"project"`
}

func IngestTask(task string) (string, error) {
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
		return "", err
	}

	courseContext := ""
	for _, course := range courses {
		courseContext = courseContext + fmt.Sprintf("\n\t - %s: %s", course.Properties.Name.Title[0].Text, course.ID)
	}

	template = strings.ReplaceAll(template, "*COURSES*", courseContext)

	// Add Project context
	projects, err := FetchProjects(filter)

	if err != nil {
		return "", err
	}

	projectContext := ""
	for _, project := range projects {
		projectContext = projectContext + fmt.Sprintf("\n\t - %s: %s", project.Properties.Name.Title[0].Text, project.ID)
	}

	template = strings.ReplaceAll(template, "*PROJECTS*", projectContext)

	// Send request to OpenAI
	res, err := Prompt(task, template)

	if err != nil {
		return "", err
	}

	var generated GeneratedTask
	err = json.Unmarshal([]byte(CleanCode(res)), &generated)

	if err != nil {
		return "", fmt.Errorf("Failed to unmarshall: " + err.Error() + "\n\n Generated Task:\n" + res)
	}

	// Creates new task object based on information

	newTask := Task{
		Parent: ParentDatabase{
			Type:       "database_id",
			DatabaseID: os.Getenv("TASKID"),
		},
		Icon: Icon{
			Type:  "emoji",
			Emoji: generated.Emoji,
		},
		Properties: TaskProperties{
			Name: NameWriteProp{
				Title: []TokenWrite{
					{
						Text: TextWrite{
							Content: generated.Task,
						},
					},
				},
			},
			Priority: SelectProp{
				Select: Select{
					Name: ConvertNumToPriority(generated.Priority),
				},
			},
			DueDate: DateProp{
				Date: Date{
					Start: generated.Deadline,
				},
			},
			Course: RelationProp{
				Pages: GeneratePageFromStrings(generated.Course),
			},
			Project: RelationProp{
				Pages: GeneratePageFromStrings(generated.Project),
			},
		},
	}

	id, err := CreateTask(newTask)

	if err != nil {
		return "", err
	}

	return id, nil
}
