package main

import (
	"fmt"
	"log"
	"os"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
)

// CubicCircuit defines a simple circuit
// x**3 + x + 5 == y
type CubicCircuit struct {
	X frontend.Variable `gnark:"x"`
	Y frontend.Variable `gnark:",public"`
}

// Define declares the circuit constraints
// x**3 + x + 5 == y
func (circuit *CubicCircuit) Define(api frontend.API) error {
	x3 := api.Mul(circuit.X, circuit.X, circuit.X)
	api.AssertIsEqual(circuit.Y, api.Add(x3, circuit.X, 5))
	return nil
}

func main() {
	// compiles our circuit into a R1CS
	var circuit CubicCircuit
	ccs, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &circuit)
	if err != nil {
		log.Printf("Failed to compile the circuit: %v\n", err)
		os.Exit(1)
	}

	// groth16 zkSNARK: Setup
	pk, vk, err := groth16.Setup(ccs)
	if err != nil {
		log.Printf("Failed to set up groth16: %v\n", err)
		os.Exit(1)
	}

	// witness definition
	assignment := CubicCircuit{X: 3, Y: 35}
	witness, err := frontend.NewWitness(&assignment, ecc.BN254.ScalarField())
	if err != nil {
		log.Printf("Failed to create a new witness: %v\n", err)
		os.Exit(1)
	}
	publicWitness, err := witness.Public()
	if err != nil {
		log.Printf("Failed to get public witness: %v\n", err)
		os.Exit(1)
	}

	// groth16: Prove & Verify
	proof, err := groth16.Prove(ccs, pk, witness)
	if err != nil {
		log.Printf("Failed to prove: %v\n", err)
		os.Exit(1)
	}
	err = groth16.Verify(proof, vk, publicWitness)
	if err != nil {
		log.Printf("Failed to verify: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Proof verified successfully!")
}
