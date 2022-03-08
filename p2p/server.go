package p2p

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/hashicorp/go-hclog"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/event"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/protocol"
	mplex "github.com/libp2p/go-libp2p-mplex"
	noise "github.com/libp2p/go-libp2p-noise"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	tls "github.com/libp2p/go-libp2p-tls"
	yamux "github.com/libp2p/go-libp2p-yamux"
	"github.com/libp2p/go-tcp-transport"
	ws "github.com/libp2p/go-ws-transport"
	"github.com/multiformats/go-multiaddr"
	"github.com/quantosnetwork/Quantos/sdk"
	"github.com/quantosnetwork/karod/workers"
	lcrypto "github.com/quantosnetwork/quantos-kyber-schnorr-go-libp2p-core/crypto"
	"math/big"
	"net"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

const (
	DefaultDialRatio = 0.2

	DefaultLibp2pPort int = 1478

	MinimumPeerConnections int64 = 1

	MinimumBootNodes int = 1

	// Priority for dial queue
	PriorityRequestedDial uint64 = 1

	PriorityRandomDial uint64 = 10
)

// regex string  to match against a valid dns/dns4/dns6 addr
const DNSRegex = `^/?(dns)(4|6)?/[^-|^/][A-Za-z0-9-]([^-|^/]?)+([\\-\\.]{1}[a-z0-9]+)*\\.[A-Za-z]{2,}(/?)$`

var (
	ErrNoBootnodes  = errors.New("no bootnodes specified")
	ErrMinBootnodes = errors.New("minimum 1 bootnode is required")
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
	peers   map[peer.ID]*Peer

	peersLock sync.Mutex

	metrics *Metrics

	dialQueue *dialQueue

	identity  *identity
	discovery *discovery

	emitterPeerEvent event.Emitter
	protocols        map[string]Protocol
	protocolsLock    sync.Mutex

	joinWatchers     map[peer.ID]chan error
	joinWatchersLock sync.Mutex

	inboundConnCount int64
	ctx              context.Context
	cancel           context.CancelFunc
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

func setupLibp2pKeys() (lcrypto.PrivKey, error) {
	priv, _, err := lcrypto.GenerateKeyPair(lcrypto.KYBER, 256)
	if err != nil {
		return nil, err
	}
	return priv, nil
}

func (s *Server) runDial() {
	notifyCh := make(chan struct{})
	err := s.SubscribeFn(func(evnt *PeerEvent) {
		// Only concerned about the listed event types
		switch evnt.Type {
		case PeerConnected, PeerFailedToConnect, PeerDisconnected, PeerDialCompleted, PeerAddedToDialQueue:
		default:
			return
		}

		select {
		case notifyCh <- struct{}{}:
		default:
		}
	})

	if err != nil {
		s.logger.Error("dial manager failed to subscribe", "err", err)
	}

	jobs := make([]workers.Job, s.numOpenSlots())
	for {
		for i := int64(0); i < s.numOpenSlots(); i++ {
			jobs[i].Descriptor = workers.JobDescriptor{
				ID:    workers.JobID(strconv.FormatInt(i, 10)),
				JType: "p2prundial",
			}
			jobs[i].Descriptor.Metadata = make(map[string]interface{})

			jobs[i].Descriptor.Metadata["peers"] = s.dialQueue.items
			jobs[i].Args = s.dialQueue.items
			jobs[i].ExecFn = func(ctx context.Context, args interface{}) (interface{}, error) {
				for _, addr := range s.dialQueue.items {
					if s.isConnected(addr.ID) {
						s.emitEvent(addr.ID, PeerAlreadyConnected)
					} else {
						if err := s.h.Connect(context.Background(), *addr); err != nil {
							s.logger.Debug("failed to dial", "addr", addr.String(), "err", err)
							s.emitEvent(addr.ID, PeerFailedToConnect)
							return nil, err
						}
					}
				}
				return nil, nil
			}
			wp := workers.New(4)
			wp.GenerateFrom(jobs)
			wp.Run(s.ctx)

		}
		select {
		case <-notifyCh:
		case <-s.closeCh:
			return
		}
	}
}

func (s *Server) numOpenSlots() int64 {
	n := s.maxOutboundConns() - s.outboundConns()
	if n < 0 {
		n = 0
	}

	return n
}

func (s *Server) inboundConns() int64 {
	count := atomic.LoadInt64(&s.inboundConnCount)
	if count < 0 {
		count = 0
	}

	return count + s.identity.pendingInboundConns()
}

func (s *Server) outboundConns() int64 {
	return (s.numPeers() - atomic.LoadInt64(&s.inboundConnCount)) + s.identity.pendingOutboundConns()
}

func (s *Server) maxInboundConns() int64 {
	return s.config.MaxInboundPeers
}

func (s *Server) maxOutboundConns() int64 {
	return s.config.MaxOutboundPeers
}

func (s *Server) isConnected(peerID peer.ID) bool {
	return s.h.Network().Connectedness(peerID) == network.Connected
}

func (s *Server) GetProtocols(peerID peer.ID) ([]string, error) {
	return s.h.Peerstore().GetProtocols(peerID)
}

func (s *Server) GetPeerInfo(peerID peer.ID) peer.AddrInfo {
	return s.h.Peerstore().PeerInfo(peerID)
}

func (s *Server) addPeer(id peer.ID, direction network.Direction) {
	s.logger.Info("Peer connected", "id", id.String())
	s.peersLock.Lock()
	defer s.peersLock.Unlock()

	p := &Peer{
		srv:           s,
		Info:          s.h.Peerstore().PeerInfo(id),
		connDirection: direction,
	}
	s.peers[id] = p

	if direction == network.DirInbound {
		atomic.AddInt64(&s.inboundConnCount, 1)
	}

	s.emitEvent(id, PeerConnected)
	s.metrics.Peers.Set(float64(len(s.peers)))
}

func (s *Server) delPeer(id peer.ID) {
	s.logger.Info("Peer disconnected", "id", id.String())

	s.peersLock.Lock()
	defer s.peersLock.Unlock()

	if peer, ok := s.peers[id]; ok {
		if peer.connDirection == network.DirInbound {
			atomic.AddInt64(&s.inboundConnCount, -1)
		}

		delete(s.peers, id)
	}

	if closeErr := s.h.Network().ClosePeer(id); closeErr != nil {
		s.logger.Error(
			fmt.Sprintf("Unable to gracefully close connection to peer [%s], %v", id.String(), closeErr),
		)
	}

	s.emitEvent(id, PeerDisconnected)
	s.metrics.Peers.Set(float64(len(s.peers)))
}

func (s *Server) Disconnect(peer peer.ID, reason string) {
	if s.h.Network().Connectedness(peer) == network.Connected {
		s.logger.Info(fmt.Sprintf("Closing connection to peer [%s] for reason [%s]", peer.String(), reason))

		if closeErr := s.h.Network().ClosePeer(peer); closeErr != nil {
			s.logger.Error(fmt.Sprintf("Unable to gracefully close peer connection, %v", closeErr))
		}
	}
}

var (
	// Anything below 35s is prone to false timeouts, as seen from empirical test data
	DefaultJoinTimeout   = 40 * time.Second
	DefaultBufferTimeout = DefaultJoinTimeout + time.Second*5
)

func (s *Server) JoinAddr(addr string, timeout time.Duration) error {
	addr0, err := multiaddr.NewMultiaddr(addr)
	if err != nil {
		return err
	}

	addr1, err := peer.AddrInfoFromP2pAddr(addr0)

	if err != nil {
		return err
	}

	return s.Join(addr1, timeout)
}

func (s *Server) Join(addr *peer.AddrInfo, timeout time.Duration) error {
	s.logger.Info("Join request", "addr", addr.String())
	s.addToDialQueue(addr, PriorityRequestedDial)

	if timeout == 0 {
		return nil
	}

	err := s.watch(addr.ID, timeout)

	return err
}

func (s *Server) watch(peerID peer.ID, dur time.Duration) error {
	ch := make(chan error)

	s.joinWatchersLock.Lock()
	if s.joinWatchers == nil {
		s.joinWatchers = map[peer.ID]chan error{}
	}

	s.joinWatchers[peerID] = ch
	s.joinWatchersLock.Unlock()

	select {
	case <-time.After(dur):
		s.joinWatchersLock.Lock()
		delete(s.joinWatchers, peerID)
		s.joinWatchersLock.Unlock()

		return fmt.Errorf("timeout %s %s", s.h.ID(), peerID)
	case err := <-ch:
		return err
	}
}

func (s *Server) runJoinWatcher() error {
	return s.SubscribeFn(func(evnt *PeerEvent) {
		switch evnt.Type {
		// only concerned about PeerConnected, PeerFailedToConnect, and PeerAlreadyConnected
		case PeerConnected, PeerFailedToConnect, PeerAlreadyConnected:
		default:
			return
		}

		// try to find a watcher for this peer
		s.joinWatchersLock.Lock()
		errCh, ok := s.joinWatchers[evnt.PeerID]
		if ok {
			errCh <- nil
			delete(s.joinWatchers, evnt.PeerID)
		}
		s.joinWatchersLock.Unlock()
	})
}

func (s *Server) Close() error {
	err := s.h.Close()
	//s.dialQueue.Close()
	close(s.closeCh)

	return err
}

func (s *Server) NewProtoStream(proto string, id peer.ID) (interface{}, error) {
	s.protocolsLock.Lock()
	defer s.protocolsLock.Unlock()

	p, ok := s.protocols[proto]
	if !ok {
		return nil, fmt.Errorf("protocol not found: %s", proto)
	}

	stream, err := s.NewStream(proto, id)

	if err != nil {
		return nil, err
	}

	return p.Client(stream), nil
}

func (s *Server) NewStream(proto string, id peer.ID) (network.Stream, error) {
	return s.h.NewStream(context.Background(), id, protocol.ID(proto))
}

type Protocol interface {
	Client(network.Stream) interface{}
	Handler() func(network.Stream)
}

func (s *Server) Register(id string, p Protocol) {
	s.protocolsLock.Lock()
	s.protocols[id] = p
	s.wrapStream(id, p.Handler())
	s.protocolsLock.Unlock()
}

func (s *Server) wrapStream(id string, handle func(network.Stream)) {
	s.h.SetStreamHandler(protocol.ID(id), func(stream network.Stream) {
		peerID := stream.Conn().RemotePeer()
		s.logger.Debug("open stream", "protocol", id, "peer", peerID)

		handle(stream)
	})
}

func (s *Server) AddrInfo() *peer.AddrInfo {
	return &peer.AddrInfo{
		ID:    s.h.ID(),
		Addrs: s.addrs,
	}
}

func (s *Server) addToDialQueue(addr *peer.AddrInfo, priority uint64) {
	//s.dialQueue.add(addr, priority)
	s.dialQueue.items = append(s.dialQueue.items, addr)
	s.emitEvent(addr.ID, PeerAddedToDialQueue)
}

func (s *Server) emitEvent(peerID peer.ID, typ PeerEventType) {
	evnt := PeerEvent{
		PeerID: peerID,
		Type:   typ,
	}

	if err := s.emitterPeerEvent.Emit(evnt); err != nil {
		s.logger.Info("failed to emit event", "peer", evnt.PeerID, "type", evnt.Type, "err", err)
	}
}

type Subscription struct {
	sub event.Subscription
	ch  chan *PeerEvent
}

func (s *Subscription) run() {
	// convert interface{} to *PeerEvent channels
	for {
		evnt := <-s.sub.Out()
		if obj, ok := evnt.(PeerEvent); ok {
			s.ch <- &obj
		}
	}
}

func (s *Subscription) GetCh() chan *PeerEvent {
	return s.ch
}

func (s *Subscription) Get() *PeerEvent {
	obj := <-s.ch

	return obj
}

func (s *Subscription) Close() {
	s.sub.Close()
}

// Subscribe starts a PeerEvent subscription
func (s *Server) Subscribe() (*Subscription, error) {
	raw, err := s.h.EventBus().Subscribe(new(PeerEvent))
	if err != nil {
		return nil, err
	}

	sub := &Subscription{
		sub: raw,
		ch:  make(chan *PeerEvent),
	}
	go sub.run()

	return sub, nil
}

// SubscribeFn is a helper method to run subscription of PeerEvents
func (s *Server) SubscribeFn(handler func(evnt *PeerEvent)) error {
	sub, err := s.Subscribe()
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case evnt := <-sub.GetCh():
				handler(evnt)

			case <-s.closeCh:
				sub.Close()

				return
			}
		}
	}()

	return nil
}

// SubscribeCh returns an event of of subscription events
func (s *Server) SubscribeCh() (<-chan *PeerEvent, error) {
	ch := make(chan *PeerEvent)

	var isClosed int32 = 0

	err := s.SubscribeFn(func(evnt *PeerEvent) {
		if atomic.LoadInt32(&isClosed) == 0 {
			ch <- evnt
		}
	})
	if err != nil {
		atomic.StoreInt32(&isClosed, 1)
		close(ch)

		return nil, errors.New(err.Error())
	}

	go func() {
		<-s.closeCh
		atomic.StoreInt32(&isClosed, 1)
		close(ch)
	}()

	return ch, nil
}

func StringToAddrInfo(addr string) (*peer.AddrInfo, error) {
	addr0, err := multiaddr.NewMultiaddr(addr)
	if err != nil {
		return nil, err
	}

	addr1, err := peer.AddrInfoFromP2pAddr(addr0)

	if err != nil {
		return nil, err
	}

	return addr1, nil
}

var (
	// Regex used for matching loopback addresses (IPv4, IPv6, DNS)
	// This regex will match:
	// /ip4/localhost/tcp/<port>
	// /ip4/127.0.0.1/tcp/<port>
	// /ip4/<any other loopback>/tcp/<port>
	// /ip6/<any loopback>/tcp/<port>
	// /dns/foobar.com/tcp/<port>
	loopbackRegex = regexp.MustCompile(
		//nolint:lll
		fmt.Sprintf(`^\/ip4\/127(?:\.[0-9]+){0,2}\.[0-9]+\/tcp\/\d+$|^\/ip4\/localhost\/tcp\/\d+$|^\/ip6\/(?:0*\:)*?:?0*1\/tcp\/\d+$|%s`, DNSRegex),
	)

	dnsRegex = "^/?(dns)(4|6)?/[^-|^/][A-Za-z0-9-]([^-|^/]?)+([\\-\\.]{1}[a-z0-9]+)*\\.[A-Za-z]{2,}(/?)$"
)

// AddrInfoToString converts an AddrInfo into a string representation that can be dialed from another node
func AddrInfoToString(addr *peer.AddrInfo) string {
	// Safety check
	if len(addr.Addrs) == 0 {
		panic("No dial addresses found")
	}

	dialAddress := addr.Addrs[0].String()

	// Try to see if a non loopback address is present in the list
	if len(addr.Addrs) > 1 && loopbackRegex.MatchString(dialAddress) {
		// Find an address that's not a loopback address
		for _, address := range addr.Addrs {
			if !loopbackRegex.MatchString(address.String()) {
				// Not a loopback address, dial address found
				dialAddress = address.String()

				break
			}
		}
	}

	// Format output and return
	return dialAddress + "/p2p/" + addr.ID.String()
}

// MultiAddrFromDNS constructs a multiAddr from the passed in DNS address and port combination
func MultiAddrFromDNS(addr string, port int) (multiaddr.Multiaddr, error) {
	var (
		version string
		domain  string
	)

	match, err := regexp.MatchString(
		dnsRegex,
		addr,
	)
	if err != nil || !match {
		return nil, errors.New("invalid DNS address")
	}

	s := strings.Trim(addr, "/")
	split := strings.Split(s, "/")

	if len(split) != 2 {
		return nil, errors.New("invalid DNS address")
	}

	switch split[0] {
	case "dns":
		version = "dns"
	case "dns4":
		version = "dns4"
	case "dns6":
		version = "dns6"
	default:
		return nil, errors.New("invalid DNS version")
	}

	domain = split[1]

	multiAddr, err := multiaddr.NewMultiaddr(
		fmt.Sprintf(
			"/%s/%s/tcp/%d",
			version,
			domain,
			port,
		),
	)

	if err != nil {
		return nil, errors.New("could not create a multi address")
	}

	return multiAddr, nil
}

func (s *Server) numPeers() int64 {
	s.peersLock.Lock()
	defer s.peersLock.Unlock()

	return int64(len(s.peers))
}

func (s *Server) getRandomBootNode() *peer.AddrInfo {
	randNum, _ := rand.Int(rand.Reader, big.NewInt(int64(len(s.discovery.bootnodes))))

	return s.discovery.bootnodes[randNum.Int64()]
}

func (s *Server) Peers() []*Peer {
	s.peersLock.Lock()
	defer s.peersLock.Unlock()

	peers := make([]*Peer, 0, len(s.peers))
	for _, p := range s.peers {
		peers = append(peers, p)
	}

	return peers
}

// hasPeer checks if the peer is present in the peers list [Thread-safe]
func (s *Server) hasPeer(peerID peer.ID) bool {
	s.peersLock.Lock()
	defer s.peersLock.Unlock()

	_, ok := s.peers[peerID]

	return ok
}
