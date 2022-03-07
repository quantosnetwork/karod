package p2p

import (
	"context"
	"errors"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	grpc2 "github.com/quantosnetwork/karod/p2p/grpc"
	pb "github.com/quantosnetwork/karod/p2p/proto"
	"go.uber.org/atomic"
	"google.golang.org/grpc"
	"sync"
)

var IdentityProtocolPrefix = "/id/1.0.0"
var (
	ErrInvalidChainID   = errors.New("invalid chain ID")
	ErrNotReady         = errors.New("not ready")
	ErrNoAvailableSlots = errors.New("no available Slots")
)

type identity struct {
	pb.UnimplementedIdentityServer
	pending            sync.Map
	pendingInboundCnt  atomic.Int64
	pendingOutboundCnt atomic.Int64
	srv                *Server
	initialized        atomic.Uint32
}

func (i *identity) updatePendingCount(direction network.Direction, delta int64) {
	switch direction {
	case network.DirInbound:
		i.pendingInboundCnt.Add(delta)
	case network.DirOutbound:
		i.pendingOutboundCnt.Add(delta)
	}
}

func (i *identity) pendingInboundConns() int64 {
	return i.pendingInboundCnt.Load()
}

func (i *identity) pendingOutboundConns() int64 {
	return i.pendingOutboundCnt.Load()
}

func (i *identity) isPending(id peer.ID) bool {
	_, ok := i.pending.Load(id)
	return ok
}

func (i *identity) delPending(id peer.ID) {
	if value, loaded := i.pending.LoadAndDelete(id); loaded {
		direction, ok := value.(network.Direction)
		if !ok {
			return
		}
		i.updatePendingCount(direction, -1)
	}
}

func (i *identity) setPending(id peer.ID, direction network.Direction) {
	if _, loaded := i.pending.LoadOrStore(id, direction); !loaded {
		i.updatePendingCount(direction, 1)
	}
}

func (i *identity) setup() {
	grpc := grpc2.NewGrpcStream()
	pb.RegisterIdentityServer(grpc.GrpcServer(), i)
	grpc.Serve()

	i.srv.Register(IdentityProtocolPrefix, grpc)

	i.srv.h.Network().Notify(&network.NotifyBundle{
		ConnectedF: func(net network.Network, conn network.Conn) {
			peerID := conn.RemotePeer()
			i.srv.logger.Debug("Conn", "peer", peerID, "direction", conn.Stat().Direction)
			initialized := i.initialized.Load()
			if initialized == 0 {
				i.srv.Disconnect(peerID, ErrNotReady.Error())
				return
			}
			if i.isPending(peerID) {
				return
			}
			if i.checkSlots(conn.Stat().Direction, peerID) {
				return
			}
			i.setPending(peerID, conn.Stat().Direction)

			go func() {
				defer func() {
					if i.isPending(peerID) {
						i.delPending(peerID)
						i.srv.emitEvent(peerID, PeerDialCompleted)
					}
				}()
				if err := i.handleConnected(peerID, conn.Stat().Direction); err != nil {
					i.srv.Disconnect(peerID, err)
				}
			}()

		},
	})

}

func (i *identity) start() error {
	i.initialized.Store(1)
	return nil
}

func (i *identity) getStatus() *pb.Status {
	return &pb.Status{
		Chain: int64(i.srv.config.Chain.Params.ChainID),
	}
}

func (i *identity) checkSlots(direction network.Direction, peerID peer.ID) (slotsFull bool) {
	switch direction {
	case network.DirInbound:
		slotsFull = i.srv.inboundConns() >= i.srv.maxInboundConns()
	case network.DirOutbound:
		slotsFull = i.srv.numOpenSlots() == 0
	default:
		i.srv.logger.Info("Invalid connection direction", "peer", peerID)
	}

	if slotsFull {
		i.srv.Disconnect(peerID, ErrNoAvailableSlots.Error())
	}

	return slotsFull
}

func (i *identity) Handshake(ctx context.Context, req *pb.Status) (*pb.Status, error) {
	return i.getStatus(), nil
}

func (i *identity) Bye(ctx context.Context, req *pb.ByeMsg) (*empty.Empty, error) {
	connContext, ok := ctx.(*grpc.Context)
	if !ok {
		return nil, errors.New("invalid type assert")
	}

	i.srv.logger.Debug("peer bye", "id", connContext.PeerID, "msg", req.Reason)

	return &empty.Empty{}, nil
}
