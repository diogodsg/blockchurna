package p2p

import (
	"blockchurna/blockchain"
	"context"
	"encoding/json"
	"fmt"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
)

func SynchronizeChain(h host.Host, bc *blockchain.Blockchain) error {
	fmt.Println("peers: ")

	fmt.Println(h.Network().Peers())
	for _, peer := range h.Network().Peers() {
		fmt.Println("Requesting block from " + peer.String())
		latestLocalBlock := bc.GetLatestBlock()

		err := requestMissingBlocks(h, peer, latestLocalBlock.Index)
		if err != nil {
			return fmt.Errorf("failed to synchronize missing blocks: %v", err)
		}

	}
	return nil
}

func requestMissingBlocks(h host.Host, peer peer.ID, fromIndex int) error {
	fmt.Println("requesting missing blocks")
	fmt.Println(peer)

	stream, err := h.NewStream(context.Background(), peer, "/blockchain/1.0.0")
	if err != nil {
		return fmt.Errorf("failed to create stream: %v", err)
	}

	err = json.NewEncoder(stream).Encode(Message{
		Type:    MessageRequestBlockchain,
		Payload: fromIndex,
	})
	if err != nil {
		return fmt.Errorf("failed to send missing blocks request: %v", err)
	}

	return nil
}
