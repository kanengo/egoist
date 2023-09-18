#!/usr/bin/env bash

PROTO_PREFIX=github.com/kanengo/egoist

protoc --proto_path=protos/3rd --proto_path=protos --go_out=. --go_opt=module="$PROTO_PREFIX" --go-grpc_out=. --go-grpc_opt=module="$PROTO_PREFIX" protos/api/v1/*.proto

#protoc --proto_path=protos --go_out=. --go_opt=module="$PROTO_PREFIX" --go-grpc_out=. --go-grpc_opt=module="$PROTO_PREFIX" protos/api/v1/*.proto