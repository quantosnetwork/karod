package api

import (
	"context"
	apipb "karod/pb/api"
)

type AuthenticationServer struct{}

func (a AuthenticationServer) mustEmbedUnimplementedAuthenticatorServer() {
	//TODO implement me
	panic("implement me")
}

func (a AuthenticationServer) Authenticate(ctx context.Context, request *apipb.AuthenticationRequest) (*apipb.AuthenticationResponse, error) {
	//TODO implement me
	panic("implement me")
}
