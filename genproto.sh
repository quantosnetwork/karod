#!/bin/bash

function install_tools() {
	go get -d \
		github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
		github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 \
		google.golang.org/protobuf/cmd/protoc-gen-go \
		google.golang.org/grpc/cmd/protoc-gen-go-grpc
}

function generate_proto_and_grpc() {
	find * -name "*.proto" | grep -v "vendor" | grep -v "google" | xargs -n1 \
		protoc -I/usr/local/include -I. -Ivendor \
		-Ithird_party/google/api \
		-Ithird_party/google/protobuf \
		-Ithird_party/protoc-gen-swagger/options \
		 --proto_path=third_party \
		--go_out=pkg/api/v1 --go_opt=paths=source_relative \
		--go-grpc_out=pkg/api/v1 --go-grpc_opt=require_unimplemented_servers=false,paths=source_relative
}

function generate_grpc_gateway() {
	protoc -I/usr/local/include -I. -Ivendor \
			-Ithird_party/google/api \
    		-Ithird_party/google/protobuf \
    		-Ithird_party/protoc-gen-swagger/options \
    		 --proto_path=third_party \
		--grpc-gateway_out=logtostderr=true,paths=source_relative:. api/proto/v1/apirpc.proto
}

function generate_openapi() {
	protoc -I/usr/local/include -I. -Ivendor \
			-Ithird_party/google/api \
    		-Ithird_party/google/protobuf \
    		-Ithird_party/protoc-gen-swagger/options \
    	 --proto_path=third_party \
		--openapiv2_out=logtostderr=true:. \
		api/proto/v1/apirpc.proto
}

#install_tools
generate_proto_and_grpc
generate_grpc_gateway
generate_openapi