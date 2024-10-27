package p2p

import (
	"fmt"

	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/peerstore"
)

func (node *P2PNode) AddPeer(peerAddr string) error {
	addr, err := peer.AddrInfoFromString(peerAddr)

	if err != nil {
		return fmt.Errorf("failed to parse peer address: %w", err)
	}

	node.Host.Peerstore().AddAddrs(addr.ID, addr.Addrs, peerstore.PermanentAddrTTL)
	fmt.Println("Added peer", addr.ID.String())

	return nil
}

func (node *P2PNode) ListPeers() {
	peers := node.Host.Network().Peers()

	fmt.Println("Connected peers:")

	for _, peerId := range peers {
		fmt.Println("- ", peerId.String())
	}
}