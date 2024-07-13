package process

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

	IngestTask(message)

	return nil
}
