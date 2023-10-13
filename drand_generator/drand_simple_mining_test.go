package drand_generator

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"sync"
	"testing"
	"time"

	"github.com/drand/drand/client"
	"github.com/drand/drand/client/http"
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

func getDrandRandomness() []byte {
	urls := []string{
		"https://api.drand.sh",
		"https://drand.cloudflare.com",
	}

	chainHash, _ := hex.DecodeString("8990e7a9aaed2ffed73dbd7092123d6f289930540d7651336225dc172e51b2ce")

	c, err := client.New(
		client.From(http.ForURLs(urls, chainHash)...),
		client.WithChainHash(chainHash),
	)

	if err != nil {
		log.Fatalf("Error creating drand client: %v", err)
	}

	randomness, err := c.Get(context.Background(), 0)
	if err != nil {
		log.Fatalf("Error getting randomness from drand: %v", err)
	}

	return randomness.Signature()
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
}
