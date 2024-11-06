package p2p

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	discovery "github.com/libp2p/go-libp2p/p2p/discovery/routing"
)	


func setupDHT(ctx context.Context, nodeHost host.Host) (*dht.IpfsDHT, error) {
	// Create a new DHT for peer discovery
	dht, err := dht.New(ctx, nodeHost)
	if err != nil {
		return nil, err
	}
	// Bootstrap the DHT for finding other peers
	err = dht.Bootstrap(ctx)
	if err != nil {
		return nil, err
	}
	
	var wg sync.WaitGroup
	for _, peerAddr := range []string{"/ip4/127.0.0.1/tcp/4001/p2p/12D3KooWSDJLc5qLQqXtt3hvx1mie5JPH2LkeGgvbFSHD4WBZFMn"} {
		peerinfo, _ := peer.AddrInfoFromString(peerAddr)

		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := nodeHost.Connect(ctx, *peerinfo); err != nil {
				log.Printf("Error while connecting to node %q: %-v", peerinfo, err)
			} else {
				log.Printf("Connection established with bootstrap node: %q", *peerinfo)
			}
		}()
	}
	wg.Wait()
	
	fmt.Println("DHT initialized and bootstrapped.")
	return dht, nil
}

func discoverPeers(ctx context.Context, dht *dht.IpfsDHT) {
	// Create a new RoutingDiscovery instance
	routingDiscovery := discovery.NewRoutingDiscovery(dht)
	
	_, err := routingDiscovery.Advertise(ctx, "p2p-discovery")
	if err != nil {
		fmt.Println("Error advertising:", err)
		return
		
	}
	fmt.Println("Node advertised in the network")

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