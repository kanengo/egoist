package components

import (
	"github.com/kanengo/egoist/pkg/resources"
	"github.com/kanengo/egoist/pkg/resources/components/v1alpha1"
)

type LocalComponents struct {
	loader resources.Loader[v1alpha1.Component]
}

func (l *LocalComponents) LoadComponents() ([]v1alpha1.Component, error) {
	return l.loader.Load()
}

func NewLocalComponents(resourcesPaths ...string) *LocalComponents {
	return &LocalComponents{loader: resources.NewLocalResourcesLoader[v1alpha1.Component](resourcesPaths...)}
}
