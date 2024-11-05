package p2p

import (
	"blockchurna/blockchain"
	"encoding/json"
)

const (
	MessageLatestBlock = "LATEST_BLOCK"
	MessageRequestBlockchain = "REQUEST_BLOCKCHAIN"
	MessageSyncronize = "SYNCRONIZE_CHAIN"
)

type Message struct {
	Type	string		`json:"type"`
	Payload interface{}	`json:"payload"`
}

func SerializeBlock(block *blockchain.Block) ([]byte, error) {
	message := Message{
		Type: MessageLatestBlock,
		Payload: block,
	}

	return json.Marshal(message)
}