package p2p

import (
	"blockchurna/blockchain"
	"context"
	"fmt"

	"github.com/libp2p/go-libp2p/core/host"
)

var Node *P2PNode

func BroadcastBlock(h host.Host, block *blockchain.Block) error {
	data, err := SerializeBlock(block)

	if err != nil {
		return fmt.Errorf("failed to serialize block: %v", err)
	}
	fmt.Println(h.Network().Peers())
	for _, peer := range h.Network().Peers() {
		fmt.Printf("Broadcasting block to %s\n", peer)
		stream, err := h.NewStream(context.Background(), peer, "/blockchain/1.0.0")

		if err != nil {
			return fmt.Errorf("failed to create stream: %v", err)
		}

		_, err = stream.Write(data)
		if err != nil {
			return fmt.Errorf("failed to send block: %v", err)
		}
		fmt.Println("Successfuly sent block")
		stream.Close()
	}

	return nil
}


func StartBlockchain() {
	// block := blockchain.BC.CreateBlock("new item")
	Node = ConnectToNetwork() 
	// BroadcastBlock(Node.Host, block)
	// StartListening(Node.Host, blockchain.BC)
	// SynchronizeChain(Node.Host, blockchain.BC)
	// for _, block := range blockchain.BC.Blocks {
	// 	fmt.Printf("Index: %d, Payload: %s, Hash: %s\n", block.Index, block.Payload, block.Id)
	// }
}
