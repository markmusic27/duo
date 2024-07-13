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
		// Add task code
	default:
		return fmt.Errorf("did not identify message type")
	}

	return nil
}
