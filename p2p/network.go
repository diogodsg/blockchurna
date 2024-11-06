package p2p

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
)


type P2PNode struct {
	Host host.Host
	PubSub *pubsub.PubSub
	Topic *pubsub.Topic
}



func NewP2PNode(ctx context.Context, listenAddr multiaddr.Multiaddr, listenport string) (*P2PNode, error) {
	priv := GetNodeKey(listenport)
    h, err := libp2p.New(libp2p.Identity(priv), libp2p.ListenAddrs(listenAddr))
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

    listenAddrStr := fmt.Sprintf("/ip4/0.0.0.0/tcp/%s", listenPort)
    listenAddr, err := multiaddr.NewMultiaddr(listenAddrStr)
    if err != nil {
        log.Fatalf("Failed to create multiaddress: %v", err)
    }

    node, err := NewP2PNode(ctx, listenAddr, listenPort) 
    if err != nil {
        log.Fatal(err)
    }

	nodeDHT, err := setupDHT(ctx, node.Host)
	if err != nil {
		fmt.Println("Error setting up DHT:", err)
		return nil
	}

    fmt.Printf("Node is listening on: %s/p2p/%s\n", listenAddrStr, node.Host.ID().String())

	ConnectToSeedNodes(ctx, node)

	go discoverPeers(ctx, nodeDHT)

	return node
}

func ConnectToSeedNodes(ctx context.Context, node *P2PNode) {
	file, err := os.Open("seed_nodes")
	if err != nil {
		log.Printf("Error opening the file:", err)
		return
	}
	defer file.Close() 

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		peerAddr := scanner.Text() 
		log.Printf("Connecting to peer: (%s)\n", peerAddr)

        if err := node.ConnectToPeer(ctx, peerAddr); err != nil {
            log.Printf("Failed to connect to peer: (%v)", err)
			continue
        }

        log.Printf("Successfully connected to peer: %s\n", peerAddr)
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Error reading lines: %v", err)
	}
}


