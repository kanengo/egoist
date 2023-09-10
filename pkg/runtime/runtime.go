package runtime

import (
	"context"
	"github.com/kanengo/egoist/pkg/components"
	"github.com/kanengo/egoist/pkg/resources/components/v1alpha1"
	"github.com/kanengo/egoist/pkg/runtime/processor"
	"github.com/kanengo/goutil/pkg/concurrency"
	"github.com/kanengo/goutil/pkg/log"
	"go.uber.org/zap"
	"sync"
	"time"
)

type Runtime struct {
	internalConfig *internalConfig
	processor      *processor.Processor

	pendingComponents chan v1alpha1.Component

	runnerCloser *concurrency.RunnerCloserManager
	wg           sync.WaitGroup
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
			if err := rt.init(ctx); err != nil {
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
			rt.wg.Wait()
			return nil
		},
	); err != nil {
		return nil, err
	}

	return rt, nil
}

func (rt *Runtime) init(ctx context.Context) error {
	if err := rt.loadComponents(ctx); err != nil {
		return err
	}

	return nil
}

func (rt *Runtime) Run(ctx context.Context) {

}

func (rt *Runtime) loadComponents(ctx context.Context) error {
	loader := components.NewLocalComponents(rt.internalConfig.resourcesPath...)
	log.Info("Loading components")

	comps, err := loader.LoadComponents()
	if err != nil {
		return err
	}

	for _, comp := range comps {
		rt.addPendingComponent(ctx, comp)
	}

	return nil
}

func (rt *Runtime) processComponents(ctx context.Context) error {
	for comp := range rt.pendingComponents {
		if comp.Name == "" {
			continue
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
