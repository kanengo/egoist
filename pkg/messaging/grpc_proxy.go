package messaging

import (
	"context"

	"google.golang.org/grpc"
)

type Proxy interface {
	Handler() grpc.StreamHandler
	SetRemoteAppFn(func(string) (remoteApp, error))
	SetTelemetryFn(func(context.Context) context.Context)
}

type remoteApp struct {
	id        string
	namespace string
	address   string
}
