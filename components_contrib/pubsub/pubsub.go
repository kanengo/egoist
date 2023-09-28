package pubsub

import (
	"context"
	"io"
)

type PubSub interface {
	io.Closer
	Init(ctx context.Context, metadata Metadata) error
	Publish(ctx context.Context, req *PublishRequest) error
	Subscribe(ctx context.Context, req *SubscribeRequest, handler Handler) error
}

type BulkPublisher interface {
	BulkPublish(ctx context.Context, req *BulkPublishRequest) (*BulkPublishResponse, error)
}

type BulkSubscriber interface {
	BulkSubscribe(ctx context.Context, req *SubscribeRequest, handler BulkHandler) error
}

type Handler func(ctx context.Context, msg *NewMessage) error

type BulkHandler func(ctx context.Context, msgs []*NewMessage) error
