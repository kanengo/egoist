package pubsub

import (
	"golang.org/x/exp/slices"
)

const (
	// FeatureMessageTTL is the feature to handle message TTL.
	FeatureMessageTTL Feature = "MESSAGE_TTL"
	// FeatureSubscribeWildcards is the feature to allow subscribing to topics/queues using a wildcard.
	FeatureSubscribeWildcards Feature = "SUBSCRIBE_WILDCARDS"
	FeatureBulkPublish        Feature = "BULK_PUBSUB"
)

// Feature names a feature that can be implemented by PubSub components.
type Feature string

// IsPresent checks if a given feature is present in the list.
func (f Feature) IsPresent(features []Feature) bool {
	return slices.Contains(features, f)
}
