package process

import ( // hola - dr. izuz ♡ 07/10/24
	"fmt"
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

	id, err := Ingest(message, instruction)
	if err != nil && len(id) == 0 {
		// Retry once
		id, err = Ingest(message, instruction)
	}
	return id, err
}
