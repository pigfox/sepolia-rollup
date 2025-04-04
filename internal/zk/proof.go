package zk

import (
	"fmt"

	"github.com/consensys/gnark"
	"github.com/consensys/gnark-crypto/ecc/bn256"
)

func GenerateProof(inputs []int64) (*gnark.Proof, error) {
	// Create the zk-SNARK circuit
	circuit, err := SumCircuit()
	if err != nil {
		return nil, fmt.Errorf("failed to create circuit: %v", err)
	}

	// Set the inputs for the proof
	witness := circuit.NewWitness()

	// Map the inputs to the circuit variables
	witness.Set("a", inputs[0])
	witness.Set("b", inputs[1])

	// Generate the proof
	proof, err := gnark.GenerateProof(circuit, witness, bn256.New())
	if err != nil {
		return nil, fmt.Errorf("failed to generate proof: %v", err)
	}

	return proof, nil
}

func VerifyProof(proof *gnark.Proof) (bool, error) {
	// Create the zk-SNARK circuit
	circuit, err := SumCircuit()
	if err != nil {
		return false, fmt.Errorf("failed to create circuit: %v", err)
	}

	// Verify the proof
	isValid, err := gnark.VerifyProof(circuit, proof, bn256.New())
	if err != nil {
		return false, fmt.Errorf("proof verification failed: %v", err)
	}

	return isValid, nil
}
