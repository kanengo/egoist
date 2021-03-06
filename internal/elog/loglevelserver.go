package zaplog

import (
	"net"
	"net/http"

	"go.uber.org/zap"
)

var setLevelPath string

func logLevelHttpServer(config *zaplogConfig, level zap.AtomicLevel) {
	logServerMux := http.NewServeMux()
	logServerMux.Handle(config.logApiPath, level)

	listener, err := net.Listen("tcp", config.listenAddr)
	if err != nil {
		Fatal("logLevelHttpServer Listen failed",
			zap.String("ipport", config.listenAddr),
			zap.Error(err))
	} else {
		setLevelPath = "http://" + listener.Addr().String() + config.logApiPath
	}

	go func() {
		err = http.Serve(listener, logServerMux)
		if err != nil {
			Fatal("logLevelHttpServer ListenAndServe failed",
				zap.String("ipport", listener.Addr().String()),
				zap.Error(err))
		}
	}()
}
