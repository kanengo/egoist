package components

import "github.com/kanengo/egoist/pkg/resources/components/v1alpha1"

type Loader interface {
	LoadComponents() ([]v1alpha1.Component, error)
}
