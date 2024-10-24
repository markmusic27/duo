package process

import ( // hola - dr. izuz ♡ 07/10/24
	"fmt"
	"log"
)

func Ingest(message string, instructions ...string) (string, error) {
	var instruction string
	if len(instructions) > 0 {
		instruction = instructions[0]
	}

	mType, err := GetType(message)

	if err != nil {
		return "", err
	}

	id := ""

	switch mType {
	case "task":
		id, err = IngestTask(message, instruction)

		if err != nil {
			return "", err
		}
	case "note":
		id, err = IngestNote(message, instruction)

		if err != nil {
			return "", err
		}
	default:
		return "", fmt.Errorf("did not identify message type")
	}

	return id, nil
}

func Process(message string, instructions ...string) (string, error) {
	var instruction string
	if len(instructions) > 0 {
		instruction = instructions[0]
	}

	var id string
	var err error
	maxRetries := 4

	for attempt := 1; attempt <= maxRetries; attempt++ {
		id, err = Ingest(message, instruction)

		log.Println(id, err)
		if err == nil || len(id) > 0 {
			break
		}
		if attempt < maxRetries {
			Log("⚠️", fmt.Sprintf("Retry Attempt %d", attempt), "Retrying...")
		}
	}

	if err != nil && len(id) == 0 {
		err = fmt.Errorf("failed after %d attempts: %w", maxRetries, err)
	}

	return id, err
}
