package components

import (
	"github.com/kanengo/egoist/components_contrib/pubsub/pulsar"
	"github.com/kanengo/egoist/pkg/components/pubsub"
)

func init() {
	pubsub.DefaultRegistry.RegisterComponent(pulsar.NewPulsar, "pulsar")
}
