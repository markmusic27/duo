package process

import (
	"encoding/json"
	"fmt"
	"net/url"
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

	raw, err := Prompt(message, template, "gpt-3.5-turbo")

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
		Children: []Child{
			{
				Object: "block",
				Type:   "paragraph",
				Paragraph: &PageText{
					Tokens: []PageToken{
						{
							Type: "text",
							Text: TextWrite{
								Content: generated.Body,
							},
						},
					},
				},
			},
		},
	}

	id, err := CreateTask(newTask)

	if err != nil {
		return "", err
	}

	return id, nil
}

// ⬇️ Notes

type GeneratedNote struct {
	Emoji       string   `json:"emoji"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Type        string   `json:"type"`
	Area        []string `json:"area"`
	Project     []string `json:"project"`
}

func IngestNote(note string) (string, error) {
	urls, message := ExtractLinksAndReplaceDomains(note)

	linkContext := ""
	if len(urls) != 0 {
		for i, urlI := range urls {
			linkData, err := ChannelLink(urlI)

			if err == nil {
				var newLine string

				if i == 0 {
					newLine = ""
				} else {
					newLine = "\n"
				}

				linkContext = linkContext + newLine + linkData
			}
		}
	}

	template := NoteTemplate

	// Add types context
	types, err := FetchNoteTypes()
	if err != nil {
		return "", err
	}

	typesContext := ""
	for _, t := range types {
		typesContext = typesContext + "\n\t- " + fmt.Sprintf(`"%s"`, t)
	}

	template = strings.ReplaceAll(template, "*TYPES*", typesContext)

	// Add projects context
	filter := &Filter{
		Property: "Status",
		Status: &StatusFilter{
			Equals: "In progress",
		},
	}
	projects, err := FetchProjects(filter)

	if err != nil {
		return "", err
	}

	projectContext := ""
	for _, project := range projects {
		projectContext = projectContext + fmt.Sprintf("\n\t - %s: %s", project.Properties.Name.Title[0].Text, project.ID)
	}

	template = strings.ReplaceAll(template, "*PROJECTS*", projectContext)

	// Add area/interest context
	areasinterests, err := FetchAreasInterests(nil)
	if err != nil {
		return "", err
	}

	areainterestContext := ""
	for _, areasinterest := range areasinterests {
		areainterestContext = areainterestContext + fmt.Sprintf("\n\t - %s: %s", areasinterest.Properties.Name.Title[0].Text, areasinterest.ID)
	}

	template = strings.ReplaceAll(template, "*AREAS*", areainterestContext)

	// Make request to OpenAI servers
	userMessage := fmt.Sprintf("Original note: %s\n\n%s", message, linkContext)

	gen, err := Prompt(userMessage, template)
	if err != nil {
		return "", err
	}

	var generated GeneratedNote
	err = json.Unmarshal([]byte(CleanCode(gen)), &generated)
	if err != nil {
		return "", err
	}

	newNote := Note{
		Parent: ParentDatabase{
			Type:       "database_id",
			DatabaseID: os.Getenv("NOTESID"),
		},
		Icon: Icon{
			Type:  "emoji",
			Emoji: generated.Emoji,
		},
		Properties: NoteProperties{
			Name: NameWriteProp{
				Title: []TokenWrite{
					{
						Text: TextWrite{
							Content: generated.Title,
						},
					},
				},
			},
			Type: SelectProp{
				Select: Select{
					Name: generated.Type,
				},
			},
			Description: RichText{
				Components: []RichTextComponent{
					{
						Text: TextContent{
							Content: generated.Description,
						},
					},
				},
			},
			AreaInterest: RelationProp{
				Pages: GeneratePageFromStrings(generated.Area),
			},
			Project: RelationProp{
				Pages: GeneratePageFromStrings(generated.Project),
			},
		},
		Children: CreateBookmarksFromURLs(urls),
	}

	id, err := CreateNote(newNote)
	if err != nil {
		return "", err
	}

	return id, nil
}

func CreateBookmarksFromURLs(urls []string) []Child {
	children := []Child{}

	for _, url := range urls {
		child := Child{
			Object: "block",
			Type:   "bookmark",
			Bookmark: &Bookmark{
				Caption: []string{},
				URL:     url,
			},
		}
		children = append(children, child)
	}

	return children
}

func ChannelLink(link string) (string, error) {
	u, err := url.Parse(link)
	if err != nil {
		return "", fmt.Errorf("URL is not a link")
	}

	switch u.Host {
	case "www.youtube.com", "youtube.com", "youtu.be":
		data, err := FetchYoutubeData(link)
		if err != nil {
			return "", err
		}

		return data, nil
	case "www.instagram.com", "instagram.com":
		// TODO: Add Instagram API integration

		return "", nil
	default:
		data, err := ContextualizeLink(link)
		if err != nil {
			return "", err
		}

		return fmt.Sprintf("Website Context:\n\t- URL: %s\n\t- Content Summary: %s", link, data), nil
	}

}
