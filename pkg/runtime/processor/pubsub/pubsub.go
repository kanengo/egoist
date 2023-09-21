package pubsub

import (
	"context"
	"sync"

	contribPubsub "github.com/kanengo/egoist/components_contrib/pubsub"
	apiv1 "github.com/kanengo/egoist/pkg/api/v1"
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

type pubsubCompItem struct {
	Component contribPubsub.PubSub
	Resource  v1alpha1.Component
}

type pubSubManager struct {
	//opts      Options
	compStore *components.CompStore[pubsubCompItem]
	meta      *meta.Meta
	registry  *pubsub.Registry

	lock sync.RWMutex
}

func NewManager(opts Options) *pubSubManager {
	pm := &pubSubManager{
		//opts:      opts,
		compStore: components.NewCompStore[pubsubCompItem](),
		meta:      opts.Meta,
		registry:  pubsub.DefaultRegistry,
	}

	return pm
}

func (p *pubSubManager) Publish(ctx context.Context, request *apiv1.PublishEventRequest) (*apiv1.PublishEventResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (p *pubSubManager) BulkPublish(ctx context.Context, request *apiv1.BulkPublishRequest) (apiv1.BulkPublishResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (p *pubSubManager) Init(ctx context.Context, comp v1alpha1.Component) error {
	spec := comp.Spec
	ps, err := p.registry.Create(spec.Type, comp.ResourceVersion)
	if err != nil {
		return err
	}

	metaBase, err := p.meta.ToBaseMetadata(comp)
	if err != nil {
		return err
	}

	err = ps.Init(ctx, contribPubsub.Metadata{Base: metaBase})
	if err != nil {
		return err
	}

	p.lock.Lock()
	defer p.lock.Unlock()
	old, ok := p.compStore.Get(comp.Name)
	if ok && old.Component != nil {
		defer func() {
			_ = old.Component.Close()
		}()
	}

	p.compStore.Set(comp.Name, pubsubCompItem{
		Component: ps,
		Resource:  comp,
	})

	return nil
}

func (p *pubSubManager) Close(comp v1alpha1.Component) error {
	p.lock.Lock()
	defer p.lock.Unlock()

	ps, ok := p.compStore.Get(comp.Name)
	if !ok {
		return nil
	}

	if err := ps.Component.Close(); err != nil {
		return err
	}

	p.compStore.Del(comp.Name)

	return nil
}
