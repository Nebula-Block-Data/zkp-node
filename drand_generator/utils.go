package drand_generator

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/drand/drand/client"
	"github.com/drand/drand/client/http"
	"log"
	"math/big"
	"sync"
	"time"
)

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
