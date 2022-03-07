package api

import (
	"context"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/quantosnetwork/Quantos/sdk"
	"github.com/quantosnetwork/Quantos/sdk/config"
	"github.com/quantosnetwork/karod/p2p"
	apipb "github.com/quantosnetwork/karod/pkg/api/proto/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type ApiServer struct {
	ctx        context.Context
	chainCfg   *config.ChainConfig
	p2pService *p2p.P2PRpc
	quitCh     chan struct{}
}

func NewApiServer(ctx context.Context, h host.Host) apipb.ApiServiceServer {
	a := &ApiServer{
		ctx:      ctx,
		chainCfg: new(config.ChainConfig),
		quitCh:   make(chan struct{}),
	}
	a.p2pService = new(p2p.P2PRpc)
	a.p2pService.New(h)
	apipb.RegisterApiServiceServer(apipb.UnimplementedApiServiceServer, a)
	return a
}

func (a *ApiServer) GetNodeInfo(context.Context, *apipb.EmptyRequest) (*apipb.NodeInfoResponse, error) {
	info := &apipb.NodeInfoResponse{}
	info.Time = timestamppb.New(time.Now())
	info.CodeVersion = "1"
	info.ServerTime = time.Now().UnixNano()
	info.BuildTime = "20220306172517"
	info.GitHash = "e20b28c5a29a"
	info.Mode = "test"
	info.Network = new(apipb.NetworkInfo)
	info.Network.NetworkName = "Quantos-Karod"
	info.Network.Id = "quantos-test-karod"
	return info, nil

}

func (a *ApiServer) GetChainInfo(context.Context, *apipb.EmptyRequest) (*apipb.ChainInfoResponse, error) {
	return &apipb.ChainInfoResponse{}, nil
}

func (a *ApiServer) Config(h string) {
	var vbuf [2]byte
	copy(vbuf[:], a.Version()[:])
	a.chainCfg.Host = h
	a.chainCfg.Version = vbuf
}

var _ sdk.BlockchainManager

func (a *ApiServer) ChainManager() {

}

func (a *ApiServer) Version() []byte {
	return []byte{sdk.SDKVERSION}
}
