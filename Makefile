.PHONY: all protobuf p2p clean

all: protobuf p2p

protobuf: protoc_gen_api

p2p:
	go build -o bin/karod

protoc_gen_api:
	protoc $$PROTO_PATH --go_out=plugins=grpc:. ./api/proto/*.proto
