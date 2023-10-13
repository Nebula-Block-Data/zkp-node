package ether_mining

import (
	"github.com/ethereum/go-ethereum/common"
	"time"
)

type Block struct {
	Data      string
	Nonce     uint64
	MixDigest common.Hash
	Hash      string
	Timestamp time.Time
}
