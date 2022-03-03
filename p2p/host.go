package p2p

import (
	"context"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/routing"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	noise "github.com/libp2p/go-libp2p-noise"
	libp2ptls "github.com/libp2p/go-libp2p-tls"
)

func CreateNewHost(sk crypto.PrivKey) (host.Host, *dht.IpfsDHT, error) {
	host1, _ := libp2p.New()
	ctx := context.Background()
	defer host1.Close()
	var kdht *dht.IpfsDHT
	routing := libp2p.Routing(func(h host.Host) (routing.PeerRouting, error) {
		kad, err := dht.New(ctx, h)
		kdht = kad
		return kad, err
	})

	h, err := libp2p.New(
		libp2p.Identity(sk),
		libp2p.ListenAddrs(host1.Addrs()...),
		libp2p.Security(libp2ptls.ID, libp2ptls.New),
		libp2p.Security(noise.ID, noise.New),
		libp2p.DefaultTransports,
		libp2p.NATPortMap(),
		routing,
		libp2p.EnableAutoRelay(),
		libp2p.EnableNATService(),
	)
	return h, kdht, err
}
