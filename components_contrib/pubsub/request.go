package pubsub

type PublishRequest struct {
	Data       []byte            `json:"data"`
	PubsubName string            `json:"pubsubName"`
	Topic      string            `json:"topic"`
	Metadata   map[string]string `json:"metadata"`
}

// BulkPublishRequest is the request to publish mutilple messages.
type BulkPublishRequest struct {
	Entries    []BulkMessageEntry `json:"entries"`
	PubsubName string             `json:"pubsubname"`
	Topic      string             `json:"topic"`
	Metadata   map[string]string  `json:"metadata"`
}

// SubscribeRequest is the request to subscribe to a topic.
type SubscribeRequest struct {
	Topic               string              `json:"topic"`
	Metadata            map[string]string   `json:"metadata"`
	BulkSubscribeConfig BulkSubscribeConfig `json:"bulkSubscribe,omitempty"`
}

// BulkMessageEntry represents a single message inside a bulk request.
type BulkMessageEntry struct {
	EntryId     string            `json:"entryId"` //nolint:stylecheck
	Event       []byte            `json:"event"`
	ContentType string            `json:"contentType,omitempty"`
	Metadata    map[string]string `json:"metadata"`
}

type BulkSubscribeConfig struct {
	MaxMessagesCount   int `json:"maxMessagesCount,omitempty"`
	MaxAwaitDurationMs int `json:"maxAwaitDurationMs,omitempty"`
}

// NewMessage is an event arriving from a message bus instance.
type NewMessage struct {
	Data        []byte            `json:"data"`
	Topic       string            `json:"topic"`
	Metadata    map[string]string `json:"metadata"`
	ContentType *string           `json:"contentType,omitempty"`
}
