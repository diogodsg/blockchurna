package p2p

import (
	"blockchurna/blockchain"
	"context"
	"encoding/json"
	"fmt"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/mitchellh/mapstructure"
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

    // Send request for missing blocks
    err = json.NewEncoder(stream).Encode(Message{
        Type:    "REQUEST_MISSING_BLOCKS",
        Payload: fromIndex,
    })
    if err != nil {
        return fmt.Errorf("failed to send missing blocks request: %v", err)
    }

    // Read and process the response with the missing blocks
    var msg Message
    err = json.NewDecoder(stream).Decode(&msg)
    if err != nil {
        return fmt.Errorf("failed to decode missing blocks: %v", err)
    }

    if msg.Type == MessageTypeBlocks {
        var missingBlocks []*blockchain.Block
        err = mapstructure.Decode(msg.Payload, &missingBlocks)
        if err != nil {
            return fmt.Errorf("failed to decode missing blocks: %v", err)
        }

        // Append the missing blocks to the local blockchain
        for _, block := range missingBlocks {
            err := blockchain.BC.AddBlock(block)
            if err != nil {
                return fmt.Errorf("failed to add block to blockchain: %v", err)
            }
        }
    }

    return nil
}
