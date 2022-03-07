package services

import (
	"context"
	gorpc "github.com/libp2p/go-libp2p-gorpc"
	"github.com/quantosnetwork/karod/services/svcs"
	"go.uber.org/atomic"
	"time"
)

type Service interface {
	Name() string
	Path() string
	GrpcServer(ctx context.Context, rpc *gorpc.Server)
	GrpcClient(ctx context.Context, rpc *gorpc.Client)
	SignedHash() []byte
	ACL(...map[string]string)
	StartServer() error
	StartClient() error
	Close()
}

type ServiceRepository struct {
	Services map[string]Service
}

type ServiceState struct {
	Current  string
	Previous string
	Next     string
}

type service struct {
	name             string
	path             string
	rpcServer        *gorpc.Server
	rpcClient        *gorpc.Client
	initialState     string
	state            *ServiceState
	upTime           atomic.Int64
	serviceTime      atomic.Int64
	clientsConnected atomic.Int32
}

func (s service) Name() string {
	return s.name
}

func (s service) Path() string {
	return s.path
}

func (s service) GrpcServer(ctx context.Context, rpc *gorpc.Server) {

}

func (s service) GrpcClient(ctx context.Context, rpc *gorpc.Client) {
	//TODO implement me
	panic("implement me")
}

func (s service) SignedHash() []byte {
	return []byte(s.name)
}

func (s service) ACL(m ...map[string]string) {
	//TODO implement me
	panic("implement me")
}

func (s service) StartServer() error {
	//TODO implement me
	panic("implement me")
}

func (s service) StartClient() error {
	//TODO implement me
	panic("implement me")
}

func (s service) Close() {
	//TODO implement me
	panic("implement me")
}

func (s service) State() *ServiceState {
	return s.state
}

func (s ServiceRepository) serviceBuilder(name, path string) Service {
	svc := &service{}
	svc.name = name
	svc.path = path
	svc.initialState = "new"
	svc.state = new(ServiceState)
	svc.upTime.Store(0)
	svc.serviceTime.Store(time.Now().UnixNano())
	svc.clientsConnected.Store(0)
	return svc
}

func NewServiceRepository() *ServiceRepository {
	return new(ServiceRepository)
}

// Global variable for ServiceRepository
var Repo *ServiceRepository

func init() {
	Repo = NewServiceRepository()
}

func CreateNewService(name, path string, f interface{}) error {
	svc := Repo.serviceBuilder(name, path)

	Repo.Services[name] = svc
	return nil
}

func init() {
	Repo.Services = map[string]Service{}
	Repo.Services["ping"] = svcs.NewPingService()

}
