package processor

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"

	compPubsub "github.com/kanengo/egoist/pkg/components/pubsub"
	"github.com/kanengo/egoist/pkg/runtime/meta"
	"github.com/kanengo/egoist/pkg/runtime/processor/pubsub"
	"github.com/kanengo/goutil/pkg/log"
	"go.uber.org/zap"

	"github.com/kanengo/egoist/components_contrib"
	apiv1 "github.com/kanengo/egoist/pkg/api/v1"
	"github.com/kanengo/egoist/pkg/resources/components/v1alpha1"
)

type componentManger interface {
	Init(ctx context.Context, component v1alpha1.Component) error
	Close(component v1alpha1.Component) error
}

type PubsubManager interface {
	Publish(ctx context.Context, request *apiv1.PublishEventRequest) (*apiv1.PublishEventResponse, error)
	BulkPublish(ctx context.Context, request *apiv1.BulkPublishRequest) (apiv1.BulkPublishResponse, error)
}

type componentMangerItem struct {
	componentManger
	mu    sync.RWMutex
	comps map[string]v1alpha1.Component
}

type Processor struct {
	pubsub       PubsubManager
	compManagers map[string]*componentMangerItem
}

func New(options Options) *Processor {
	p := &Processor{
		compManagers: make(map[string]*componentMangerItem),
	}

	p.compManagers[components_contrib.TypePubsub] = &componentMangerItem{
		componentManger: pubsub.NewManager(pubsub.Options{
			ID:            options.ID,
			Namespace:     options.NameSpace,
			PodName:       options.PodName,
			ResourcesPath: nil,
			Registry:      compPubsub.DefaultRegistry,
			Meta:          &meta.Meta{},
		}),
		comps: make(map[string]v1alpha1.Component, 3),
	}

	return p
}

func parseComponentType(comp v1alpha1.Component) (string, error) {
	typ := comp.Spec.Type
	if typ == "" {
		return "", errors.New("component type must have value")
	}

	ct := strings.Split(typ, ".")[0]

	return ct, nil
}

func (p *Processor) InitComponent(ctx context.Context, comp v1alpha1.Component) error {
	log.Debug("InitComponent", zap.Any("comp", comp))
	compTyp, err := parseComponentType(comp)
	if err != nil {
		return err
	}

	mgr, ok := p.compManagers[compTyp]
	if !ok {
		return fmt.Errorf("not support this component type:%s", compTyp)
	}

	if err := mgr.Init(ctx, comp); err != nil {
		return err
	}

	mgr.mu.Lock()
	mgr.comps[comp.Name] = comp
	mgr.mu.Unlock()

	return nil
}

func (p *Processor) CloseAllComponents(ctx context.Context) error {
	var errs []error
	for _, mgr := range p.compManagers {
		func() {
			defer mgr.mu.Unlock()
			mgr.mu.Lock()

			for _, comp := range mgr.comps {
				if err := mgr.Close(comp); err != nil {
					errs = append(errs, err)
				}
			}
		}()
	}
	return errors.Join(errs...)
}
