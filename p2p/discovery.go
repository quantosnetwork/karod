package p2p

import (
	"github.com/libp2p/go-libp2p-core/peer"
	"sync"
)

var discoveryProtocolPrefix = "disc/1.0.0"

const defaultBucketSize = 20

type referencePeer struct {
	id     peer.ID
	stream interface{}
}

type referencePeers struct {
	mux   sync.RWMutex
	peers []*referencePeer
}
