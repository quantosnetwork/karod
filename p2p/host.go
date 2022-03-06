package p2p

import (
	"context"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/routing"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

func CreateNewHost() (host.Host, *dht.IpfsDHT, error) {

	ctx := context.Background()
	var kdht *dht.IpfsDHT
	routing := libp2p.Routing(func(h host.Host) (routing.PeerRouting, error) {
		kad, err := dht.New(ctx, h)
		kdht = kad
		return kad, err
	})

	h, err := libp2p.New(
		libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/0"),

		libp2p.NATPortMap(),
		routing,
		libp2p.EnableAutoRelay(),
		libp2p.EnableNATService(),
	)
	_, err = pubsub.NewGossipSub(ctx, h)
	if err != nil {
		panic(err)
	}

	return h, kdht, err
}

func init() {

	_, _, _ = CreateNewHost()
}
