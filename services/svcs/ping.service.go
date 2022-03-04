package svcs

import (
	"context"
	gorpc "github.com/libp2p/go-libp2p-gorpc"
	"log"
)

type PingArgs struct {
	Data []byte
}
type PingReply struct {
	Data []byte
}
type PingService struct{}

func (t *PingService) Path() string {
	//TODO implement me
	panic("implement me")
}

func (t *PingService) GrpcServer(ctx context.Context, rpc *gorpc.Server) {
	//TODO implement me
	panic("implement me")
}

func (t *PingService) GrpcClient(ctx context.Context, rpc *gorpc.Client) {
	//TODO implement me
	panic("implement me")
}

func (t *PingService) SignedHash() []byte {
	//TODO implement me
	panic("implement me")
}

func (t *PingService) ACL(m ...map[string]string) {
	//TODO implement me
	panic("implement me")
}

func (t *PingService) StartServer() error {
	//TODO implement me
	panic("implement me")
}

func (t *PingService) StartClient() error {
	//TODO implement me
	panic("implement me")
}

func (t *PingService) Close() {
	//TODO implement me
	panic("implement me")
}

func (t *PingService) State() interface{} {
	//TODO implement me
	panic("implement me")
}

func (t *PingService) Ping(ctx context.Context, argType PingArgs, replyType *PingReply) error {
	log.Println("Received a Ping call")
	replyType.Data = argType.Data
	return nil
}

func (t *PingService) Name() string {
	return "ping"
}

func NewPingService() *PingService {
	return &PingService{}
}
