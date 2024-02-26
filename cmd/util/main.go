package util

import (
	"encoding/json"
	"os"
)

func SaveStructToJSON(data interface{}, path string) error {
	jsonData, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return err
	}

	err = os.WriteFile(path, jsonData, 0644)
	if err != nil {
		return err
	}

	return nil
}
