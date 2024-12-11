package p2p

import (
	"context"
	"fmt"
	"time"

	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/host"
	discovery "github.com/libp2p/go-libp2p/p2p/discovery/routing"
)	


func setupDHT(ctx context.Context, nodeHost host.Host) (*dht.IpfsDHT, error) {
	// Create a new DHT for peer discovery
	dht, err := dht.New(ctx, nodeHost, dht.Mode(dht.ModeServer))
	if err != nil {
		return nil, err
	}
	err = dht.Bootstrap(ctx)
	if err != nil {
		return nil, err
	}
	
	fmt.Println("DHT initialized and bootstrapped.")
	return dht, nil
}

func discoverPeers(ctx context.Context, dht *dht.IpfsDHT) {
	routingDiscovery := discovery.NewRoutingDiscovery(dht)
	
	for {
		_, err := routingDiscovery.Advertise(ctx, "p2p-discovery")
		if err != nil {
			fmt.Println("No peers found in DHT. Retrying advertisement...")
			time.Sleep(5 * time.Second) 
			continue
		}
		break
	}
	fmt.Println("Successfully entered the network")

	// Find other peers in the network
	peerChan, err := routingDiscovery.FindPeers(ctx, "p2p-discovery")
	if err != nil {
		fmt.Println("Error finding peers:", err)
		return
	}

	for peerInfo := range peerChan {
		if peerInfo.ID == dht.Host().ID() {
			// Skip connecting to self
			continue
		}

		// Attempt to connect to each discovered peer
		err := dht.Host().Connect(ctx, peerInfo)
		if err != nil {
			fmt.Printf("Failed to connect to peer %s: %v\n", peerInfo.ID, err)
		} else {
			fmt.Println("Connected to peer:", peerInfo.ID)
		}
		time.Sleep(1 * time.Second) // Optional: limit discovery rate
	}
}

func FindPeers(ctx context.Context, host host.Host) {
	nodeDHT, err := setupDHT(ctx, host)
	if err != nil {
		fmt.Println("Error setting up DHT:", err)
		return
	}

	go discoverPeers(ctx, nodeDHT)

	select {}
}