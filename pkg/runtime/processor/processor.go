package processor

import (
	"context"
	"errors"
	"fmt"
	compPubsub "github.com/kanengo/egoist/pkg/components/pubsub"
	"github.com/kanengo/egoist/pkg/runtime/processor/pubsub"
	"strings"

	"github.com/kanengo/egoist/components_contrib"
	apiv1 "github.com/kanengo/egoist/pkg/api/v1"
	"github.com/kanengo/egoist/pkg/resources/components/v1alpha1"
)

type componentManger interface {
	Init(ctx context.Context, component v1alpha1.Component) error
	Close(component v1alpha1.Component) error
}

type PubsubManager interface {
	Publish(ctx context.Context, request apiv1.PublishEventRequest) error
	BulkPublish(ctx context.Context, request apiv1.BulkPublishRequest) (apiv1.BulkPublishResponse, error)
}

type Processor struct {
	pubsub       PubsubManager
	compManagers map[string]componentManger
}

func New(options Options) *Processor {
	p := &Processor{
		compManagers: make(map[string]componentManger),
	}

	p.compManagers[components_contrib.TypePubsub] = pubsub.New(pubsub.Options{
		ID:            options.ID,
		Namespace:     options.NameSpace,
		PodName:       options.PodName,
		ResourcesPath: nil,
		Registry:      compPubsub.DefaultRegistry,
		Meta:          nil,
	})

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

	return nil
}
