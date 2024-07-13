package process

import "log"

func Process(message string) error {
	// mType, err := GetType(message)

	// if err != nil {
	// 	return err
	// }

	// switch mType {
	// case "task":
	// 	// Add task code
	// default:
	// 	return fmt.Errorf("did not identify message type")
	// }

	err := IngestTask(message)

	if err != nil {
		log.Println("ERROR: " + err.Error())
	}

	return nil
}
