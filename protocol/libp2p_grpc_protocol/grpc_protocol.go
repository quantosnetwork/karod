package libp2p_grpc_protocol

import (
	"context"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/protocol"
	"google.golang.org/grpc"
	"karod/api"
)

const Protocol protocol.ID = "/grpc/1.0.0"

type GRPCProtocol struct {
	ctx        context.Context
	host       host.Host
	grpcServer *grpc.Server
	apiServer  *api.GrpcApiServer
	streamCh   chan network.Stream
}

func NewGrpcProtocol(ctx context.Context, h host.Host) *GRPCProtocol {
	s := grpc.NewServer()
	a := new(api.ApiManager)
	as := a.NewServer()
	grpcProtocol := &GRPCProtocol{
		ctx:        ctx,
		host:       h,
		grpcServer: s,
		apiServer:  as,
		streamCh:   make(chan network.Stream),
	}
	h.SetStreamHandler(Protocol, grpcProtocol.HandleStream)
	go s.Serve(newGrpcListener(grpcProtocol))
	return grpcProtocol
}

// GetGRPCServer returns the grpc server.
func (p *GRPCProtocol) GetGRPCServer() *grpc.Server {
	return p.grpcServer
}

// HandleStream handles an incoming stream.
func (p *GRPCProtocol) HandleStream(stream network.Stream) {
	select {
	case <-p.ctx.Done():
		return
	case p.streamCh <- stream:
	}
}
