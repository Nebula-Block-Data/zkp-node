package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/drand/drand/client"
	"github.com/drand/drand/client/http"
)

var urls = []string{
	"https://api.drand.sh",
	"https://drand.cloudflare.com",
}

var chainHash, _ = hex.DecodeString("8990e7a9aaed2ffed73dbd7092123d6f289930540d7651336225dc172e51b2ce")

func main() {
	c, err := client.New(
		client.From(http.ForURLs(urls, chainHash)...),
		client.WithChainHash(chainHash),
	)

	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}

	// Fetching randomness
	randomness, err := c.Get(context.Background(), 0)
	if err != nil {
		log.Fatalf("Error getting randomness: %v", err)
	}

	fmt.Printf("Randomness round: %d\n", randomness.Round())
	fmt.Printf("Randomness value: %x\n", randomness.Signature())
}
