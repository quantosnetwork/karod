package p2p

import (
	"context"
	"fmt"
	"github.com/hashicorp/go-hclog"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/event"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	mplex "github.com/libp2p/go-libp2p-mplex"
	noise "github.com/libp2p/go-libp2p-noise"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	tls "github.com/libp2p/go-libp2p-tls"
	yamux "github.com/libp2p/go-libp2p-yamux"
	"github.com/libp2p/go-tcp-transport"
	ws "github.com/libp2p/go-ws-transport"
	"github.com/multiformats/go-multiaddr"
	"github.com/quantosnetwork/Quantos/sdk"
	"net"
	"sync"
)

type ServerConfig struct {
	NoDiscover       bool
	Addr             *net.TCPAddr
	NatAddr          net.IP
	DNS              multiaddr.Multiaddr
	DataDir          string
	MaxPeers         int64
	MaxInboundPeers  int64
	MaxOutboundPeers int64
	Chain            sdk.BlockchainManager
	SecretsManager   sdk.KeyManager
	Metrics          *Metrics
}

func SetDefaultConfig() *ServerConfig {
	return &ServerConfig{
		NoDiscover:       false,
		Addr:             &net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 55655},
		MaxPeers:         40,
		MaxInboundPeers:  32,
		MaxOutboundPeers: 8,
	}
}

type Server struct {
	logger  hclog.Logger
	ps      *pubsub.PubSub
	config  *ServerConfig
	closeCh chan struct{}
	h       host.Host
	addrs   []multiaddr.Multiaddr
	peers   map[peer.ID]*peer.PeerRecord

	peersLock sync.Mutex

	metrics *Metrics

	dialQueue *dialQueue

	identity  *identity
	discovery *discovery

	emitterPeerEvent event.Emitter

	inboundConnCount int64
	ctx              context.Context
	cancel           context.CancelFunc
}

func (s *Server) Disconnect(id peer.ID, s2 string) {

}

type Peer struct {
	srv *Server

	Info peer.AddrInfo

	connDirection network.Direction
}

func NewServer(logger hclog.Logger, config *ServerConfig) (*Server, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	s := &Server{}
	s.config = SetDefaultConfig()
	s.ctx = ctx
	s.cancel = cancel
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

	addrsFactory := func(addrs []multiaddr.Multiaddr) []multiaddr.Multiaddr {
		if config.NatAddr != nil {
			addr, _ := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/%s/tcp/%d", config.NatAddr.String(), config.Addr.Port))

			if addr != nil {
				addrs = []multiaddr.Multiaddr{addr}
			}
		} else if config.DNS != nil {
			addrs = []multiaddr.Multiaddr{config.DNS}
		}

		return addrs
	}

	h, _ := libp2p.New(transports, muxers, security, listenAddrs, libp2p.AddrsFactory(addrsFactory))
	emitter, err := h.EventBus().Emitter(new(PeerEvent))
	if err != nil {
		return nil, err
	}
	s.emitterPeerEvent = emitter
	s.h = h

	s.closeCh = make(chan struct{})

	return s, nil

}

type PeerEventType uint

const (
	PeerConnected        PeerEventType = iota // Emitted when a peer connected
	PeerFailedToConnect                       // Emitted when a peer failed to connect
	PeerDisconnected                          // Emitted when a peer disconnected from node
	PeerAlreadyConnected                      // Emitted when a peer already connected on dial
	PeerDialCompleted                         // Emitted when a peer completed dial
	PeerAddedToDialQueue                      // Emitted when a peer is added to dial queue
)

var peerEventToName = map[PeerEventType]string{
	PeerConnected:        "PeerConnected",
	PeerFailedToConnect:  "PeerFailedToConnect",
	PeerDisconnected:     "PeerDisconnected",
	PeerAlreadyConnected: "PeerAlreadyConnected",
	PeerDialCompleted:    "PeerDialCompleted",
	PeerAddedToDialQueue: "PeerAddedToDialQueue",
}

func (s PeerEventType) String() string {
	name, ok := peerEventToName[s]
	if !ok {
		return "unknown"
	}

	return name
}

type PeerEvent struct {
	// PeerID is the id of the peer that triggered
	// the event
	PeerID peer.ID

	// Type is the type of the event
	Type PeerEventType
}
