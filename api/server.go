package api

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	apipb "karod/pb/api"
	"net"
	"os"
)

type GrpcApiServer struct {
	netConnection net.Listener
	server        *grpc.Server
	conns         []*grpc.ClientConn
	streams       []*grpc.ClientStream
}

type apiServer interface {
	GetRouter() interface{}
	GetRPCEndpoint(path string)
	Authorize()
	Authenticate()
	VerifyToken()
	RedirectTo(url string)
	ProxyTo(proxyHost string, proxyPort string)
}

type ApiManager struct {
	Ctx context.Context
	apiServer
}

func (a *ApiManager) NewServer() *GrpcApiServer {
	server := &GrpcApiServer{}
	server.server = grpc.NewServer()
	apipb.RegisterAuthorizationServerServer(server.server, a.NewAuthorizationServer())
	apipb.RegisterAuthenticatorServer(server.server, a.NewAuthenticationServer())
	return server
}

func (a *ApiManager) NewAuthorizationServer() apipb.AuthorizationServerServer {
	return AuthorizationServer{}
}

func (a *ApiManager) NewAuthenticationServer() apipb.AuthenticatorServer {
	return AuthenticationServer{}
}

func (a *ApiManager) StartApiServer(host string, port int) {

	conn, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	a.Ctx = context.Background()
	ctx, cancel := context.WithCancel(a.Ctx)
	defer cancel()
	defer ctx.Done()
	s := a.NewServer()
	s.netConnection = conn
	err = s.server.Serve(s.netConnection)
	if err != nil {
		fmt.Printf("Error serving api server on port %d (error): %v", port, err)
		cancel()
		os.Exit(0)
		return
	}

}
