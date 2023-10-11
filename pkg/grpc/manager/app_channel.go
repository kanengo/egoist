package manager

import (
	"context"
	"fmt"
	"time"

	"github.com/kanengo/goutil/pkg/log"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AppChannelConfig struct {
	AppID              string
	Port               int
	UnixDomainSocket   string
	MaxReadBufferSize  int
	MaxWriteBufferSize int
	DialTimeout        time.Duration
}

func NewAppChannel(ctx context.Context, config AppChannelConfig) (*Channel, error) {
	channel := &Channel{
		appId: config.AppID,
	}

	if config.DialTimeout == 0 {
		config.DialTimeout = time.Second * 30
	}

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	}

	if config.MaxWriteBufferSize != 0 {
		opts = append(opts, grpc.WithWriteBufferSize(config.MaxWriteBufferSize<<20))
	}

	if config.MaxReadBufferSize != 0 {
		opts = append(opts, grpc.WithReadBufferSize(config.MaxReadBufferSize<<20))
	}

	var target string

	if config.UnixDomainSocket != "" {
		socket := fmt.Sprintf("%s/egoist-app-%s-grpc.socket", config.UnixDomainSocket, config.AppID)
		target = fmt.Sprintf("unix://%s", socket)
	} else {
		target = fmt.Sprintf("localhost:%d", config.Port)
	}

	log.Debug(fmt.Sprintf("app channel target:%s", target))
	channel.target = target

	go func() {
		for {
			ctx, cancel := context.WithTimeout(ctx, config.DialTimeout)
			c, err := grpc.DialContext(ctx, target, opts...)
			cancel()
			if err != nil {
				log.Debug("app channel dial failed, retry", zap.Error(err))
				continue
			}
			log.Debug("app channel init success")
			channel.cli = c
			break
		}
	}()

	return channel, nil
}
