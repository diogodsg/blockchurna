package p2p

import (
	"blockchurna/blockchain"
	"context"
	"encoding/json"
	"fmt"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/mitchellh/mapstructure"
)

func BroadcastBlock(h host.Host, block *blockchain.Block) error {
	data, err := SerializeBlock(block)

	if err != nil {
		return fmt.Errorf("failed to serialize block: %v", err)
	}

	for _, peer := range h.Network().Peers() {
		fmt.Println("Broadcasting block")
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

func StartListening(h host.Host, bc *blockchain.Blockchain) {
	h.SetStreamHandler("/blockchain/1.0.0", func(stream network.Stream) {
		var msg Message
		err := json.NewDecoder(stream).Decode(&msg)
		if err != nil {
			fmt.Printf("failed to decode message: %v\n", err) 
			stream.Reset()
			return
		}
		fmt.Printf("\nreceiving message %s\n\n", msg.Type)

		switch msg.Type {
		case MessageTypeBlock:
            var block blockchain.Block
            err := mapstructure.Decode(msg.Payload, &block)
            if err != nil {
				fmt.Printf("failed to decode block: %v\n", err)
                return
            }
            // Validate and add the block to the local blockchain
			fmt.Printf("\nadding block %s\n\n", block.Payload)
            bc.AddBlock(block.Payload)
        case MessageTypeTransaction:
            var tx string
            err := mapstructure.Decode(msg.Payload, &tx)
            if err != nil {
                fmt.Printf("failed to decode transaction: %v\n", err)
                return
            }
            // Validate and process the transaction
            // You could add it to a mempool here
        
		case MessageRequestLatestBlock:
			var tx string
			err := mapstructure.Decode(msg.Payload, &tx)
			if err != nil {
				fmt.Printf("failed to decode transaction: %v\n", err)
				return
			}
			// Validate and process the transaction
			// You could add it to a mempool here
		}
	})
}

func SynchronizeChain(h host.Host, bc *blockchain.Blockchain) error {
	for _, peer := range h.Network().Peers() {
        stream, err := h.NewStream(context.Background(), peer, "/blockchain/1.0.0")
        if err != nil {
            return fmt.Errorf("failed to create stream: %v", err)
        }
        // Request the peer's latest block
        // You may want to define specific request types (e.g., "REQUEST_BLOCK")
        err = json.NewEncoder(stream).Encode(Message{
            Type: "REQUEST_LATEST_BLOCK",
        })
        if err != nil {
            return fmt.Errorf("failed to send request: %v", err)
        }
        
        // Read the peer's response
        var msg Message
        err = json.NewDecoder(stream).Decode(&msg)
        if err != nil {
            return fmt.Errorf("failed to decode message: %v", err)
        }

        // Compare the received block with the local chain
        if msg.Type == MessageTypeBlock {
            var block blockchain.Block
            err = mapstructure.Decode(msg.Payload, &block)
            if err != nil {
                return fmt.Errorf("failed to decode block: %v", err)
            }
            
            latestLocalBlock := bc.GetLatestBlock()
            if block.Index > latestLocalBlock.Index {
                // Request missing blocks or sync from this point
                // (This part could be implemented depending on your sync strategy)
            }
        }
    }
    return nil
}

func StartBlockchain() {
	block := blockchain.BC.AddBlock("new item")
	p2pNode := ConnectToNetwork() 
	// p2p.SynchronizeChain(p2pNode.Host, bc)
	BroadcastBlock(p2pNode.Host, block)
	StartListening(p2pNode.Host, blockchain.BC)

	for _, block := range blockchain.BC.Blocks {
		fmt.Printf("Index: %d, Payload: %s, Hash: %s\n", block.Index, block.Payload, block.Id)
	}
}
