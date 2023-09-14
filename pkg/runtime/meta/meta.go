package meta

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/kanengo/egoist/components_contrib/metadata"
	"github.com/kanengo/egoist/pkg/resources/common"
	"github.com/kanengo/egoist/pkg/resources/components/v1alpha1"
	"strings"
)

type Meta struct {
	id        string
	podName   string
	namespace string
}

func (m *Meta) ToBaseMetadata(comp v1alpha1.Component) (metadata.Base, error) {
	// Add global wasm strict sandbox config to the wasm component metadata
	//if components.IsWasmComponentType(comp.Spec.Type) {
	//	m.AddWasmStrictSandbox(&comp)
	//}

	props, err := m.convertItemsToProps(comp.Spec.Metadata)
	if err != nil {
		return metadata.Base{}, err
	}

	return metadata.Base{
		Properties: props,
		Name:       comp.Name,
	}, nil
}

func (m *Meta) convertItemsToProps(items []common.NameValuePair) (map[string]string, error) {
	properties := map[string]string{}
	for _, c := range items {
		val := c.Value.String()
		for strings.Contains(val, "{uuid}") {
			val = strings.Replace(val, "{uuid}", uuid.New().String(), 1)
		}
		if strings.Contains(val, "{podName}") {
			if m.podName == "" {
				return nil, fmt.Errorf("failed to parse metadata: property %s refers to {podName} but podName is not set", c.Name)
			}
			val = strings.ReplaceAll(val, "{podName}", m.podName)
		}
		val = strings.ReplaceAll(val, "{namespace}", fmt.Sprintf("%s.%s", m.namespace, m.id))
		val = strings.ReplaceAll(val, "{appID}", m.id)
		properties[c.Name] = val
	}
	return properties, nil
}
