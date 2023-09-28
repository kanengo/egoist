package processor

import (
	"context"

	apiv1 "github.com/kanengo/egoist/pkg/api/v1"
)

func (p *Processor) PublishEvent(ctx context.Context, request *apiv1.PublishEventRequest) (*apiv1.PublishEventResponse, error) {
	return p.pubsub.Publish(ctx, request)
}

func (p *Processor) SubscribeStream(request *apiv1.SubscribeRequest, streamServer apiv1.API_SubscribeStreamServer) error {
	return p.pubsub.SubscribeStream(request, streamServer)
}
