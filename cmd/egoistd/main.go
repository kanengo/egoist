package main

import (
	"github.com/kanengo/egoist/pkg/runtime"
	"github.com/kanengo/egoist/pkg/signals"
	"github.com/kanengo/goutil/pkg/log"
	"go.uber.org/zap"

	_ "github.com/kanengo/egoist/cmd/egoistd/components"
)

func main() {
	config := &runtime.Config{
		Namespace: "test",
		AppId:     "test",
		AppPort:   50167,
		//GrpcPort:                0,
		ResourcesPath:           "",
		GracefulShutdownSeconds: 0,
		UnixDomainSocket:        "",
		APIListenAddress:        "",
	}

	ctx := signals.Context()
	rt, err := config.NewRuntime(ctx)
	if err != nil {
		log.Fatal("failed to new runtime", zap.Error(err), zap.Any("config", config))
	}

	if err := rt.Run(ctx); err != nil {
		return
	}
}
