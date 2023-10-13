package drand_generator

import (
	"context"
	"encoding/hex"
	"testing"

	"github.com/drand/drand/client"
	"github.com/drand/drand/client/http"
	"github.com/stretchr/testify/assert"
)

var urls = []string{
	"https://api.drand.sh",
	"https://drand.cloudflare.com",
}

var chainHash, _ = hex.DecodeString("8990e7a9aaed2ffed73dbd7092123d6f289930540d7651336225dc172e51b2ce")

func TestFetchRandomness(t *testing.T) {
	c, err := client.New(
		client.From(http.ForURLs(urls, chainHash)...),
		client.WithChainHash(chainHash),
	)
	assert.NoError(t, err, "Error creating client")

	// Fetching randomness
	randomness, err := c.Get(context.Background(), 0)
	assert.NoError(t, err, "Error getting randomness")

	t.Logf("Randomness round: %d\n", randomness.Round())
	t.Logf("Randomness value: %x\n", randomness.Signature())
}
