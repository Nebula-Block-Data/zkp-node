# Blockchain Implementations in Go

This repository provides implementations and simulations of two critical aspects of blockchain technology: Ethereum mining using Proof of Work (PoW) and Zero-Knowledge Proofs using the `gnark` library. Each implementation is written in Go and serves as an educational resource to understand the underlying mechanisms of these technologies.

## Ethereum Mining Simulation in Go

This section simulates the Ethereum mining process using a simplified Proof of Work mechanism. It's an educational implementation, aiding in the understanding of how Ethereum mining, block validation, and consensus are achieved in a PoW blockchain.

### Key Features

- **Block Mining:** Simulates the process of mining a block by finding a nonce that satisfies specific conditions.
- **Multi-Threading:** Utilizes multiple threads (goroutines) to enhance the mining process.
- **Block Validation:** Includes a validation process to verify if the mined block meets the required conditions.

## Zero-Knowledge Proofs with gnark

This part of the repository demonstrates the implementation of zero-knowledge proofs using the `gnark` library. It’s an example of how a prover can prove the knowledge of secret information without revealing the information itself.

### Key Features

- **Cubic Circuit Definition:** A simple cubic equation circuit is defined to simulate the zero-knowledge proof.
- **Proof Generation:** Generates a zk-SNARK proof that the prover knows the values satisfying the circuit without revealing them.
- **Proof Verification:** Verifies the generated proof, validating that the prover knows the secret values.

## Getting Started

### Prerequisites

- Go programming language installed.
- Basic understanding of blockchain technology, Ethereum, and zero-knowledge proofs.

### Running the Simulations

1. **Clone the repository:**
   ```sh
   git clone https://github.com/Nebula-Block-Data/zkp-node.git
   cd zkp-node
   ```
2. Run the Ethereum Mining Simulation
3. Run the Zero-Knowledge Proof Simulation

## Acknowledgements
- The Ethereum community for extensive documentation and resources.
- The gnark library and ConsenSys for providing tools to work with zk-SNARKs.