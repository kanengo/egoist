package main

import (
	"context"
	"errors"
	"io"
	"net"
	"os"
	"time"

	apiv1 "github.com/kanengo/egoist/pkg/api/v1"
	"github.com/kanengo/goutil/pkg/log"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	defer func() {
		_ = os.Remove("/Users/kanonlee/.egoist/egoist-app-test-grpc.socket")
	}()
	l, err := net.Listen("unix", "/Users/kanonlee/.egoist/egoist-app-test-grpc.socket")
	if err != nil {
		log.Fatal("Failed to listen gRPC server on Unix",
			zap.Error(err))
	}
	svr := grpc.NewServer()
	defer func() {
		svr.GracefulStop()
	}()
	go func() {
		if err := svr.Serve(l); err != nil {
			log.Fatal("serve failed", zap.Error(err))
		}
	}()

	cli, err := grpc.DialContext(ctx, "unix:///Users/kanonlee/.egoist/egoist-test-grpc.socket", grpc.WithTransportCredentials(
		insecure.NewCredentials()))
	if err != nil {
		log.Fatal("", zap.Error(err))
	}
	defer cli.Close()

	apiCli := apiv1.NewAPIClient(cli)

	subStreamClient, err := apiCli.SubscribeStream(context.Background(), &apiv1.SubscribeRequest{
		Configs: []*apiv1.SubscribeConfig{
			{
				PubsubName:          "test",
				Topic:               "topic/test",
				EnableBulk:          true,
				MaxBulkEventCount:   100,
				MaxBulkEventAwaitMs: 100,
				Metadata: map[string]string{
					"pubsubConsumerID": "test",
				},
			},
		},
	})
	if err != nil {
		log.Error("failed to SubscribeStream", zap.Error(err))
		return
	}

	defer func() {
		_ = subStreamClient.CloseSend()
	}()

	go func() {
		for {
			msg, err := subStreamClient.Recv()
			if err != nil {
				if errors.Is(err, io.EOF) {
					log.Info("SubscribeStream Recv eof")
				} else {
					log.Error("SubscribeStream Recv failed", zap.Error(err))
				}
				break
			}
			log.Debug("SubscribeStream Recv", zap.Any("msg", msg))
		}
	}()

	_, _ = apiCli.PublishEvent(context.Background(), &apiv1.PublishEventRequest{
		Data:        []byte("test pub sub"),
		PubsubName:  "test",
		Topic:       "topic/test",
		Metadata:    nil,
		ContentType: "",
	})

	time.Sleep(time.Second * 10)
}
