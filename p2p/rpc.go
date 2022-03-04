package p2p

import (
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/protocol"
	gorpc "github.com/libp2p/go-libp2p-gorpc"
	"karod/services"
)

const P2PRpcPrefix protocol.ID = "/quantos/p2p/rpc/1.0.0"

type P2PRpc struct {
	host        *gorpc.Server
	client      *gorpc.Client
	localClient *gorpc.Client
}

func (prpc *P2PRpc) New(h host.Host) {
	prpc.host = gorpc.NewServer(h, P2PRpcPrefix)
	prpc.localClient = gorpc.NewClientWithServer(h, P2PRpcPrefix, prpc.host)

}

func (prpc *P2PRpc) RegisterServices(repo *services.ServiceRepository) (err error) {
	for _, elem := range repo.Services {
		for k, svc := range elem {
			// We register all services in the ServiceRepository by Name
			err = prpc.host.RegisterName(k, svc)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
