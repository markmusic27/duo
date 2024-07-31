package process

import (
	"regexp"
	"sort"
	"strings"
)

type Block struct {
	Object       string                  `json:"object"`
	Type         string                  `json:"type"`
	Paragraph    *map[string]interface{} `json:"paragraph,omitempty"`
	Code         *map[string]interface{} `json:"code,omitempty"`
	Heading1     *map[string]interface{} `json:"heading_1,omitempty"`
	Heading2     *map[string]interface{} `json:"heading_2,omitempty"`
	Heading3     *map[string]interface{} `json:"heading_3,omitempty"`
	Bullets      *map[string]interface{} `json:"bulleted_list_item,omitempty"`
	NumberedList *map[string]interface{} `json:"numbered_list_item,omitempty"`
	ToDo         *map[string]interface{} `json:"to_do,omitempty"`
	Quote        *map[string]interface{} `json:"quote,omitempty"`
}

type MarkdownRichText struct {
	Type       string                  `json:"type"`
	Text       *map[string]interface{} `json:"text,omitempty"`
	Equation   *map[string]interface{} `json:"equation,omitempty"`
	Annotation *map[string]interface{} `json:"annotations,omitempty"`
}

func ConvertMarkdownToNotion(markdown string) ([]Block, error) {
	lines := strings.Split(markdown, "\n")
	lines = RemoveEmptyLines(lines)

	for i, line := range lines {
		lines[i] = RemoveLeadingWhitespace(line)
	}

	var blocks []Block

	// Code block elemenets
	inCodeBlock := false
	var language strings.Builder
	var code strings.Builder

	for _, line := range lines {
		// Handle addition of code block
		if !inCodeBlock && strings.Contains(line, "```") {
			inCodeBlock = true
			language.WriteString(line[3:])
		} else if inCodeBlock && strings.Contains(line, "```") {
			// Appends the code block to the list
			blocks = append(blocks, createCodeBlock(language.String(), code.String()))

			inCodeBlock = false
			language.Reset()
			code.Reset()
		} else {
			indent := "\n"
			if len(code.String()) == 0 {
				indent = ""
			}
			code.WriteString(indent + line)
		}

		// Handle addition of single line blocks
		switch {
		case strings.HasPrefix(line, "# "):
			blocks = append(blocks, createHeading1Block(line[2:]))
		case strings.HasPrefix(line, "## "):
			blocks = append(blocks, createHeading2Block(line[3:]))
		case strings.HasPrefix(line, "### "):
			blocks = append(blocks, createHeading3Block(line[4:]))
		case strings.HasPrefix(line, "- ") && !strings.HasPrefix(line, "- [ ") && !strings.HasPrefix(line, "- [x"):
			blocks = append(blocks, createBulletList(line[2:]))
		case validateNumberedList(line):
			blocks = append(blocks, createNumberedList(line))
		case strings.HasPrefix(line, "- [ ]"):
			blocks = append(blocks, createTodo(line[5:], false))
		case strings.HasPrefix(line, "- [x]"):
			blocks = append(blocks, createTodo(line[5:], true))
		case strings.HasPrefix(line, "> "):
			blocks = append(blocks, createQuote(line[2:]))
		default:
			blocks = append(blocks, createParagraph(line))
		}
	}

	// Handles the case where a code block is not ended
	if inCodeBlock {
		blocks = append(blocks, createCodeBlock(language.String(), code.String()))
	}

	return blocks, nil
}

func createQuote(line string) Block {
	return Block{
		Object: "block",
		Type:   "quote",
		Quote: &map[string]interface{}{
			"rich_text": createRichText(line),
		},
	}
}

func createTodo(line string, complete bool) Block {
	return Block{
		Object: "block",
		Type:   "to_do",
		ToDo: &map[string]interface{}{
			"rich_text": createRichText(line),
			"checked":   complete,
		},
	}
}

func createNumberedList(line string) Block {
	itemPattern := regexp.MustCompile(`^\d+\. (.*)`)
	itemMatch := itemPattern.FindStringSubmatch(line)
	content := itemMatch[1]

	return Block{
		Object: "block",
		Type:   "numbered_list_item",
		NumberedList: &map[string]interface{}{
			"rich_text": createRichText(content),
		},
	}
}

func validateNumberedList(line string) bool {
	if len(line) < 3 {
		return false
	}

	numberPattern := regexp.MustCompile(`^(\d+)\.`)
	match := numberPattern.FindStringSubmatch(line)

	if !strings.Contains(line, ". ") {
		return false
	}

	return len(match) > 1
}

func createBulletList(item string) Block {
	return Block{
		Object: "block",
		Type:   "bulleted_list_item",
		Bullets: &map[string]interface{}{
			"rich_text": createRichText(item),
		},
	}
}

func createParagraph(paragraph string) Block {
	return Block{
		Object: "block",
		Type:   "paragraph",
		Paragraph: &map[string]interface{}{
			"rich_text": createRichText(paragraph),
		},
	}
}

func createRichText(paragraph string) []MarkdownRichText {
	components := FormatMarkdownParagraph(paragraph)
	text := []MarkdownRichText{}

	for _, component := range components {
		content := &map[string]interface{}{
			"content": component.Text,
		}

		if component.Link != nil {
			content = &map[string]interface{}{
				"content": component.Text,
				"link": &map[string]interface{}{
					"url": *component.Link,
				},
			}
		}
		isBold := false
		isStrikethrough := false
		isItalic := false
		isCode := false

		switch component.Format {
		case "B":
			isBold = true
		case "I":
			isItalic = true
		case "S":
			isStrikethrough = true
		case "C":
			isCode = true
		}

		text = append(text, MarkdownRichText{
			Type: "text",
			Text: content,
			Annotation: &map[string]interface{}{
				"bold":          isBold,
				"italic":        isItalic,
				"strikethrough": isStrikethrough,
				"underline":     false,
				"code":          isCode,
				"color":         "default",
			},
		})
	}

	return text
}

func createCodeBlock(language string, code string) Block {
	return Block{
		Object: "block",
		Type:   "code",
		Code: &map[string]interface{}{
			"language": language,
			"rich_text": []MarkdownRichText{
				{
					Type: "text",
					Text: &map[string]interface{}{
						"content": code,
					},
				},
			},
		},
	}
}

func createHeading1Block(line string) Block {
	return Block{
		Object: "block",
		Type:   "heading_1",
		Heading1: &map[string]interface{}{
			"rich_text":     createRichText(line),
			"is_toggleable": false,
			"color":         "default",
		},
	}
}

func createHeading2Block(line string) Block {
	return Block{
		Object: "block",
		Type:   "heading_2",
		Heading2: &map[string]interface{}{
			"rich_text":     createRichText(line),
			"is_toggleable": false,
			"color":         "default",
		},
	}
}

func createHeading3Block(line string) Block {
	return Block{
		Object: "block",
		Type:   "heading_3",
		Heading3: &map[string]interface{}{
			"rich_text":     createRichText(line),
			"is_toggleable": false,
			"color":         "default",
		},
	}
}

type FormattedText struct {
	Index  int
	End    int
	Text   string
	Format string
	Link   *string
}

type FormattedTextSlice []FormattedText

func (s FormattedTextSlice) Len() int {
	return len(s)
}

func (s FormattedTextSlice) Less(i, j int) bool {
	return s[i].Index < s[j].Index
}

func (s FormattedTextSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Function to sort and return a sorted slice
func SortFormattedTexts(texts []FormattedText) []FormattedText {
	// Create a copy of the slice to avoid modifying the original
	sortedTexts := make([]FormattedText, len(texts))
	copy(sortedTexts, texts)

	// Sort the copied slice
	sort.Sort(FormattedTextSlice(sortedTexts))

	return sortedTexts
}

func FormatMarkdownParagraph(input string) []FormattedText {
	// Indentify all formatable objects
	boldPattern := regexp.MustCompile(`\*\*(.*?)\*\*`)
	strikethroughPattern := regexp.MustCompile(`~~(.*?)~~`)
	linkPattern := regexp.MustCompile(`\[([^\n]+)\]\(([^\n]+)\)`)
	italicsPattern := regexp.MustCompile(`\*(.*?)\*`)
	codePattern := regexp.MustCompile("`(.*?)`")

	boldMatches := boldPattern.FindAllStringSubmatchIndex(input, -1)
	codeMatches := codePattern.FindAllStringSubmatchIndex(input, -1)
	strikethroughMatches := strikethroughPattern.FindAllStringSubmatchIndex(input, -1)
	linkMatches := linkPattern.FindAllStringSubmatchIndex(input, -1)
	rawItalicMatches := italicsPattern.FindAllStringSubmatchIndex(input, -1)
	italicMatches := [][]int{}

	cummulativeLength := len(boldMatches) + len(strikethroughMatches) + len(linkMatches) + len(italicMatches) + len(codeMatches)

	if cummulativeLength == 0 {
		return []FormattedText{
			{
				Index:  0,
				End:    len(input),
				Text:   input,
				Format: "N/A",
			},
		}
	}

	for _, irMatch := range rawItalicMatches {
		if (irMatch[3] - irMatch[2]) != 0 {
			italicMatches = append(italicMatches, irMatch)
		}
	}

	corpus := []FormattedText{}

	for _, bMatch := range boldMatches {
		corpus = append(corpus, FormattedText{
			Index:  bMatch[0],
			End:    bMatch[1],
			Text:   input[bMatch[2]:bMatch[3]],
			Format: "B",
		})
	}

	for _, sMatch := range strikethroughMatches {
		corpus = append(corpus, FormattedText{
			Index:  sMatch[0],
			End:    sMatch[1],
			Text:   input[sMatch[2]:sMatch[3]],
			Format: "S",
		})
	}

	for _, cMatch := range codeMatches {
		corpus = append(corpus, FormattedText{
			Index:  cMatch[0],
			End:    cMatch[1],
			Text:   input[cMatch[2]:cMatch[3]],
			Format: "C",
		})
	}

	for _, lMatch := range linkMatches {
		link := input[lMatch[4]:lMatch[5]]
		corpus = append(corpus, FormattedText{
			Index:  lMatch[0],
			End:    lMatch[1],
			Text:   input[lMatch[2]:lMatch[3]],
			Link:   &link,
			Format: "L",
		})
	}

	for _, iMatch := range italicMatches {
		corpus = append(corpus, FormattedText{
			Index:  iMatch[0],
			End:    iMatch[1],
			Text:   input[iMatch[2]:iMatch[3]],
			Format: "I",
		})
	}

	// Fill in the gaps
	corpus = SortFormattedTexts(corpus)
	unformattedCorpus := FormattedTextSlice{}

	if corpus[0].Index != 0 {
		unformattedCorpus = append(unformattedCorpus, FormattedText{
			Index:  0,
			End:    corpus[0].Index - 1,
			Text:   input[:(corpus[0].Index)],
			Format: "N/A",
		})
	}

	for i := 0; i < len(corpus); i++ {
		if !(len(corpus)-1 == i) {
			if (corpus[i+1].Index - corpus[i].End) != 0 {
				// Fill the gap
				unformattedCorpus = append(unformattedCorpus, FormattedText{
					Index:  corpus[i].End,
					End:    corpus[i+1].Index,
					Text:   input[corpus[i].End:corpus[i+1].Index],
					Format: "N/A",
				})
			}
		}
	}

	if corpus[len(corpus)-1].End != len(input) {
		unformattedCorpus = append(unformattedCorpus, FormattedText{
			Index:  corpus[len(corpus)-1].End,
			End:    len(input) - 1,
			Text:   input[corpus[len(corpus)-1].End:],
			Format: "N/A",
		})
	}

	corpus = append(corpus, unformattedCorpus...)
	corpus = SortFormattedTexts(corpus)

	return corpus
}
