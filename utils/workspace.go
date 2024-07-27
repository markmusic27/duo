package process

import (
	"encoding/json"
	"fmt"
	"log"
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

	raw, err := Prompt(message, template, SecondaryModel)

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

type GeneratedDeadline struct {
	Deadline string `json:"deadline"`
}

type GeneratedTask struct {
	Emoji    string   `json:"emoji"`
	Task     string   `json:"task"`
	Priority int64    `json:"priority"`
	Body     string   `json:"body"`
	Course   []string `json:"course"`
	Project  []string `json:"project"`
}

func IngestTask(task string, instructions ...string) (string, error) {
	var instruction string
	if len(instructions) > 0 {
		instruction = instructions[0]
	}

	template := TaskTemplate

	// Loads the adequate
	location, err := time.LoadLocation(os.Getenv("LOCATION"))
	if err != nil {
		return "", err
	}

	// Add date information
	t := time.Now()
	currentTime := t.In(location)
	template = strings.ReplaceAll(template, "*DATE*", currentTime.Format(time.RFC3339))
	template = strings.ReplaceAll(template, "*WEEKDAY*", currentTime.Weekday().String())

	// Add Course context
	courseFilter := &Filter{
		Property: "Status",
		Status: &StatusFilter{
			Equals: "In progress",
		},
	}
	courses, err := FetchCourses(courseFilter)

	if err != nil {
		return "", err
	}

	courseContext := ""
	for _, course := range courses {
		courseContext = courseContext + fmt.Sprintf("\n\t - %s: %s", course.Properties.Name.Title[0].Text, course.ID)
	}

	template = strings.ReplaceAll(template, "*COURSES*", courseContext)

	// Add Project context
	projectFilter := &Filter{
		Property: "Status",
		Status: &StatusFilter{
			Equals: "Active",
		},
	}
	projects, err := FetchProjects(projectFilter)

	if err != nil {
		return "", err
	}

	projectContext := ""
	for _, project := range projects {
		projectContext = projectContext + fmt.Sprintf("\n\t - %s: %s", project.Properties.Name.Title[0].Text, project.ID)
	}

	// Obtain deadline and priority
	deadlineTemplate := strings.ReplaceAll(DeadlineTemplate, "*DATE*", currentTime.Format(time.RFC3339))
	deadlineTemplate = strings.ReplaceAll(deadlineTemplate, "*WEEKDAY*", currentTime.Weekday().String())

	timezone, err := ExtractTimezone(currentTime.Format(time.RFC3339))
	if err != nil {
		return "", err
	}

	deadlineTemplate = strings.ReplaceAll(deadlineTemplate, "*TIMEZONE*", timezone)

	deadlineRes, err := Prompt(task, deadlineTemplate)
	if err != nil {
		return "", err
	}

	var generatedDeadline GeneratedDeadline
	err = json.Unmarshal([]byte(CleanCode(deadlineRes)), &generatedDeadline)
	if err != nil {
		return "", err
	}

	template = strings.ReplaceAll(template, "*PROJECTS*", projectContext)

	// Add instructions if necessary
	if len(instruction) > 0 {
		template = template + InstructionPreamble + instruction
	}

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

	emoji := string([]rune(generated.Emoji)[0])
	body, err := ConvertMarkdownToNotion(generated.Body)
	if err != nil {
		return "", nil
	}

	newTask := Task{
		Parent: ParentDatabase{
			Type:       "database_id",
			DatabaseID: os.Getenv("TASKID"),
		},
		Icon: Icon{
			Type:  "emoji",
			Emoji: emoji,
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
					Start: generatedDeadline.Deadline,
				},
			},
			Course: RelationProp{
				Pages: GeneratePageFromStrings(generated.Course),
			},
			Project: RelationProp{
				Pages: GeneratePageFromStrings(generated.Project),
			},
		},
		Children: body,
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

func IngestNote(note string, instructions ...string) (string, error) {
	var instruction string
	if len(instructions) > 0 {
		instruction = instructions[0]
	}

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
			Equals: "Active",
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

	// Add instructions if necessary
	if len(instruction) > 0 {
		template = template + InstructionPreamble + instruction
	}

	gen, err := Prompt(userMessage, template)
	if err != nil {
		return "", err
	}

	var generated GeneratedNote
	err = json.Unmarshal([]byte(CleanCode(gen)), &generated)
	if err != nil {
		return "", err
	}

	if len(generated.Area) == 0 && len(generated.Project) == 0 {
		generated.Area = append(generated.Area, os.Getenv("COURSE_PAGE_ID"))
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

func CreateBookmarksFromURLs(urls []string) []Block {
	var markdown strings.Builder

	markdown.WriteString("### Resources")

	for _, url := range urls {
		markdown.WriteString(fmt.Sprintf("\n- [%s](%s)", url, url))
	}

	log.Println(markdown.String())

	b, err := ConvertMarkdownToNotion(markdown.String())
	if err != nil {
		b = []Block{}
	}

	return b
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
