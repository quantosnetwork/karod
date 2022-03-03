package libp2p_grpc_protocol

import (
	"github.com/libp2p/go-libp2p-core/network"
	manet "github.com/multiformats/go-multiaddr/net"
	"net"
)

type streamConn struct {
	network.Stream
}

func (c *streamConn) LocalAddr() net.Addr {
	addr, err := manet.ToNetAddr(c.Stream.Conn().LocalMultiaddr())
	if err != nil {
		return fakeLocalAddr()
	}
	return addr
}

// RemoteAddr returns the remote address.
func (c *streamConn) RemoteAddr() net.Addr {
	addr, err := manet.ToNetAddr(c.Stream.Conn().RemoteMultiaddr())
	if err != nil {
		return fakeRemoteAddr()
	}
	return addr
}

var _ net.Conn = &streamConn{}
