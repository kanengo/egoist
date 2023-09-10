#!/usr/bin/env bash

PROTO_PREFIX=github.com/kanengo/egoist

protoc --go_out=. --go_opt=module="$PROTO_PREFIX" --go-grpc_out=. --go-grpc_opt=module="$PROTO_PREFIX" .\api\components\v1\*.proto