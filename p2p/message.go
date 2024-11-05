package p2p

import (
	"blockchurna/blockchain"
	"encoding/json"
)

const (
	MessageTypeBlock = "BLOCK"
	MessageTypeTransaction = "TRANSACTION"
	MessageRequestLatestBlock = "REQUEST_LATEST_BLOCK"
	MessageRequestMissingBlocks = "REQUEST_MISSING_BLOCKS"
	MessageTypeBlocks = "TYPE_BLOCKS"
)

type Message struct {
	Type	string		`json:"type"`
	Payload interface{}	`json:"payload"`
}

func SerializeBlock(block *blockchain.Block) ([]byte, error) {
	message := Message{
		Type: MessageTypeBlock,
		Payload: block,
	}

	return json.Marshal(message)
}