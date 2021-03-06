package gossip

import (
	"github.com/iotaledger/hive.go/autopeering/peer"
	"github.com/iotaledger/hive.go/autopeering/peer/service"
)

// IsSupported returns whether the peer supports the gossip service.
func IsSupported(p *peer.Peer) bool {
	return p.Services().Get(service.GossipKey) != nil
}

// GetAddress returns the address of the gossip service.
func GetAddress(p *peer.Peer) string {
	gossip := p.Services().Get(service.GossipKey)
	if gossip == nil {
		panic("peer does not support gossip")
	}
	return gossip.String()
}
