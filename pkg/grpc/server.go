package grpc

import (
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/kanengo/egoist/pkg/messaging"
	"github.com/kanengo/goutil/pkg/log"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

const (
	GrpcServerKindApi = "apiserver"
)

type Server interface {
	io.Closer
	StartNonBlocking() error
}

type server struct {
	api              API
	config           ServerConfig
	servers          []*grpc.Server
	closed           atomic.Bool
	closeCh          chan struct{}
	kind             string
	maxConnectionAge *time.Duration
	wg               sync.WaitGroup
	proxy            messaging.Proxy
}

func (s *server) Close() error {
	defer s.wg.Wait()
	if s.closed.CompareAndSwap(false, true) {
		close(s.closeCh)
	}
	s.wg.Add(len(s.servers))
	for _, server := range s.servers {
		go func(server *grpc.Server) {
			defer s.wg.Done()
			server.GracefulStop()
		}(server)
	}

	if s.api != nil {
		if err := s.api.Close(); err != nil {
			return err
		}
	}

	return nil
}

func (s *server) getMiddlewareOptions() []grpc.ServerOption {
	return []grpc.ServerOption{}
}

func (s *server) getGRPCServer() (*grpc.Server, error) {
	opts := s.getMiddlewareOptions()
	if s.maxConnectionAge != nil {
		opts = append(opts, grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionAge: *s.maxConnectionAge,
		}))
	}

	if s.config.MaxRequestBodySizeMB != 0 {
		opts = append(opts, grpc.MaxRecvMsgSize(s.config.MaxRequestBodySizeMB<<20))
	}

	if s.config.ReadBufferSizeKB != 0 {
		opts = append(opts, grpc.MaxHeaderListSize(uint32(s.config.ReadBufferSizeKB<<10)))
	}

	if s.proxy != nil {

	}

	return grpc.NewServer(opts...), nil
}

func (s *server) StartNonBlocking() error {
	var listeners []net.Listener
	if s.config.UnixDomainSocket != "" && s.kind == GrpcServerKindApi {
		socket := fmt.Sprintf("%s/egoist-%s-grpc.socket", s.config.UnixDomainSocket, s.config.AppID)
		l, err := net.Listen("unix", socket)
		if err != nil {
			return err
		}
		log.Info("gRpc server listening on UNIX socket", zap.String("socket", socket))
		listeners = append(listeners, l)
	} else {
		for _, apiListenAddress := range s.config.APIListenAddresses {
			addr := apiListenAddress + ":" + strconv.Itoa(s.config.Port)
			l, err := net.Listen("tcp", addr)
			if err != nil {
				log.Error("Failed to listen gRPC server on TCP", zap.String("address", addr),
					zap.Error(err))
			} else {
				log.Info("gRPC server listening on TCP", zap.String("address", addr))
				listeners = append(listeners, l)
			}
		}
	}

	if len(listeners) == 0 {
		return errors.New("could not listen on any endpoint")
	}

	for _, listener := range listeners {
		server, err := s.getGRPCServer()
		if err != nil {
			return err
		}
		s.servers = append(s.servers, server)

		//RegisterService
		s.wg.Add(1)
		go func(server *grpc.Server, l net.Listener) {
			defer s.wg.Done()
			if err := server.Serve(l); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
				log.Fatal("gRPC server error", zap.Error(err))
			}
		}(server, listener)
	}

	return nil
}

func NewAPIServer(api API, config ServerConfig, proxy messaging.Proxy) Server {
	return &server{
		api:     api,
		config:  config,
		servers: nil,
		closed:  atomic.Bool{},
		closeCh: nil,
		proxy:   proxy,
	}
}

func NewInternalServer() Server {
	return &server{
		config:  ServerConfig{},
		servers: nil,
		closed:  atomic.Bool{},
		closeCh: nil,
	}
}
