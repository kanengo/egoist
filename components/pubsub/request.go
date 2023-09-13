package pubsub

type PublishRequest struct {
	Data       []byte            `json:"data"`
	PubsubName string            `json:"pubsubName"`
	Topic      string            `json:"topic"`
	Metadata   map[string]string `json:"metadata"`
}
