package gnark2zkv

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"

	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
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
	err = ioutil.WriteFile(filePath, jsonData, 0644)
	if err != nil {
		return err
	}

	return nil
}

// ConvertPublicWitnessToJSON takes a public witness, extracts its data, and saves it as a JSON file.
// It returns the generated JSON file's data as a string array.
func ConvertPublicWitnessToJSON(publicWitness interface{}, outputFilePath string) ([]string, error) {
	// Convert public witness to binary
	bPublicWitness, err := publicWitness.(interface{ MarshalBinary() ([]byte, error) }).MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal public witness: %w", err)
	}

	// Skip the first 12 bytes (gnark specific header)
	bPublicWitness = bPublicWitness[12:]
	publicWitnessStr := hex.EncodeToString(bPublicWitness)

	// Decode hex to bytes
	inputBytes, err := hex.DecodeString(publicWitnessStr)
	if err != nil {
		return nil, fmt.Errorf("failed to decode public witness hex: %w", err)
	}

	// Validate input length
	if len(inputBytes)%fr.Bytes != 0 {
		return nil, fmt.Errorf("inputBytes mod fr.Bytes != 0")
	}

	// Convert bytes to big.Int slices
	nbInputs := len(inputBytes) / fr.Bytes
	input := make([]*big.Int, nbInputs)
	for i := 0; i < nbInputs; i++ {
		var e fr.Element
		e.SetBytes(inputBytes[fr.Bytes*i : fr.Bytes*(i+1)])
		input[i] = new(big.Int)
		e.BigInt(input[i])
	}

	// Convert big.Int to string array
	var witnessPublic []string
	for _, bi := range input {
		witnessPublic = append(witnessPublic, bi.String())
	}

	// Save to JSON
	err = SaveToJSON(outputFilePath, witnessPublic)
	if err != nil {
		return nil, fmt.Errorf("failed to save witness to JSON: %w", err)
	}

	return witnessPublic, nil
}
