package processor

import (
	"context"

	"github.com/kanengo/egoist/pkg/components/pubsub"
	"github.com/kanengo/egoist/pkg/resources/components/v1alpha1"
)

type pubSubManager struct {
}

func (p *pubSubManager) Init(ctx context.Context, component v1alpha1.Component) error {
	spec := component.Spec
	pub, err := pubsub.DefaultRegistry.Create(spec.Type, component.ResourceVersion)
	if err != nil {
		return err
	}
}

func (p *pubSubManager) Close(component v1alpha1.Component) error {
	//TODO implement me
	panic("implement me")
}
