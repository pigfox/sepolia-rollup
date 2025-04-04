package zk

import (
	"github.com/consensys/gnark"
	"github.com/consensys/gnark-crypto/ecc"
)

// Define the zk circuit structure
func SumCircuit() (*gnark.Circuit, error) {
	// Create a new circuit for the zk-SNARK
	circuit := gnark.NewCircuit()

	// Define the two inputs
	a := circuit.CreateVariable()
	b := circuit.CreateVariable()

	// Define the output sum
	sum := circuit.CreateVariable()

	// Add the constraint sum = a + b
	circuit.AddConstraint(a.Add(b).Equal(sum))

	return circuit, nil
}
