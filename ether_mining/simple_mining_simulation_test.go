package ether_mining

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
	"sync"
	"testing"
	"time"
)

func mine(blockData string, target *big.Int, startNonce uint64, nonceStep uint64, resultChan chan<- Block, wg *sync.WaitGroup) {
	defer wg.Done()

	var hashInt big.Int
	nonce := startNonce

	for {
		data := fmt.Sprintf("%s%d", blockData, nonce)
		hashBytes := sha256.Sum256([]byte(data))
		hashInt.SetBytes(hashBytes[:])

		if hashInt.Cmp(target) == -1 {
			resultChan <- Block{
				Data:      blockData,
				Nonce:     nonce,
				Hash:      hex.EncodeToString(hashBytes[:]),
				Timestamp: time.Now(),
			}
			return
		}

		nonce += nonceStep
	}
}

func validateBlock(block Block, target *big.Int) bool {
	data := fmt.Sprintf("%s%d", block.Data, block.Nonce)
	hashBytes := sha256.Sum256([]byte(data))
	var hashInt big.Int
	hashInt.SetBytes(hashBytes[:])

	// Verify if the hash is less than the target
	return hashInt.Cmp(target) == -1
}

func TestEthBlockMining(t *testing.T) {
	target := new(big.Int).Exp(big.NewInt(2), big.NewInt(256), nil)
	blockData := "This is some block data"

	resultChan := make(chan Block, 1) // Buffer to hold the mined block

	var wg sync.WaitGroup

	numGoroutines := uint64(16) // Changed type to uint64
	startTime := time.Now()     // Record the start time

	for i := uint64(0); i < numGoroutines; i++ { // Changed type of i to uint64
		wg.Add(1)
		go mine(blockData, target, i, numGoroutines, resultChan, &wg)
	}
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	block := <-resultChan

	endTime := time.Now() // Record the end time
	t.Logf("Time taken to mine: %s", endTime.Sub(startTime))
	t.Logf("Target: %d", target)
	t.Logf("Nonce: %d", block.Nonce)
	t.Logf("Hash: %s", block.Hash)

	isValid := validateBlock(block, target)
	if isValid {
		t.Log("Block is valid.")
	} else {
		t.Error("Block is invalid.")
	}
}
