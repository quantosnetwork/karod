package libp2p_grpc_protocol

import (
	"context"
	manet "github.com/multiformats/go-multiaddr/net"
	"io"
	"net"
)

type grpcListener struct {
	*GRPCProtocol
	listenerCtx       context.Context
	listenerCtxCancel context.CancelFunc
}

// newGrpcListener builds a new GRPC listener.
func newGrpcListener(proto *GRPCProtocol) net.Listener {
	l := &grpcListener{
		GRPCProtocol: proto,
	}
	l.listenerCtx, l.listenerCtxCancel = context.WithCancel(proto.ctx)
	return l
}

// Accept waits for and returns the next connection to the listener.
func (l *grpcListener) Accept() (net.Conn, error) {
	select {
	case <-l.listenerCtx.Done():
		return nil, io.EOF
	case stream := <-l.streamCh:
		return &streamConn{Stream: stream}, nil
	}
}

// Addr returns the listener's network address.
func (l *grpcListener) Addr() net.Addr {
	listenAddrs := l.host.Network().ListenAddresses()
	if len(listenAddrs) > 0 {
		for _, addr := range listenAddrs {
			if na, err := manet.ToNetAddr(addr); err == nil {
				return na
			}
		}
	}
	return fakeLocalAddr()
}

// Close closes the listener.
// Any blocked Accept operations will be unblocked and return errors.
func (l *grpcListener) Close() error {
	l.listenerCtxCancel()
	return nil
}

func fakeLocalAddr() net.Addr {
	localIp := net.ParseIP("127.0.0.1")
	return &net.TCPAddr{IP: localIp, Port: 0}
}

// fakeRemoteAddr returns a dummy remote address.
func fakeRemoteAddr() net.Addr {
	remoteIp := net.ParseIP("127.1.0.1")
	return &net.TCPAddr{IP: remoteIp, Port: 0}
}
