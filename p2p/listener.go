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

		defer stream.Close()

		msg, err := DecodeMessage(stream)

		if err != nil {
			fmt.Println("Error decoding message")
			return
		}

		fmt.Printf("\nReceived message: %s\n", msg.Type)

		switch msg.Type {
		case MessageLatestBlock:
			HandleMessageLatestBlock(stream, msg, bc)

		case MessageRequestBlockchain:
			HandleMessageRequestBlockchain(h, stream, msg, bc)

		case MessageSyncronize:
			HandleMessageSyncronize(msg, bc)
		default:
			fmt.Printf("received unknown message type: %s\n", msg.Type)
		}
	})
}

func DecodeMessage(stream network.Stream) (*Message, error) {
	var rawMessage json.RawMessage
	if err := json.NewDecoder(stream).Decode(&rawMessage); err != nil {
		fmt.Println("Error reading raw message:", err)
		return nil, err
	}
	fmt.Println("Raw JSON received:", string(rawMessage))

	var msg Message
	if err := json.Unmarshal(rawMessage, &msg); err != nil {
		fmt.Println("Error decoding message:", err)
		return nil, err
	}

	return &msg, nil
}

func HandleMessageLatestBlock(stream network.Stream, msg *Message, bc *blockchain.Blockchain) error {
	fmt.Println("Receiving Latest Block")
	var block *blockchain.Block
	err := mapstructure.Decode(msg.Payload, &block)
	if err != nil {
		fmt.Printf("failed to decode block: %v\n", err)
		stream.Reset()
		return err
	}

	bc.AddBlock(block)

	return nil
}

func HandleMessageRequestBlockchain(h host.Host, stream network.Stream, msg *Message, bc *blockchain.Blockchain) error {
	fmt.Println("Sending to Remote Peer ID:", stream.Conn().RemotePeer())

	var fromIndex int
	err := mapstructure.Decode(msg.Payload, &fromIndex)
	if err != nil {
		return fmt.Errorf("failed to decode block index: %v", err)
	}

	missingBlocks := bc.GetBlocksAfterIndex(0)

	jsonData, err := json.Marshal(missingBlocks)
	if err != nil {
		return fmt.Errorf("error serializing missingBlocks: %v", err)
	}
	fmt.Println("Serialized missingBlocks:", string(jsonData))

	responseMessage := Message{
		Type:    MessageSyncronize,
		Payload: missingBlocks,
	}
	fmt.Println("Sending MessageTypeBlocks to:", stream.Conn().RemotePeer())
	receiverStream, err := h.NewStream(context.Background(), stream.Conn().RemotePeer(), "/blockchain/1.0.0")
	defer receiverStream.Close()

	if err != nil {
		return fmt.Errorf("failed to create stream: %v", err)

	}
	data, _ := json.Marshal(responseMessage)
	_, err = receiverStream.Write(data)
	if err != nil {
		return fmt.Errorf("failed to send block: %v", err)
	}
	fmt.Println("Sent missing blocks:", missingBlocks)

	return nil
}

func HandleMessageSyncronize(msg *Message, bc *blockchain.Blockchain) error {
	fmt.Println("Receiving blocks to sync")
	var receivedBlocks []*blockchain.Block
	err := mapstructure.Decode(msg.Payload, &receivedBlocks)
	if err != nil {
		return fmt.Errorf("failed to decode received blocks: %v", err)
	}
	bc.ReplaceBlockchain(receivedBlocks)

	return nil
}
