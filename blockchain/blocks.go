package blockchain

import (
	"crypto/sha256"
	"fmt"
	"time"
)

type Block struct {
	Id           string `json:"id"`
	Index        int    `json:"index"`
	Timestamp    int64  `json:"timestamp"`
	Payload      string `json:"payload"`
	PreviousNode string `json:"previous_node"`
}

func NewBlock(index int, payload string, previousNode string) *Block {
	block := &Block{
		Index:     index,
		Timestamp: time.Now().Unix(),
		Payload: payload,
		PreviousNode: previousNode,
	}
	block.Id = block.calculateHash()

	return block
}

func (b *Block) calculateHash() string {
	record := string(b.Index) + string(b.Timestamp) + b.PreviousNode + string(b.Payload)
	hash := sha256.New()
	hash.Write([]byte(record))
	return fmt.Sprintf("%x", hash.Sum(nil))
}