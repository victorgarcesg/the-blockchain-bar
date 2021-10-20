package database

import (
	"crypto/sha256"
	"encoding/json"
)

type Hash [32]byte

type BlockHeader struct {
	Parent Hash   `json:"parent"`
	Time   uint64 `json:"time"`
}

type Block struct {
	Header BlockHeader `json:"header"`
	TXs    []Tx        `json:"payload"`
}

type BlockFS struct {
	Key   Hash  `json:"key"`
	Value Block `json:"value"`
}

func NewBlock(hash Hash, time uint64, txs []Tx) Block {
	return Block{BlockHeader{hash, time}, txs}
}

func (b Block) Hash() (Hash, error) {
	blockJson, err := json.Marshal(b)
	if err != nil {
		return Hash{}, err
	}

	return sha256.Sum256(blockJson), nil
}
