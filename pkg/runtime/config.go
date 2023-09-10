package runtime

import (
	"context"
	"time"
)

const (
	DefaultGracefulShutdownDuration = time.Second * 5
)

type Config struct {
	Namespace               string
	AppId                   string
	AppPort                 int
	GrpcPort                int
	ResourcesPath           []string
	GracefulShutdownSeconds int
}

type internalConfig struct {
	namespace                string
	id                       string
	appPort                  int
	grpcPort                 int
	gracefulShutdownDuration time.Duration

	resourcesPath []string
}

func (c *Config) toInternalConfig() *internalConfig {
	intCfg := &internalConfig{
		namespace:     c.Namespace,
		id:            c.AppId,
		appPort:       c.AppPort,
		grpcPort:      c.GrpcPort,
		resourcesPath: c.ResourcesPath,
	}

	if c.GracefulShutdownSeconds <= 0 {
		intCfg.gracefulShutdownDuration = DefaultGracefulShutdownDuration
	} else {
		intCfg.gracefulShutdownDuration = time.Duration(c.GracefulShutdownSeconds) * time.Second
	}

	return intCfg
}

func (c *Config) NewRuntime(ctx context.Context) (*Runtime, error) {
	intCfg := c.toInternalConfig()
	return newRuntime(intCfg)
}
