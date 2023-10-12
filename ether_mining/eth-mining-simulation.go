package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
	"sync"
	"time"
)

type Block struct {
	Data      string
	Nonce     int
	Hash      string
	Timestamp time.Time
}

func mine(blockData string, target *big.Int, startNonce int, nonceStep int, resultChan chan<- Block, wg *sync.WaitGroup) {
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

func main() {
	target := new(big.Int).Exp(big.NewInt(2), big.NewInt(230), nil)
	blockData := "This is some block data"

	resultChan := make(chan Block)

	var wg sync.WaitGroup

	numGoroutines := 4

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go mine(blockData, target, i, numGoroutines, resultChan, &wg)
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	block := <-resultChan

	fmt.Printf("Time taken to mine: %s\n", time.Since(block.Timestamp))
	fmt.Printf("Nonce: %d\n", block.Nonce)
	fmt.Printf("Hash: %s\n", block.Hash)

	isValid := validateBlock(block, target)
	if isValid {
		fmt.Println("Block is valid.")
	} else {
		fmt.Println("Block is invalid.")
	}
}
