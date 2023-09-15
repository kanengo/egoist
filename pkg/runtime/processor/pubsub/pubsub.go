package pubsub

import (
	"context"
	"sync"

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

type pubsubCompItem struct {
	Component contribPubsub.PubSub
	Resource  v1alpha1.Component
}

type pubSubManager struct {
	compStore components.CompStore[pubsubCompItem]
	meta      *meta.Meta
	registry  *pubsub.Registry

	lock sync.RWMutex
}

func New(opts Options) *pubSubManager {
	ps := &pubSubManager{
		compStore: components.CompStore[pubsubCompItem]{},
		meta:      opts.Meta,
		registry:  opts.Registry,
	}

	return ps
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
