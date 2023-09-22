package runtime

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/kanengo/egoist/pkg/grpc"
	"github.com/kanengo/egoist/pkg/messaging"

	"github.com/kanengo/egoist/pkg/components"
	"github.com/kanengo/egoist/pkg/resources/components/v1alpha1"
	"github.com/kanengo/egoist/pkg/runtime/processor"
	"github.com/kanengo/goutil/pkg/concurrency"
	"github.com/kanengo/goutil/pkg/log"
	"go.uber.org/zap"
)

type Runtime struct {
	internalConfig *internalConfig
	processor      *processor.Processor

	pendingComponents chan v1alpha1.Component
	proxy             messaging.Proxy

	runnerCloser *concurrency.RunnerCloserManager
	wg           sync.WaitGroup

	api grpc.API
}

func newRuntime(intCfg *internalConfig) (*Runtime, error) {
	rt := &Runtime{
		internalConfig:    intCfg,
		pendingComponents: make(chan v1alpha1.Component),
		processor: processor.New(processor.Options{
			ID:        intCfg.id,
			NameSpace: intCfg.namespace,
		}),
	}

	rt.runnerCloser = concurrency.NewRunnerCloserManager(rt.internalConfig.gracefulShutdownDuration,
		rt.processComponents,
		func(ctx context.Context) error {
			defer func() {
				close(rt.pendingComponents)
			}()

			start := time.Now()
			log.Info("egoist runtime init", zap.String("appId", rt.internalConfig.id))
			if err := rt.initRuntime(ctx); err != nil {
				return err
			}

			d := time.Since(start)
			log.Info("egoist initialized. Status: Running", zap.Duration("initElapsed", d))

			<-ctx.Done()

			return nil
		},
	)

	if err := rt.runnerCloser.AddCloser(
		func() error {
			log.Info("egoist is shutting down")
			var errs []error
			rt.wg.Wait()
			errs = append(errs, rt.cleanSocket())
			return errors.Join(errs...)
		},
	); err != nil {
		return nil, err
	}

	return rt, nil
}

func (rt *Runtime) initRuntime(ctx context.Context) error {
	if err := rt.loadComponents(ctx); err != nil {
		//return err
	}

	api := grpc.NewAPI(grpc.APIOptions{})
	if err := rt.startGRPCAPIServer(api); err != nil {
		return fmt.Errorf("faild to start API gRPC server: %w", err)
	}

	rt.api = api

	if rt.internalConfig.unixDomainSocket != "" {
		log.Info("API gRPC server is running on Unix Domain Socket")
	} else {
		log.Info("API gRPC server is running", zap.Any("port", rt.internalConfig.grpcPort))
	}

	return nil
}

func (rt *Runtime) Run(ctx context.Context) error {
	if err := rt.runnerCloser.Run(ctx); err != nil {
		log.Error("failed to run runtime", zap.Error(err))
		return err
	}

	return nil
}

func (rt *Runtime) getDefaultGPRCServerConfig() grpc.ServerConfig {
	return grpc.ServerConfig{
		AppID:                rt.internalConfig.id,
		HostAddress:          rt.internalConfig.hostAddress,
		Port:                 rt.internalConfig.appPort,
		APIListenAddresses:   rt.internalConfig.apiListenAddresses,
		NameSpace:            rt.internalConfig.namespace,
		MaxRequestBodySizeMB: rt.internalConfig.maxRequestBodySize,
		UnixDomainSocket:     rt.internalConfig.unixDomainSocket,
		ReadBufferSizeKB:     rt.internalConfig.readBufferSize,
	}
}

func (rt *Runtime) startGRPCAPIServer(api grpc.API) error {

	serverConf := rt.getDefaultGPRCServerConfig()
	serverConf.Port = rt.internalConfig.grpcPort

	server := grpc.NewAPIServer(api, serverConf, rt.proxy)
	if err := server.StartNonBlocking(); err != nil {
		return err
	}
	if err := rt.runnerCloser.AddCloser(server); err != nil {
		log.Error("start gRPC API server failed to add closer", zap.Error(err))
		return err
	}

	return nil
}

func (rt *Runtime) loadComponents(ctx context.Context) error {
	loader := components.NewLocalComponents(rt.internalConfig.resourcesPath...)
	log.Info("Loading components")

	comps, err := loader.LoadComponents()
	if err != nil {
		log.Error("failed to loadComponents", zap.Error(err), zap.Strings("resourcePath", rt.internalConfig.resourcesPath))
		return err
	}

	for _, comp := range comps {
		rt.addPendingComponent(ctx, comp)
	}

	return nil
}

func (rt *Runtime) processComponents(ctx context.Context) error {
	log.Info("start to process components")
	for comp := range rt.pendingComponents {
		if comp.Name == "" {
			continue
		}
		if err := rt.processor.InitComponent(ctx, comp); err != nil {
			log.Error("failed to update component", zap.Any("comp", comp))
		}
	}

	return nil
}

func (rt *Runtime) addPendingComponent(ctx context.Context, comp v1alpha1.Component) bool {
	select {
	case <-ctx.Done():
		return false
	case rt.pendingComponents <- comp:
		return true
	}
}

func (rt *Runtime) cleanSocket() error {
	if rt.internalConfig.unixDomainSocket != "" {
		_ = os.Remove(fmt.Sprintf("%s/egoist-%s-grpc.socket", rt.internalConfig.unixDomainSocket,
			rt.internalConfig.id))
		//return err
	}
	return nil
}
