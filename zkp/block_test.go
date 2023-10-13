// block_test.go
package zkp

import (
	"testing"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
)

func TestBlockValidation(t *testing.T) {
	// Compiling the circuit into R1CS
	var circuit CubicCircuit
	ccs, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &circuit)
	if err != nil {
		t.Fatalf("Failed to compile the circuit: %v", err)
	}

	// Groth16 zkSNARK setup
	pk, vk, err := groth16.Setup(ccs)
	if err != nil {
		t.Fatalf("Failed to set up groth16: %v", err)
	}

	// Defining the witness (assignment of specific values to circuit variables)
	assignment := CubicCircuit{X: 3, Y: 35}
	witness, err := frontend.NewWitness(&assignment, ecc.BN254.ScalarField())
	if err != nil {
		t.Fatalf("Failed to create a new witness: %v", err)
	}

	publicWitness, err := witness.Public()
	if err != nil {
		t.Fatalf("Failed to get public witness: %v", err)
	}

	// Generating a Groth16 proof (simulating block mining)
	proof, err := groth16.Prove(ccs, pk, witness)
	if err != nil {
		t.Fatalf("Failed to prove: %v", err)
	}

	// Creating a new block with the generated proof
	block := Block{
		Data:  "Block data",
		Proof: proof,
	}

	// Verifying the Groth16 proof (simulating block validation)
	err = groth16.Verify(block.Proof, vk, publicWitness) // Fixed this line
	if err != nil {
		t.Fatalf("Failed to verify the block: %v", err)
	}

	t.Log("Block validated successfully!")
}
