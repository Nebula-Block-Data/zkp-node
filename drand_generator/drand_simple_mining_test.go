package drand_generator

import (
	"encoding/hex"
	"math/big"
	"sync"
	"testing"
	"time"
)

func validateBlock(block Block, target *big.Int) bool {
	var hashInt big.Int
	hashBytes, _ := hex.DecodeString(block.Hash)
	hashInt.SetBytes(hashBytes)

	// Check if the hash is less than the target
	return hashInt.Cmp(target) == -1
}

func TestEthBlockMiningWithDrandAndMultithreading(t *testing.T) {
	randomness := getDrandRandomness()
	randomInt := new(big.Int).SetBytes(randomness[:16])                // Using first 16 bytes to get a big int
	maxTarget := new(big.Int).Exp(big.NewInt(2), big.NewInt(355), nil) // Maximum possible target
	target := new(big.Int).Div(maxTarget, randomInt)                   // Adjust the target based on the randomness

	blockData := "This is some block data"

	resultChan := make(chan Block, 1) // Buffer to hold the mined block

	var wg sync.WaitGroup

	numGoroutines := 16
	startTime := time.Now() // Record the start time

	for i := uint64(0); i < uint64(numGoroutines); i++ {
		wg.Add(1)
		go mine(blockData, target, i, uint64(numGoroutines), resultChan, &wg)
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	block := <-resultChan

	endTime := time.Now() // Record the end time
	t.Logf("Time taken to mine: %s", endTime.Sub(startTime))
	t.Logf("Target: %d", target)
	t.Logf("Randomness: %d", randomness)
	t.Logf("Nonce: %d", block.Nonce)
	t.Logf("Hash: %s", block.Hash)

	// Validate the block
	isValid := validateBlock(block, target)
	if isValid {
		t.Log("Block is valid.")
	} else {
		t.Error("Block is invalid.")
	}
}
