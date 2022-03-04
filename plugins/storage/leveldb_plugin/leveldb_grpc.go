package plugins

import (
	"github.com/google/uuid"
	"github.com/hashicorp/go-plugin"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	proto "karod/pb/plugins"
	"time"
)

type LevelDBGrpc struct {
	plugin.Plugin
	Impl LevelDBPluginInterface
}

type LevelDBGrpcClient struct {
	client proto.LevelDBGrpcClient
}

type LevelDBGrpcServer struct {
	Impl LevelDBPluginInterface
}

func (p *LevelDBGrpc) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	proto.RegisterLevelDBGrpcServer(s, &LevelDBGrpcServer{Impl: p.Impl})
	return nil
}

func (p *LevelDBGrpc) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &LevelDBGrpcClient{client: proto.NewLevelDBGrpcClient(c)}, nil
}

func (l *LevelDBGrpcClient) Put(key []byte, value []byte) (*proto.PutResponse, error) {
	res, err := l.client.Put(context.Background(), &proto.PutRequest{
		Key: key, Value: value,
	})
	return res, err
}

func (l *LevelDBGrpcClient) Get(key []byte) (*proto.GetResponse, error) {
	res, err := l.client.Get(context.Background(), &proto.GetRequest{Key: key})

	return res, err
}

func (ls *LevelDBGrpcServer) Put(ctx context.Context, key []byte, value []byte) (*proto.PutResponse, error) {
	err := ls.Impl.Put(string(key), string(value))
	if err != nil {
		return &proto.PutResponse{Entry: key, Success: false, Error: "an error occured putting the data", Timestamp: time.Now().UnixNano()}, nil
	}
	return &proto.PutResponse{Entry: key, Success: true, Error: "", Timestamp: time.Now().UnixNano()}, nil
}

func (ls *LevelDBGrpcServer) Get(ctx context.Context, key []byte) (*proto.GetResponse, error) {
	entry := ls.Impl.Get(string(key))
	res := &proto.GetResponse{}
	res.Entry = &proto.Entry{}
	id, _ := uuid.NewUUID()
	res.Entry.Id = id.String()
	res.Entry.Key = res.Entry.Id
	res.Entry.IsEncrypted = false
	res.Entry.Prefix = ""
	res.Entry.Value = string(entry.Value)
	res.Success = true
	res.Error = ""
	return res, nil
}
