package proxy

import (
	"google.golang.org/grpc"
)

type handler struct {
}

func (h *handler) handler(svr any, stream grpc.ServerStream) error {

	return nil
}
