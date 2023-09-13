package pubsub

import (
	"context"
	"io"
)

type PubSub interface {
	io.Closer
	Init(ctx context.Context, metadata Metadata)
	Publish(ctx context.Context, req *PublishRequest)
	Subscribe(ctx context.Context)
}
