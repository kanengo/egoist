package nameresolver

import (
	"context"
)

type Resolver interface {
	Init(ctx context.Context, metadata Metadata) error
	ResolveID(ctx context.Context, id string) (string, error)
}
