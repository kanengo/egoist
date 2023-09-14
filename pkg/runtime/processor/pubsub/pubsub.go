package pubsub

import (
	"context"
	contribPubsub "github.com/kanengo/egoist/components_contrib/pubsub"
	"github.com/kanengo/egoist/pkg/components"
	"github.com/kanengo/egoist/pkg/runtime/meta"

	"github.com/kanengo/egoist/pkg/components/pubsub"
	"github.com/kanengo/egoist/pkg/resources/components/v1alpha1"
)

type Options struct {
	ID            string
	Namespace     string
	IsHTTP        bool
	PodName       string
	ResourcesPath []string

	Registry *pubsub.Registry
	Meta     *meta.Meta
	//TracingSpec    *config.TracingSpec
	//GRPC           *manager.Manager
	//Channels       *channels.Channels
	//OperatorClient operatorv1.OperatorClient
}

type pubSubManager struct {
	compStore components.CompStore[contribPubsub.PubSub]
	meta      *meta.Meta
	registry  *pubsub.Registry
}

func New(opts Options) *pubSubManager {
	ps := &pubSubManager{
		compStore: components.CompStore[contribPubsub.PubSub]{},
		meta:      opts.Meta,
		registry:  opts.Registry,
	}

	return ps
}

func (p *pubSubManager) Init(ctx context.Context, component v1alpha1.Component) error {
	spec := component.Spec
	ps, err := p.registry.Create(spec.Type, component.ResourceVersion)
	if err != nil {
		return err
	}

	metaBase, err := p.meta.ToBaseMetadata(component)
	if err != nil {
		return err
	}

	err = ps.Init(ctx, contribPubsub.Metadata{Base: metaBase})
	if err != nil {
		return err
	}

	p.compStore.Set(component.Name, ps)

	return nil
}

func (p *pubSubManager) Close(component v1alpha1.Component) error {
	//TODO implement me
	panic("implement me")
}
