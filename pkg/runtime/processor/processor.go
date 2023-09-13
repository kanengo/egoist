package processor

import (
	"context"
	apiv1 "github.com/kanengo/egoist/pkg/api/v1"
	"github.com/kanengo/egoist/pkg/resources/components/v1alpha1"
)

type ComponentManger interface {
	Init(ctx context.Context, component v1alpha1.Component) error
	Close(component v1alpha1.Component) error
}

type PubsubManager interface {
	Publish(ctx context.Context, request apiv1.PublishEventRequest) error
	BulkPublish(ctx context.Context, request apiv1.BulkPublishRequest) (apiv1.BulkPublishResponse, error)
}

type Processor struct {
	pubsub PubsubManager
}

func New(options Options) *Processor {
	processor := &Processor{}
	return processor
}

func (p *Processor) UpdateComponent(comp v1alpha1.Component) error {

	return nil
}
