package grpc

import (
	"context"
	"errors"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	manet "github.com/multiformats/go-multiaddr/net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	grpcPeer "google.golang.org/grpc/peer"
	"io"
	"net"
)

type GrpcStream struct {
	ctx        context.Context
	streamCh   chan network.Stream
	grpcServer *grpc.Server
}

func NewGrpcStream() *GrpcStream {
	g := &GrpcStream{
		ctx:        context.Background(),
		streamCh:   make(chan network.Stream),
		grpcServer: grpc.NewServer(grpc.UnaryInterceptor(interceptor)),
	}
	return g
}

type Context struct {
	context.Context
	PeerID peer.ID
}

func interceptor(ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {
	fromContext, _ := grpcPeer.FromContext(ctx)
	addr, ok := fromContext.Addr.(*wrapLibp2pAddr)
	if !ok {
		return nil, errors.New("invalid type assertion")
	}
	ctx2 := &Context{
		ctx,
		addr.id,
	}
	h, err := handler(ctx2, req)

	return h, err
}

type streamConn struct {
	network.Stream
}

func (c streamConn) LocalAddr() net.Addr {
	addr, err := manet.ToNetAddr(c.Stream.Conn().LocalMultiaddr())
	if err != nil {
		return fakeRemoteAddr()
	}

	return &wrapLibp2pAddr{Addr: addr, id: c.Stream.Conn().LocalPeer()}
}

func (c streamConn) RemoteAddr() net.Addr {
	addr, err := manet.ToNetAddr(c.Stream.Conn().RemoteMultiaddr())
	if err != nil {
		return fakeRemoteAddr()
	}

	return &wrapLibp2pAddr{Addr: addr, id: c.Stream.Conn().RemotePeer()}
}

type wrapLibp2pAddr struct {
	id peer.ID
	net.Addr
}

func (g *GrpcStream) Client(stream network.Stream) interface{} {
	return WrapClient(stream)
}

func (g *GrpcStream) Serve() {
	go func() {
		_ = g.grpcServer.Serve(g)
	}()
}

func (g *GrpcStream) Handler() func(network.Stream) {
	return func(stream network.Stream) {
		select {
		case <-g.ctx.Done():
			return
		case g.streamCh <- stream:
		}
	}
}

func (g *GrpcStream) RegisterService(sd *grpc.ServiceDesc, ss interface{}) {
	g.grpcServer.RegisterService(sd, ss)
}

func (g *GrpcStream) GrpcServer() *grpc.Server {
	return g.grpcServer
}

func (g *GrpcStream) Accept() (net.Conn, error) {
	select {
	case <-g.ctx.Done():
		return nil, io.EOF
	case stream := <-g.streamCh:
		return &streamConn{Stream: stream}, nil
	}
}

// Addr implements the net.Listener interface
func (g *GrpcStream) Addr() net.Addr {
	return fakeLocalAddr()
}

func (g *GrpcStream) Close() error {
	return nil
}

// --- conn ---

func WrapClient(s network.Stream) *grpc.ClientConn {
	opts := grpc.WithContextDialer(func(ctx context.Context, peerIdStr string) (net.Conn, error) {
		return &streamConn{s}, nil
	})
	conn, err := grpc.Dial("", grpc.WithTransportCredentials(insecure.NewCredentials()), opts)

	if err != nil {
		// TODO: this should not fail at all
		panic(err)
	}

	return conn
}

var _ net.Conn = &streamConn{}

// fakeLocalAddr returns a dummy local address.
func fakeLocalAddr() net.Addr {
	return &net.TCPAddr{
		IP:   net.ParseIP("127.0.0.1"),
		Port: 0,
	}
}

// fakeRemoteAddr returns a dummy remote address.
func fakeRemoteAddr() net.Addr {
	return &net.TCPAddr{
		IP:   net.ParseIP("127.1.0.1"),
		Port: 0,
	}
}
