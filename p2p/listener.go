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
func StartListening(h host.Host, bc *blockchain.Blockchain) {
	h.SetStreamHandler("/blockchain/1.0.0", func(stream network.Stream) {
		fmt.Println("Connection received from:", stream.Conn().RemotePeer())
		fmt.Println("Local Peer ID:", stream.Conn().LocalPeer())

		defer stream.Close() // Ensure the stream is properly closed when done

		 // Read raw data from the stream
		 var rawMessage json.RawMessage
		 if err := json.NewDecoder(stream).Decode(&rawMessage); err != nil {
			 fmt.Println("Error reading raw message:", err)
			 return
		 }
		 fmt.Println("Raw JSON received:", string(rawMessage))

		 var msg Message
		if err := json.Unmarshal(rawMessage, &msg); err != nil {
			fmt.Println("Error decoding message:", err)
			return
		}

		fmt.Printf("\nreceiving message: %s\n", msg.Type)

		switch msg.Type {
		case MessageTypeBlock:
			fmt.Println("Receiving Latest Block")
			var block *blockchain.Block
			err := mapstructure.Decode(msg.Payload, &block)
			if err != nil {
				fmt.Printf("failed to decode block: %v\n", err)
				stream.Reset()
				return
			}
			// Validate and add the block to the local blockchain
			fmt.Println(block)

			bc.AddBlock(block)

		case MessageTypeTransaction:
			var tx string
			err := mapstructure.Decode(msg.Payload, &tx)
			if err != nil {
				fmt.Printf("failed to decode transaction: %v\n", err)
				stream.Reset()
				return
			}
			// Validate and process the transaction (e.g., add it to a mempool)
			fmt.Printf("received transaction: %s\n", tx)
			// Process transaction here

		case MessageRequestLatestBlock:
			// Respond with the latest block in the blockchain
			latestBlock := bc.GetLatestBlock()
			err := json.NewEncoder(stream).Encode(Message{
				Type:    MessageTypeBlock,
				Payload: latestBlock,
			})
			if err != nil {
				fmt.Printf("failed to send latest block: %v\n", err)
				stream.Reset()
				return
			}
			fmt.Println("sent latest block")

		case MessageRequestMissingBlocks:
			fmt.Println("Sending to Remote Peer ID:", stream.Conn().RemotePeer())

			var fromIndex int
			err := mapstructure.Decode(msg.Payload, &fromIndex)
			if err != nil {
				fmt.Printf("Failed to decode block index: %v\n", err)
				return
			}
			// Retrieve missing blocks and encode response
			missingBlocks := bc.GetBlocksAfterIndex(0)

			jsonData, err := json.Marshal(missingBlocks)
			if err != nil {
				fmt.Printf("Error serializing missingBlocks: %v\n", err)
				return
			}
			fmt.Println("Serialized missingBlocks:", string(jsonData))

			responseMessage := Message{
				Type:    MessageTypeBlocks,
				Payload: missingBlocks,
			}
			fmt.Println("Sending MessageTypeBlocks to:", stream.Conn().RemotePeer())
			stream, err := h.NewStream(context.Background(), stream.Conn().RemotePeer(), "/blockchain/1.0.0")
			defer stream.Close()
			if err != nil {
				fmt.Errorf("failed to create stream: %v", err)
				return
			}
			data, _ := json.Marshal(responseMessage)
			_, err = stream.Write(data)
			if err != nil {
				fmt.Errorf("failed to send block: %v", err)
				return 
			}		
			fmt.Println("Sent missing blocks:", missingBlocks)
		
		case MessageTypeBlocks:
			fmt.Println("Receiving blocks to sync")
			var receivedBlocks []*blockchain.Block
			err := mapstructure.Decode(msg.Payload, &receivedBlocks)
			if err != nil {
				fmt.Printf("Failed to decode received blocks: %v\n", err)
				return
			}
			bc.ReplaceBlockchain(receivedBlocks)
		
		default:
			fmt.Printf("received unknown message type: %s\n", msg.Type)
		}
	})
}
