package process

import (
	"fmt"
)

func Ingest(message string) (string, error) {
	mType, err := GetType(message)

	if err != nil {
		return "", err
	}

	id := ""

	switch mType {
	case "task":
		id, err = IngestTask(message)

		if err != nil {
			return "", err
		}
	case "note":
		id, err = IngestNote(message)

		if err != nil {
			return "", err
		}
	default:
		return "", fmt.Errorf("did not identify message type")
	}

	return id, nil
}

func Process(message string) (string, error) {
	id, err := Ingest(message)
	if err != nil && len(id) == 0 {
		// Retry once
		id, err = Ingest(message)
	}
	return id, err
}
