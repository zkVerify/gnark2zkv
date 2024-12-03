package gnark2zkv

import (
	"encoding/json"
	"os"
)

// SaveToJSON saves a given data structure (v) to a JSON file at the specified filePath.
// It returns an error if the operation fails.
func SaveToJSON(filePath string, v interface{}) error {
	// MarshalIndent generates a JSON string with indentation for better readability.
	jsonData, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}

	// Write the JSON data to the specified file.
	err = os.WriteFile(filePath, jsonData, 0644)
	if err != nil {
		return err
	}

	return nil
}
