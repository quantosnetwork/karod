package p2p

import (
	"context"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/routing"
	dht "github.com/libp2p/go-libp2p-kad-dht"
)

const DefaultKadPrefix = "/quantos/karod/kad"

var DHTContext = context.Background()

func NewKarodDHT(h host.Host) (routing.PeerRouting, error) {
	kDht, err := dht.New(DHTContext, h, dht.ProtocolPrefix(DefaultKadPrefix))
	return kDht, err
}
