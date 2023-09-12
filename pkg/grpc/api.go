package grpc

import (
	"context"
	apiv1 "github.com/kanengo/egoist/pkg/api/v1"
	"io"
	"sync"
	"sync/atomic"
)

type API interface {
	io.Closer
	apiv1.APIServer
}

type api struct {
	apiv1.UnimplementedAPIServer
	closed  atomic.Bool
	closeCh chan struct{}
	wg      sync.WaitGroup
}

func (a *api) Close() error {
	defer a.wg.Wait()
	if a.closed.CompareAndSwap(false, true) {
		close(a.closeCh)
	}
	return nil
}

func (a *api) PublishEvent(ctx context.Context, request *apiv1.PublishEventRequest) (*apiv1.PublishEventResponse, error) {
	//TODO implement me
	panic("implement me")
}

type APIOptions struct {
}

func NewAPI(opts APIOptions) API {
	return &api{
		closeCh: make(chan struct{}),
	}
}
