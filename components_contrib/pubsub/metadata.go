package pubsub

import "github.com/kanengo/egoist/components_contrib/metadata"

type Metadata struct {
	metadata.Base `json:",inline"`
}

const (
	PubsubMetadataConsumerID = "pubsubConsumerID"
)
