package gnark2zkv

import (
	"encoding/json"
	"fmt"
	"os"

	// "github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/consensys/gnark/backend/witness"
)

// SaveToJSON saves a given data structure (v) or a public witness to a JSON file at the specified filePath.
// It intelligently handles publicWitness if the parameter is of type witness.Witness.
func SaveToJSON(filePath string, v interface{}) error {
	if witness, ok := v.(witness.Witness); ok {
		// rawVector, ok := witness.Vector().(fr.Vector)
		
		// if !ok {
		// 	return fmt.Errorf("failed to assert type of publicWitness.Vector() to fr.Vector")
		// }

		// witnessPublicStrings := make([]string, len(rawVector))
		// for i, val := range rawVector {
		// 	witnessPublicStrings[i] = val.String()
		// }

		return SaveToJSON(filePath, witness.Vector())
	}

	jsonData, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal data to JSON: %v", err)
	}

	err = os.WriteFile(filePath, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("failed to save JSON file: %v", err)
	}

	return nil
}
