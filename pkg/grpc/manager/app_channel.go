package manager

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AppChannelConfig struct {
	AppID              string
	Port               int
	UnixDomainSocket   string
	MaxReadBufferSize  int
	MaxWriteBufferSize int
}

func NewAPPChannel(ctx context.Context, config AppChannelConfig) (*Channel, error) {
	channel := &Channel{
		appId: config.AppID,
	}

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	if config.MaxWriteBufferSize != 0 {
		opts = append(opts, grpc.WithWriteBufferSize(config.MaxWriteBufferSize<<20))
	}

	if config.MaxReadBufferSize != 0 {
		opts = append(opts, grpc.WithReadBufferSize(config.MaxReadBufferSize<<20))
	}

	if config.UnixDomainSocket != "" {
		socket := fmt.Sprintf("%s/egoist-app-%s-grpc.socket", config.UnixDomainSocket, config.AppID)
		c, err := grpc.DialContext(ctx, fmt.Sprintf("unix://%s", socket), opts...)
		if err != nil {
			return nil, err
		}
		channel.cli = c
	} else {
		c, err := grpc.DialContext(ctx, fmt.Sprintf("localhost:%d", config.Port), opts...)
		if err != nil {
			return nil, err
		}
		channel.cli = c
	}

	return channel, nil
}
