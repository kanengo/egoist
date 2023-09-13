#!/usr/bin/env bash

PROTO_PREFIX=github.com/kanengo/egoist

protoc --proto_path=. --go_out=. --go_opt=module="$PROTO_PREFIX" --go-grpc_out=. --go-grpc_opt=module="$PROTO_PREFIX" ./api/v1/*.proto

#protoc --proto_path=. --go_out=. --go_opt=module=github.com/kanengo/egoist --go-grpc_out=. --go-grpc_opt=module=github.com/kanengo/egoist ./api/v1/*.proto