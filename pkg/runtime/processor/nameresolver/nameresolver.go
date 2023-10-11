package nameresolver

import (
	"context"
	"sync/atomic"

	contribNameResolver "github.com/kanengo/egoist/components_contrib/nameresolver"
	"github.com/kanengo/egoist/pkg/components/nameresolver"
	"github.com/kanengo/egoist/pkg/resources/components/v1alpha1"
	"github.com/kanengo/egoist/pkg/runtime/meta"
)

type manager struct {
	component *componentItem
	registry  nameresolver.Registry
	meta      *meta.Meta
}

type componentItem struct {
	component contribNameResolver.Resolver
	resource  v1alpha1.Component

	closed  atomic.Bool
	closeCh chan struct{}
}

func (m *manager) Init(ctx context.Context, comp v1alpha1.Component) error {
	spec := comp.Spec
	ns, err := m.registry.Create(spec.Type, comp.ResourceVersion)
	if err != nil {
		return err
	}

	metaBase, err := m.meta.ToBaseMetadata(comp)
	if err != nil {
		return err
	}

	if err := ns.Init(ctx, contribNameResolver.Metadata{Base: metaBase}); err != nil {
		return err
	}

	old := m.component
	if old != nil {
		defer func() {
			_ = old.component.Close()
		}()
	}

	m.component = &componentItem{
		component: ns,
		resource:  comp,
		closeCh:   make(chan struct{}),
	}

	return nil
}

func (m *manager) Close(component v1alpha1.Component) error {
	if !m.component.closed.CompareAndSwap(false, true) {
		return nil
	}

	if m.component == nil {
		return nil
	}

	if component.Name != m.component.resource.Name {
		return nil
	}

	if err := m.component.component.Close(); err != nil {
		return err
	}

	close(m.component.closeCh)

	return nil
}
