package main

import (
	"fmt"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"

	"github.com/zkVerify/gnark2zkv"
)

type Circuit struct {
	A frontend.Variable
	B frontend.Variable `gnark:",public"`
	C frontend.Variable `gnark:",public"`
}

func (circuit *Circuit) Define(api frontend.API) error {
	A2 := api.Mul(circuit.A, circuit.A)
	C2 := api.Mul(circuit.C, circuit.C)
	api.AssertIsEqual(circuit.B, api.Add(A2, C2))
	return nil
}

func main() {
	var circuit Circuit

	// 1. Circuit definition
	_r1cs, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &circuit)
	if err != nil {
		panic(err)
	}

	// 2. Setup
	pk, vk, err := groth16.Setup(_r1cs)
	if err != nil {
		panic(err)
	}

	// 3. Witness definition
	assignment := Circuit{
		A: "3",
		B: "10",
		C: "1",
	}
	witness, _ := frontend.NewWitness(&assignment, ecc.BN254.ScalarField())
	publicWitness, _ := witness.Public()

	// 4. Proof creation
	// groth16: Prove & Verify
	proof, err := groth16.Prove(_r1cs, pk, witness)
	if err != nil {
		panic(err)
	}

	// 5. Proof verification
	err = groth16.Verify(proof, vk, publicWitness)
	if err != nil {
		panic(err)
	}

	// ******************* Save publicWitness as JSON file **********************
	rawVector, ok := publicWitness.Vector().(fr.Vector)
	if !ok {
		panic("Failed to assert type of publicWitness.Vector() to fr.Vector")
	}
	witnessPublicStrings := make([]string, len(rawVector))
	for i, val := range rawVector {
		witnessPublicStrings[i] = val.String()
	}
	err = gnark2zkv.SaveToJSON("WitnessPublic.json", witnessPublicStrings)
	if err != nil {
		fmt.Printf("Error converting public witness to JSON: %v\n", err)
		return
	}

	fmt.Println("Witness JSON saved to WitnessPublic.json")
	// ******************* Save Proof as JSON file ******************************
	err = gnark2zkv.SaveToJSON("proof.json", proof)
	if err != nil {
		fmt.Println("Error saving proof to JSON:", err)
		return
	}

	fmt.Println("Proof saved to proof.json")
	// ******************* Save VK as JSON file *********************************
	err = gnark2zkv.SaveToJSON("vk.json", vk)
	if err != nil {
		fmt.Println("Error saving VK to JSON:", err)
		return
	}

	fmt.Println("VK saved to vk.json")
	// **************************************************************************

}
