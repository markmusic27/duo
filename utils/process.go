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

	id, err := IngestTask(message)

	if err != nil {
		log.Println(err)
	}

	log.Println(id)

	return nil
}
