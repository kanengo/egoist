package resources

import (
	"github.com/kanengo/egoist/pkg/resources/components/v1alpha1"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLocalResourcesLoader(t *testing.T) {
	request := NewLocalResourcesLoader[v1alpha1.Component]("test_component_path")

	yaml := `
apiVersion: dapr.io/v1alpha1
kind: Component
metadata:
   name: statestore
spec:
   type: state.couchbase
   metadata:
   - name: prop1
     value: value1
   - name: prop2
     value: value2
---
apiVersion: dapr.io/v1alpha1
kind: Subscription
metadata:
   name: pubsub
spec:
   type: pubsub.in-memory
   metadata:
   - name: rawPayload
     value: "true"
`
	components, errs := request.decodeYaml([]byte(yaml))
	assert.Empty(t, errs)
	assert.Len(t, components, 1)
	assert.Equal(t, "statestore", components[0].Name)
	assert.Equal(t, "state.couchbase", components[0].Spec.Type)
	assert.Len(t, components[0].Spec.Metadata, 2)
	assert.Equal(t, "prop1", components[0].Spec.Metadata[0].Name)
	assert.Equal(t, "value1", components[0].Spec.Metadata[0].Value.String())
}
