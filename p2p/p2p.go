package p2p

import (
	"context"
	"fmt"
	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-datastore"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/routing"
	kaddht "github.com/libp2p/go-libp2p-kad-dht"
	mplex "github.com/libp2p/go-libp2p-mplex"
	noise "github.com/libp2p/go-libp2p-noise"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	tls "github.com/libp2p/go-libp2p-tls"
	yamux "github.com/libp2p/go-libp2p-yamux"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
	"github.com/libp2p/go-tcp-transport"
	ws "github.com/libp2p/go-ws-transport"
	"github.com/multiformats/go-multiaddr"
	"github.com/multiformats/go-multihash"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type mdnsNotifee struct {
	h   host.Host
	ctx context.Context
}

func StartNetwork() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	transports := libp2p.ChainOptions(
		libp2p.Transport(tcp.NewTCPTransport),
		libp2p.Transport(ws.New),
	)

	muxers := libp2p.ChainOptions(
		libp2p.Muxer("/yamux/1.0.0", yamux.DefaultTransport),
		libp2p.Muxer("/mplex/6.7.0", mplex.DefaultTransport),
	)

	security := libp2p.ChainOptions(
		libp2p.Security(tls.ID, tls.New),
		libp2p.Security(noise.ID, noise.New),
	)

	listenAddrs := libp2p.ListenAddrStrings(
		"/ip4/0.0.0.0/tcp/0",
		"/ip4/0.0.0.0/tcp/0/ws",
	)
	var dht *kaddht.IpfsDHT
	newDHT := func(h host.Host) (routing.PeerRouting, error) {
		var err error
		dht, err = kaddht.New(ctx, h)
		return dht, err
	}

	routing2 := libp2p.Routing(newDHT)

	nat := libp2p.ChainOptions(
		libp2p.NATPortMap(),
		libp2p.EnableAutoRelay(nil),
		libp2p.EnableNATService(),
		libp2p.AutoNATServiceRateLimit(30, 10, time.Second*30),
	)

	h, _ := libp2p.New(
		transports,
		listenAddrs,
		muxers,
		security,
		nat,
		routing2,
	)

	dht2 := kaddht.NewDHTClient(ctx, h, datastore.NewMapDatastore())

	peers := h.Network().Peers()
	log.Println(peers)
	peerid := h.ID().String()
	ps, err := pubsub.NewGossipSub(ctx, h)
	if err != nil {
		panic(err)
	}
	topic, err := ps.Join("quantos")
	defer func(topic *pubsub.Topic) {
		_ = topic.Close()
	}(topic)
	sub, err := topic.Subscribe()
	go pubsubHandler(ctx, sub)

	for _, addr := range h.Addrs() {
		fmt.Println("Listening on", addr)
	}

	mdnsService := mdns.NewMdnsService(h, "", &mdnsNotifee{h: h, ctx: ctx})
	if err := mdnsService.Start(); err != nil {
		panic(err)
	}

	log.Println(peerid)
	err = dht.Bootstrap(ctx)
	if err != nil {
		panic(err)
	}
	var DefaultBootstrapPeers []multiaddr.Multiaddr
	for _, s := range []string{
		"/ip4/104.131.131.82/tcp/4001/ipfs/QmaCpDMGvV2BGHeYERUEnRQAwe3N8SzbUtfsmvsqQLuvuJ",
		"/ip4/104.236.179.241/tcp/4001/ipfs/QmSoLPppuBtQSGwKDZT2M73ULpjvfd3aZ6ha4oFGL1KrGM",
		"/ip4/104.236.76.40/tcp/4001/ipfs/QmSoLV4Bbm51jM9C4gDYZQ9Cy3U6aXMJDAbzgu2fzaDs64",
		"/ip4/128.199.219.111/tcp/4001/ipfs/QmSoLSafTMBsPKadTEgaXctDQVcqN88CNLHXMkTNwMKPnu",
		"/ip4/178.62.158.247/tcp/4001/ipfs/QmSoLer265NRgSp2LA3dPaeykiS1J6DifTC88f5uVQKNAd",
	} {
		ma, err := multiaddr.NewMultiaddr(s)
		if err != nil {
			panic(err)
		}
		DefaultBootstrapPeers = append(DefaultBootstrapPeers, ma)
	}

	for _, addr := range DefaultBootstrapPeers {
		iaddr := addr

		peerinfo, _ := peer.AddrInfoFromP2pAddr(iaddr)

		if err := h.Connect(ctx, *peerinfo); err != nil {
			continue
		} else {
			fmt.Println("Connection established: ", *peerinfo)
		}
	}

	fmt.Println("announcing ourselves...")
	ch, _ := multihash.Sum([]byte("meet me here"), multihash.SHA2_256, -1)
	c := cid.NewCidV1(cid.Raw, ch)
	log.Println(c.Hash())
	tctx, _ := context.WithTimeout(ctx, time.Second*90)
	if err := dht.Provide(tctx, c, true); err != nil {
		log.Println(err.Error())
	}
	log.Println("starting karod")

	fmt.Println("searching for other peers...")
	tctx, _ = context.WithTimeout(ctx, time.Second*10)
	peers2, err := dht2.FindProviders(tctx, c)
	if err != nil {
		log.Println(err.Error())
	}
	fmt.Printf("Found %d peers!\n", len(peers2))

	var wg sync.WaitGroup
	for _, peerAddr := range kaddht.GetDefaultBootstrapPeerAddrInfos() {
		peerinfo := peerAddr
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := h.Connect(ctx, peerinfo); err != nil {
				log.Println(err)
			} else {
				log.Println("Connection established with bootstrap node:", peerinfo)
			}
		}()
	}
	wg.Wait()

	donec := make(chan struct{}, 1)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT)

	select {
	case <-stop:
		h.Close()
		os.Exit(0)
	case <-donec:
		h.Close()
	}
	select {}
}

func init() {
	StartNetwork()
}

func pubsubHandler(ctx context.Context, sub *pubsub.Subscription) {

}

func (m *mdnsNotifee) HandlePeerFound(pi peer.AddrInfo) {
	m.h.Connect(m.ctx, pi)
}
