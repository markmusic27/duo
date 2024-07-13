package process

import (
	"encoding/json"
	"strings"
)

var Databases = []string{"task", "note"}

// ⬇️ Get Type

type GetTypeResponseBody struct {
	Type string `json:"type"`
}

func GetType(message string) (string, error) {
	db := ""

	for i := 0; i < len(Databases); i++ {
		var filler string
		if i == 0 {
			filler = ""
		} else {
			filler = "\n"
		}

		db = db + filler + "- " + Databases[i]
	}

	template := strings.ReplaceAll(TypeTemplate, "*DB*", db)

	raw, err := Prompt(message, template)

	if err != nil {
		return "", err
	}

	var res GetTypeResponseBody

	err = json.Unmarshal([]byte(CleanCode(raw)), &res)

	if err != nil {
		return "", err
	}

	return res.Type, nil
}

// ⬇️ Tasks
