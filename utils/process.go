package process

import (
	"fmt"
)

const SuccessMessage = "Logged ✅"

func Process(message string, from string) error {
	mType, err := GetType(message)

	if err != nil {
		return err
	}

	switch mType {
	case "task":
		_, err := IngestTask(message)

		if err != nil {
			return err
		}

		Message(from, SuccessMessage)
	default:
		return fmt.Errorf("did not identify message type")
	}

	return nil
}
