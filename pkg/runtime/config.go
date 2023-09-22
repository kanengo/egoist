package runtime

import (
	"context"
	"os"
	"os/user"
	"path"
	"strings"
	"time"

	"github.com/kanengo/egoist/utils"
	"github.com/kanengo/goutil/pkg/log"
	"go.uber.org/zap"
)

const (
	DefaultGracefulShutdownDuration = time.Second * 5

	DefaultMaxRequestBodySize = 4
	// DefaultAPIListenAddress is which address to listen for the Dapr HTTP and GRPC APIs. Empty string is all addresses.
	DefaultAPIListenAddress = ""
	// DefaultReadBufferSize is the default option for the maximum header size in KB for Dapr HTTP servers.
	DefaultReadBufferSize = 4

	DefaultRuntimeAppPort  = 50001
	DefaultRuntimeGrpcPort = 50017
)

type Config struct {
	Namespace               string
	AppId                   string
	AppPort                 int
	GrpcPort                int
	ResourcesPath           string
	GracefulShutdownSeconds int
	UnixDomainSocket        string
	APIListenAddress        string
}

type internalConfig struct {
	namespace                string
	id                       string
	appPort                  int
	grpcPort                 int
	gracefulShutdownDuration time.Duration
	hostAddress              string
	maxRequestBodySize       int
	readBufferSize           int
	unixDomainSocket         string
	apiListenAddresses       []string

	resourcesPath []string
}

func getNamespace() string {
	return os.Getenv("NAMESPACE")
}

func (c *Config) toInternalConfig() *internalConfig {
	intCfg := &internalConfig{
		namespace:        c.Namespace,
		id:               c.AppId,
		appPort:          c.AppPort,
		grpcPort:         c.GrpcPort,
		unixDomainSocket: c.UnixDomainSocket,
	}

	if intCfg.namespace == "" {
		intCfg.namespace = getNamespace()
	}

	intCfg.apiListenAddresses = strings.Split(c.APIListenAddress, ",")
	if len(intCfg.apiListenAddresses) == 0 {
		intCfg.apiListenAddresses = []string{""}
	}

	if c.GracefulShutdownSeconds <= 0 {
		intCfg.gracefulShutdownDuration = DefaultGracefulShutdownDuration
	} else {
		intCfg.gracefulShutdownDuration = time.Duration(c.GracefulShutdownSeconds) * time.Second
	}

	hostAddress, err := utils.GetHostAddress()
	if err != nil {
		log.Fatal("Failed to GetHostAddress", zap.Error(err))
	}

	if intCfg.maxRequestBodySize == 0 {
		intCfg.maxRequestBodySize = DefaultMaxRequestBodySize
	}

	if intCfg.readBufferSize == 0 {
		intCfg.readBufferSize = DefaultReadBufferSize
	}

	intCfg.hostAddress = hostAddress

	if intCfg.unixDomainSocket == "" {
		u, err := user.Current()
		if err != nil {
			log.Fatal("failed to init config when get current user", zap.Error(err))
		}
		intCfg.unixDomainSocket = path.Join(u.HomeDir, ".egoist", intCfg.unixDomainSocket)
	}

	if intCfg.appPort == 0 {
		intCfg.appPort = DefaultRuntimeAppPort
	}

	if intCfg.grpcPort == 0 {
		intCfg.grpcPort = DefaultRuntimeGrpcPort
	}

	if c.ResourcesPath == "" {
		u, err := user.Current()
		if err != nil {
			log.Fatal("failed to init config when get current user", zap.Error(err))
		}
		intCfg.resourcesPath = []string{path.Join(u.HomeDir, ".egoist/components")}
	} else {
		intCfg.resourcesPath = strings.Split(c.ResourcesPath, ",")
	}

	log.Debug("init internal config success", zap.String("appId", intCfg.id), zap.String("namespace", intCfg.namespace))

	return intCfg
}

func (c *Config) NewRuntime(ctx context.Context) (*Runtime, error) {
	intCfg := c.toInternalConfig()
	return newRuntime(intCfg)
}
