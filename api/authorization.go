package api

import (
	"context"
	apipb "karod/pb/api"
)

type AuthorizationServer struct{}

func (a AuthorizationServer) Authorize(ctx context.Context, request *apipb.AuthorizationRequest) (*apipb.AuthorizationResponse, error) {
	//TODO implement me
	panic("implement me")
}
