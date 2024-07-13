package process

// ⬇️ Database Query Body
type DatabaseQueryBody struct {
	Filter *Filter `json:"filter,omitempty"`
}

type Filter struct {
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
	Title []Token `json:"rich_text"`
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
	Name           TextProp    `json:"Name"`
	Description    TextProp    `json:"Description"`
	Professor      TextProp    `json:"Professor"`
	ProfessorEmail EmailProp   `json:"Professor Email"`
	Location       TextProp    `json:"Location"`
	Website        WebsiteProp `json:"Website"`
}

func FetchCourses() []Course {

}

// ⬇️ Projects
