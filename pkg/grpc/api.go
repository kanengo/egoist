package grpc

import (
	"context"
	"io"
	"sync"
	"sync/atomic"

	apiv1 "github.com/kanengo/egoist/pkg/api/v1"
	"github.com/kanengo/egoist/pkg/runtime/processor"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type API interface {
	io.Closer
	apiv1.APIServer
}

type api struct {
	apiv1.UnimplementedAPIServer
	closed    atomic.Bool
	closeCh   chan struct{}
	wg        sync.WaitGroup
	processor *processor.Processor
}

func (a *api) PublishEvent(ctx context.Context, request *apiv1.PublishEventRequest) (*apiv1.PublishEventResponse, error) {
	if request.Topic == "" {
		return nil, status.Error(codes.InvalidArgument, "topic must not empty")
	}
	return a.processor.PublishEvent(ctx, request)
}

func (a *api) SubscribeStream(request *apiv1.SubscribeRequest, streamServer apiv1.API_SubscribeStreamServer) error {
	if len(request.Configs) == 0 {
		return nil
	}

	return a.processor.SubscribeStream(request, streamServer)
}

func (a *api) Close() error {
	defer a.wg.Wait()
	if a.closed.CompareAndSwap(false, true) {
		close(a.closeCh)
	}
	return nil
}

type APIOptions struct {
	Processor *processor.Processor
}

func NewAPI(opts APIOptions) API {
	return &api{
		closeCh:   make(chan struct{}),
		processor: opts.Processor,
	}
}
