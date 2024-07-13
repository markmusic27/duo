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

	//err := IngestTask(message)

	task := Task{
		Parent: ParentDatabase{
			Type:       "database_id",
			DatabaseID: "fa423041a28b4d7a9e8501372358177d",
		},
		Icon: Icon{
			Type:  "emoji",
			Emoji: "🌴",
		},
		Properties: TaskProperties{
			Name: NameWriteProp{
				Title: []TokenWrite{
					{
						Text: TextWrite{
							Content: "LETS GOOOO",
						},
					},
				},
			},
			Priority: SelectProp{
				Select: Select{
					Name: "P2",
				},
			},
			DueDate: DateProp{
				Date: Date{
					Start: "2024-07-13T22:47:00.000Z",
				},
			},
			Course: RelationProp{
				Pages: []Page{
					{
						ID: "4ec49e3e-f79b-4564-acba-9b4a5f596acf",
					},
				},
			},
			Project: RelationProp{
				Pages: []Page{
					{
						ID: "323db484-e9a3-419d-88e5-adf83b005b82",
					},
				},
			},
		},
	}

	id, err := CreateTask(task)

	if err != nil {
		log.Println("ERROR: " + err.Error())
	}

	log.Println(id)

	return nil
}
