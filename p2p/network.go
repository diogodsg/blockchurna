package p2p

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/libp2p/go-libp2p" // Core libp2p library
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"

	// Network interface
	// Peer management
	// Connection manager
	pubsub "github.com/libp2p/go-libp2p-pubsub" // PubSub for gossiping messages
	// Multiaddr for peer addresses
)


type P2PNode struct {
	Host host.Host
	PubSub *pubsub.PubSub
	Topic *pubsub.Topic
}

func NewP2PNode(ctx context.Context, listenAddr multiaddr.Multiaddr) (*P2PNode, error) {
    h, err := libp2p.New(libp2p.ListenAddrs(listenAddr))
    if err != nil {
        return nil, fmt.Errorf("failed to create libp2p host: %w", err)
    }

    fmt.Printf("Node created with Peer ID: %s\n", h.ID().String())

    ps, err := pubsub.NewGossipSub(ctx, h)
    if err != nil {
        return nil, fmt.Errorf("failed to create pubsub: %w", err)
    }

    topic, err := ps.Join("blockchurna-topic")
    if err != nil {
        return nil, fmt.Errorf("failed to join topic: %w", err)
    }

    fmt.Printf("Node is now subscribed to topic 'blockchurna-topic'\n")

    return &P2PNode{
        Host:   h,
        PubSub: ps,
        Topic:  topic,
    }, nil
}


func (node *P2PNode) ConnectToPeer(ctx context.Context, peerAddr string) error {
	addr, err := multiaddr.NewMultiaddr(peerAddr)

	if err != nil {
		return fmt.Errorf("failed to parse multiaddr: %w", err)
	}

	peerInfo, err := peer.AddrInfoFromP2pAddr(addr)	
	if err != nil {
		return fmt.Errorf("failed to get peer info: %w", err)
	}

	err = node.Host.Connect(ctx, *peerInfo);

	if err != nil {
		return fmt.Errorf("failed to connect to peer: %w", err)
	}

	fmt.Printf("Connected to peer: %s\n", peerInfo.ID.String())
	return nil
}

func (node *P2PNode) BroadcastMessage(ctx context.Context, msg []byte) error {
	return node.Topic.Publish(ctx, msg)
}


func (node *P2PNode) ListenMessages(ctx context.Context, handler func([]byte)) error {
	sub, err := node.Topic.Subscribe()
	
	if err != nil {
		return fmt.Errorf("failed to subscribe to topic: %w", err)
	}

	go func() {
		for {
			msg, err := sub.Next(ctx)

			if err != nil {
				log.Println("Error receiving message: ", err)
				return
			}

			handler(msg.Data)
		}
	}()

	return nil
}


func ConnectToNetwork() *P2PNode {
    ctx := context.Background()

    if len(os.Args) < 2 {
        fmt.Println("Usage: go run main.go <port> [<peer-multiaddr>]")
        return nil
    }

    listenPort := os.Args[1]

    listenAddrStr := fmt.Sprintf("/ip4/127.0.0.1/tcp/%s", listenPort)
    listenAddr, err := multiaddr.NewMultiaddr(listenAddrStr)
    if err != nil {
        log.Fatalf("Failed to create multiaddress: %v", err)
    }

    node, err := NewP2PNode(ctx, listenAddr) 
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Node is listening on: %s/p2p/%s\n", listenAddrStr, node.Host.ID().String())

    if len(os.Args) > 2 {
        peerAddr := os.Args[2]
        fmt.Printf("Connecting to peer: %s\n", peerAddr)

        if err := node.ConnectToPeer(ctx, peerAddr); err != nil {
            log.Fatalf("Failed to connect to peer: %v", err)
        }

        fmt.Println("Successfully connected to peer!")
    }

    return node
}
