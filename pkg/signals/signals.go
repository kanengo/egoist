package signals

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/kanengo/goutil/pkg/log"
	"go.uber.org/zap"
	"golang.org/x/net/context"
)

var (
	onlyOneSignalHandler = make(chan struct{})
	shutdownSignals      = []os.Signal{os.Interrupt, syscall.SIGTERM}
)

func Context() context.Context {
	close(onlyOneSignalHandler)

	ctx, cancel := context.WithCancel(context.Background())
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, shutdownSignals...)

	go func() {
		sig := <-sigCh
		log.Info("Received signal, beginning shutdown", zap.Any("sig", sig))
		cancel()

		sig = <-sigCh
		log.Fatal("Received signal during shutdown; existing immediately", zap.Any("sig", sig))
	}()

	return ctx
}
