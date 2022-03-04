.PHONY: all protobuf p2p clean

all: protobuf p2p

protobuf: protoc_gen_api protoc_gen_grpc_plugins

p2p:
	go build -o bin/karod

protoc_gen_api:
	protoc $$PROTO_PATH --go_out=:. --go-grpc_out=:. ./api/proto/*.proto

protoc_gen_grpc_plugins:
	protoc $$PROTO_PATH --go_out=:. --go-grpc_out=:. ./plugins/storage/leveldb_plugin/proto/*.proto