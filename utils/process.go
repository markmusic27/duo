package process

import (
	"fmt"
)

func Process(message string) error {
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
	default:
		return fmt.Errorf("did not identify message type")
	}

	return nil
}
