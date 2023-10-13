package zkp

import (
	"testing"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
)

// TestGroth16ProofVerification is a test function for Groth16 proof generation and verification
func TestGroth16ProofVerification(t *testing.T) {
	// compiles our circuit into a R1CS
	var circuit CubicCircuit
	ccs, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &circuit)
	if err != nil {
		t.Fatalf("Failed to compile the circuit: %v", err)
	}

	// groth16 zkSNARK: Setup
	pk, vk, err := groth16.Setup(ccs)
	if err != nil {
		t.Fatalf("Failed to set up groth16: %v", err)
	}

	// witness definition
	assignment := CubicCircuit{X: 3, Y: 35}
	witness, err := frontend.NewWitness(&assignment, ecc.BN254.ScalarField())
	if err != nil {
		t.Fatalf("Failed to create a new witness: %v", err)
	}
	publicWitness, err := witness.Public()
	if err != nil {
		t.Fatalf("Failed to get public witness: %v", err)
	}

	// groth16: Prove & Verify
	proof, err := groth16.Prove(ccs, pk, witness)
	if err != nil {
		t.Fatalf("Failed to prove: %v", err)
	}
	err = groth16.Verify(proof, vk, publicWitness)
	if err != nil {
		t.Fatalf("Failed to verify: %v", err)
	}

	t.Log("Proof verified successfully!")
}
