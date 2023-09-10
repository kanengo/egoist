package processor

import (
	"context"
	apipubsub "github.com/kanengo/egoist/pkg/api/components/pubsub/v1"
	"github.com/kanengo/egoist/pkg/resources/components/v1alpha1"
)

type ComponentManger interface {
	Init(ctx context.Context, component v1alpha1.Component) error
	Close(component v1alpha1.Component) error
}

type PubsubManager interface {
	Publish(ctx context.Context, request apipubsub.PublishRequest) error
	BulkPublish(ctx context.Context, request apipubsub.BulkPublishRequest) (apipubsub.BulkPublishResponse, error)
}

type Processor struct {
	pubsub PubsubManager
}

func New(options Options) *Processor {
	processor := &Processor{}
	return processor
}
