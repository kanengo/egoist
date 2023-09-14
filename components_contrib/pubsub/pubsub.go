package pubsub

import (
	"context"
	"io"
)

type PubSub interface {
	io.Closer
	Init(ctx context.Context, metadata Metadata) error
	Publish(ctx context.Context, req *PublishRequest) error
	Subscribe(ctx context.Context) error
}
