package pubsub

import (
	"context"

	contribPubsub "github.com/kanengo/egoist/components_contrib/pubsub"
	apiv1 "github.com/kanengo/egoist/pkg/api/v1"
	"github.com/kanengo/goutil/pkg/log"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

type bulkPubsub struct {
	ps contribPubsub.PubSub
}

func (b *bulkPubsub) BulkPublish(ctx context.Context, req *contribPubsub.BulkPublishRequest) (*contribPubsub.BulkPublishResponse, error) {
	var failedEntry []contribPubsub.FailedEntry
	for _, entry := range req.Entries {
		publishRequest := &contribPubsub.PublishRequest{
			Data:       entry.Event,
			PubsubName: req.PubsubName,
			Topic:      req.Topic,
			Metadata:   req.Metadata,
		}

		if err := b.ps.Publish(ctx, publishRequest); err != nil {
			failedEntry = append(failedEntry, contribPubsub.FailedEntry{
				EntryId: entry.EntryId,
				Err:     err,
			})
		}
	}

	return &contribPubsub.BulkPublishResponse{FailedEntry: failedEntry}, nil
}

func (b *bulkPubsub) BulkSubscribe(ctx context.Context, req *contribPubsub.SubscribeRequest, handler contribPubsub.BulkHandler) error {
	return nil
}

type subscribeStreamHandler struct {
	streamServer apiv1.API_SubscribeStreamServer
	ctx          context.Context
	//cancel       context.CancelFunc
}

func (h *subscribeStreamHandler) handler() contribPubsub.Handler {
	return func(ctx context.Context, msg *contribPubsub.NewMessage) error {
		cloudEvent := &apiv1.CloudEvent{}
		err := proto.Unmarshal(msg.Data, cloudEvent)
		if err != nil {
			log.Error("SubscribeStream handler proto.Unmarshal failed", zap.Error(err), zap.String("topic", msg.Topic))
			return err
		}
		entry := &apiv1.SubscribeEntry{
			Topic:  msg.Topic,
			Events: cloudEvent,
		}
		err = h.streamServer.Send(&apiv1.SubscribeResponse{Entries: []*apiv1.SubscribeEntry{entry}})
		if err != nil {
			//if s, ok := status.FromError(err); ok {
			//	if s.Code() == codes.Unavailable {
			//		h.cancel()
			//		return err
			//	}
			//}
			log.Error("SubscribeStream handler streamServer.Send failed", zap.Error(err), zap.String("topic", msg.Topic))
			return err
		}
		return nil
	}
}

func (h *subscribeStreamHandler) bulkHandler() contribPubsub.BulkHandler {
	return func(ctx context.Context, msgs []*contribPubsub.NewMessage) error {
		entries := make([]*apiv1.SubscribeEntry, 0, len(msgs))
		for _, msg := range msgs {
			cloudEvent := &apiv1.CloudEvent{}
			err := proto.Unmarshal(msg.Data, cloudEvent)
			if err != nil {
				log.Error("SubscribeStream handler proto.Unmarshal failed", zap.Error(err), zap.String("topic", msg.Topic))
				return err
			}
			entry := &apiv1.SubscribeEntry{
				Topic:  msg.Topic,
				Events: cloudEvent,
			}
			entries = append(entries, entry)
		}
		err := h.streamServer.Send(&apiv1.SubscribeResponse{Entries: entries})
		if err != nil {
			//if s, ok := status.FromError(err); ok {
			//	if s.Code() == codes.Unavailable {
			//		h.cancel()
			//		return err
			//	}
			//}
			log.Error("SubscribeStream handler streamServer.Send failed", zap.Error(err))
			return err
		}
		return nil
	}
}
