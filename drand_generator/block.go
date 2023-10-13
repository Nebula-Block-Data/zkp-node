package drand_generator

import (
	"time"
)

type Block struct {
	Data      string
	Nonce     uint64
	Hash      string
	Timestamp time.Time
}
