package v1

import (
	"context"
	v1 "github.com/quantosnetwork/karod/pkg/api/v1/api/proto/v1"
)

const (
	apiVersion = "v1"
)

type apiServiceServer struct {
}

func (a apiServiceServer) SendRequest(ctx context.Context, request *v1.ApiRequest) (*v1.ApiResponse, error) {
	//TODO implement me
	panic("implement me")
}

func NewApiServiceServer() v1.ApiTestServiceServer {
	return &apiServiceServer{}
}
