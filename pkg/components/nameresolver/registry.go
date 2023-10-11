package nameresolver

import (
	"fmt"
	"strings"

	"github.com/kanengo/egoist/components_contrib/nameresolver"
	"github.com/kanengo/egoist/pkg/components"
)

type Registry struct {
	messageBuses map[string]func() nameresolver.Resolver
}

// DefaultRegistry is the singleton with the registry.
var DefaultRegistry *Registry = NewRegistry()

// NewRegistry returns a new pub sub registry.
func NewRegistry() *Registry {
	return &Registry{
		messageBuses: map[string]func() nameresolver.Resolver{},
	}
}

// RegisterComponent adds a new message bus to the registry.
func (p *Registry) RegisterComponent(componentFactory func() nameresolver.Resolver, names ...string) {
	for _, name := range names {
		p.messageBuses[createFullName(name)] = componentFactory
	}
}

// Create instantiates a pub/sub based on `name`.
func (p *Registry) Create(name, version string) (nameresolver.Resolver, error) {
	if method, ok := p.getResolver(name, version); ok {
		return method(), nil
	}
	return nil, fmt.Errorf("couldn't find message bus %s/%s", name, version)
}

func (p *Registry) getResolver(name, version string) (func() nameresolver.Resolver, bool) {
	nameLower := strings.ToLower(name)
	versionLower := strings.ToLower(version)
	pubSubFn, ok := p.messageBuses[nameLower+"/"+versionLower]
	if ok {
		return p.wrapFn(pubSubFn), true
	}
	if components.IsInitialVersion(versionLower) {
		pubSubFn, ok = p.messageBuses[nameLower]
		if ok {
			return p.wrapFn(pubSubFn), true
		}
	}

	return nil, false
}

func (p *Registry) wrapFn(componentFactory func() nameresolver.Resolver) func() nameresolver.Resolver {
	return func() nameresolver.Resolver {
		return componentFactory()
	}
}

func createFullName(name string) string {
	return strings.ToLower("nameresolver." + name)
}
