package pubsub

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	_struct "github.com/golang/protobuf/ptypes/struct"
	contribPubsub "github.com/kanengo/egoist/components_contrib/pubsub"
	apiv1 "github.com/kanengo/egoist/pkg/api/v1"
	"github.com/kanengo/egoist/pkg/components"
	"github.com/kanengo/egoist/pkg/runtime/meta"
	"github.com/kanengo/goutil/pkg/log"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/kanengo/egoist/pkg/components/pubsub"
	"github.com/kanengo/egoist/pkg/resources/components/v1alpha1"
)

type Options struct {
	ID            string
	Namespace     string
	IsHTTP        bool
	PodName       string
	ResourcesPath []string

	Registry *pubsub.Registry
	Meta     *meta.Meta
	//TracingSpec    *config.TracingSpec
	//GRPC           *manager.Manager
	//Channels       *channels.Channels
	//OperatorClient operatorv1.OperatorClient
}

type pubsubCompItem struct {
	Component contribPubsub.PubSub
	BulkPub   contribPubsub.BulkPublisher
	BulkSub   contribPubsub.BulkSubscriber
	Resource  v1alpha1.Component

	closed  atomic.Bool
	closeCh chan struct{}
}

type pubSubManager struct {
	//opts      Options
	compStore *components.CompStore[*pubsubCompItem]
	meta      *meta.Meta
	registry  *pubsub.Registry

	lock sync.RWMutex
}

func NewManager(opts Options) *pubSubManager {
	pm := &pubSubManager{
		//opts:      opts,
		compStore: components.NewCompStore[*pubsubCompItem](),
		meta:      opts.Meta,
		registry:  pubsub.DefaultRegistry,
	}

	return pm
}

func (p *pubSubManager) Publish(ctx context.Context, request *apiv1.PublishEventRequest) (*apiv1.PublishEventResponse, error) {
	ps, ok := p.compStore.Get(request.PubsubName)
	if !ok {
		return nil, fmt.Errorf("no %s pubusb component fonund", request.PubsubName)
	}

	cloudEventType := ""
	cloudEventSource := ""
	cloudEventId := ""

	if cet, ok := request.Metadata[components.CloudEventMetadataType]; ok {
		cloudEventType = cet
	}

	if ces, ok := request.Metadata[components.CloudEventMetadataSource]; ok {
		cloudEventSource = ces
	}

	if ceid, ok := request.Metadata[components.CloudEventMetadataId]; ok {
		cloudEventId = ceid
	} else {
		cloudEventId = gonanoid.Must()
	}

	cloudEvent := &apiv1.CloudEvent{
		Id:              cloudEventId,
		Data:            request.Data,
		Source:          cloudEventSource,
		SpecVersion:     "1.0",
		Type:            cloudEventType,
		DataContentType: "application/protobuf",
		Timestamp:       time.Now().Unix(),
		Extensions:      make(map[string]*_struct.Value, len(request.Metadata)),
	}

	for k, v := range request.Metadata {
		if _, ok := components.CloudEventMetadataKeys[k]; ok {
			continue
		}
		cloudEvent.Extensions[k] = structpb.NewStringValue(v)
	}

	data, _ := proto.Marshal(cloudEvent)
	publishRequest := &contribPubsub.PublishRequest{
		Data:       data,
		PubsubName: request.PubsubName,
		Topic:      request.Topic,
		Metadata:   request.Metadata,
	}
	if err := ps.Component.Publish(ctx, publishRequest); err != nil {
		return nil, err
	}

	return &apiv1.PublishEventResponse{}, nil
}

func (p *pubSubManager) BulkPublish(ctx context.Context, request *apiv1.BulkPublishRequest) (*apiv1.BulkPublishResponse, error) {
	psItem, ok := p.compStore.Get(request.PubsubName)
	if !ok {
		return nil, fmt.Errorf("no %s pubusb component fonund", request.PubsubName)
	}

	bulkRequest := &contribPubsub.BulkPublishRequest{
		Entries:    make([]contribPubsub.BulkMessageEntry, 0, len(request.Entries)),
		PubsubName: request.PubsubName,
		Topic:      request.Topic,
		Metadata:   request.Metadata,
	}

	for _, entry := range request.Entries {
		cloudEventType := ""
		cloudEventSource := ""
		cloudEventId := ""

		if cet, ok := entry.Metadata[components.CloudEventMetadataType]; ok {
			cloudEventType = cet
		}

		if ces, ok := entry.Metadata[components.CloudEventMetadataSource]; ok {
			cloudEventSource = ces
		}

		if ceid, ok := request.Metadata[components.CloudEventMetadataId]; ok {
			cloudEventId = ceid
		} else {
			cloudEventId = gonanoid.Must()
		}

		cloudEvent := &apiv1.CloudEvent{
			Id:              cloudEventId,
			Data:            entry.Event,
			Source:          cloudEventSource,
			SpecVersion:     "1.0",
			Type:            cloudEventType,
			DataContentType: "application/protobuf",
			Timestamp:       time.Now().Unix(),
			Extensions:      make(map[string]*_struct.Value, len(request.Metadata)),
		}

		for k, v := range request.Metadata {
			if _, ok := components.CloudEventMetadataKeys[k]; ok {
				continue
			}
			cloudEvent.Extensions[k] = structpb.NewStringValue(v)
		}

		data, _ := proto.Marshal(cloudEvent)

		message := contribPubsub.BulkMessageEntry{
			EntryId:     cloudEventId,
			Event:       data,
			ContentType: "application/protobuf",
			Metadata:    nil,
		}

		bulkRequest.Entries = append(bulkRequest.Entries, message)
	}

	bulkResponse, err := psItem.BulkPub.BulkPublish(ctx, bulkRequest)
	if err != nil {
		return nil, err
	}

	response := &apiv1.BulkPublishResponse{FailedEntries: make([]*apiv1.BulkPublishResponseFailedEntry, 0, len(bulkResponse.FailedEntry))}
	for _, failedEntry := range bulkResponse.FailedEntry {
		response.FailedEntries = append(response.FailedEntries, &apiv1.BulkPublishResponseFailedEntry{
			EntryId: failedEntry.EntryId,
			Error:   err.Error(),
		})
	}

	return response, nil
}

func (p *pubSubManager) SubscribeStream(request *apiv1.SubscribeRequest, streamServer apiv1.API_SubscribeStreamServer) error {
	log.Debug("SubscribeStream start request", zap.Any("request", request))
	wg := sync.WaitGroup{}
	for _, config := range request.Configs {
		if config.Topic == "" {
			return status.Error(codes.InvalidArgument, "topic must not empty")
		}
		psItem, ok := p.compStore.Get(config.PubsubName)
		if !ok {
			return status.Error(codes.InvalidArgument, fmt.Sprintf("no %s pubusb component fonund", config.PubsubName))
		}

		if config.EnableBulk {
			subscribeRequest := &contribPubsub.SubscribeRequest{
				Topic:    config.Topic,
				Metadata: config.Metadata,
				BulkSubscribeConfig: contribPubsub.BulkSubscribeConfig{
					MaxMessagesCount:   int(config.MaxBulkEventCount),
					MaxAwaitDurationMs: int(config.MaxBulkEventAwaitMs),
				},
			}
			if subscribeRequest.BulkSubscribeConfig.MaxMessagesCount == 0 {
				subscribeRequest.BulkSubscribeConfig.MaxMessagesCount = 1000
			}
			if subscribeRequest.BulkSubscribeConfig.MaxAwaitDurationMs == 0 {
				subscribeRequest.BulkSubscribeConfig.MaxAwaitDurationMs = 100
			}

			wg.Add(1)
			if psItem.BulkSub != nil {
				go func() {
					defer wg.Done()
					ctx := streamServer.Context()
					h := &subscribeStreamHandler{streamServer: streamServer, ctx: ctx}
					err := psItem.BulkSub.BulkSubscribe(ctx, subscribeRequest, h.bulkHandler())
					if err != nil {
						log.Error("SubscribeStream BulkSubscribe failed", zap.Error(err), zap.String("topic", config.Topic))
						return
					}
					select {
					case <-psItem.closeCh:
						log.Info("SubscribeStream component closed", zap.String("topic", config.Topic), zap.String("pubsubName", config.PubsubName))
					case <-ctx.Done():
						log.Info("SubscribeStream context done", zap.String("topic", config.Topic), zap.String("pubsubName", config.PubsubName))
					}
				}()
			} else {
				go func() {
					ticker := time.NewTicker(time.Duration(subscribeRequest.BulkSubscribeConfig.MaxAwaitDurationMs) * time.Millisecond)
					bulkEntries := make([]*apiv1.SubscribeEntry, 0, subscribeRequest.BulkSubscribeConfig.MaxMessagesCount)
					defer func() {
						ticker.Stop()
						wg.Done()
					}()
					ctx := streamServer.Context()
					bulkCh := make(chan *apiv1.SubscribeEntry, subscribeRequest.BulkSubscribeConfig.MaxMessagesCount)
					err := psItem.Component.Subscribe(ctx, subscribeRequest, func(ctx context.Context, msg *contribPubsub.NewMessage) error {
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
						//log.Debug("Subscribe cloudevent", zap.Any("cloudEvent", cloudEvent))
						bulkCh <- entry
						return nil
					})
					if err != nil {
						log.Error("SubscribeStream Subscribe failed", zap.Error(err), zap.String("topic", config.Topic))
						return
					}

					for {
						select {
						case <-ticker.C:
							if len(bulkEntries) > 0 {
								log.Debug("SubscribeStream bulk max await duration", zap.Int("count", len(bulkEntries)))
								err := streamServer.Send(&apiv1.SubscribeResponse{Entries: bulkEntries})
								if err != nil {
									log.Error("SubscribeStream handler streamServer.Send failed", zap.Error(err))
									return
								}
								bulkEntries = bulkEntries[:0]
							}
						case entry := <-bulkCh:
							bulkEntries = append(bulkEntries, entry)
							if len(bulkEntries) >= subscribeRequest.BulkSubscribeConfig.MaxMessagesCount {
								log.Debug("SubscribeStream bulk max message count", zap.Int("count", len(bulkEntries)))
								err := streamServer.Send(&apiv1.SubscribeResponse{Entries: bulkEntries})
								if err != nil {
									log.Error("SubscribeStream handler streamServer.Send failed", zap.Error(err))
									return
								}
								bulkEntries = bulkEntries[:0]
							}
						case <-psItem.closeCh:
							log.Info("SubscribeStream component closed", zap.String("topic", config.Topic), zap.String("pubsubName", config.PubsubName))
							return
						case <-ctx.Done():
							log.Info("SubscribeStream stream done", zap.String("topic", config.Topic), zap.String("pubsubName", config.PubsubName))
							return
						}
					}
				}()
			}
		} else {
			subscribeRequest := &contribPubsub.SubscribeRequest{
				Topic:               config.Topic,
				Metadata:            config.Metadata,
				BulkSubscribeConfig: contribPubsub.BulkSubscribeConfig{},
			}
			go func() {
				defer func() {
					wg.Done()
				}()
				ctx := streamServer.Context()
				h := &subscribeStreamHandler{streamServer: streamServer, ctx: ctx}
				if err := psItem.Component.Subscribe(ctx, subscribeRequest, h.handler()); err != nil {
					log.Error("SubscribeStream Subscribe failed", zap.Error(err), zap.String("topic", config.Topic))
					return
				}
				select {
				case <-psItem.closeCh:
					log.Info("SubscribeStream component closed", zap.String("topic", config.Topic), zap.String("pubsubName", config.PubsubName))
				case <-ctx.Done():
					log.Info("SubscribeStream context done", zap.String("topic", config.Topic), zap.String("pubsubName", config.PubsubName))
				}
			}()
		}
	}
	log.Debug("SubscribeStream listening")
	wg.Wait()
	log.Debug("SubscribeStream done")
	return nil
}

func (p *pubSubManager) Init(ctx context.Context, comp v1alpha1.Component) error {
	spec := comp.Spec
	ps, err := p.registry.Create(spec.Type, comp.ResourceVersion)
	if err != nil {
		return err
	}

	metaBase, err := p.meta.ToBaseMetadata(comp)
	if err != nil {
		return err
	}

	err = ps.Init(ctx, contribPubsub.Metadata{Base: metaBase})
	if err != nil {
		return err
	}

	p.lock.Lock()
	defer p.lock.Unlock()
	old, ok := p.compStore.Get(comp.Name)
	if ok && old.Component != nil {
		defer func() {
			_ = old.Component.Close()
		}()
	}

	item := &pubsubCompItem{
		Component: ps,
		Resource:  comp,
		closeCh:   make(chan struct{}),
	}

	if bc, ok := ps.(contribPubsub.BulkPublisher); ok {
		item.BulkPub = bc
	} else {
		item.BulkPub = &bulkPubsub{ps: ps}
	}

	if bs, ok := ps.(contribPubsub.BulkSubscriber); ok {
		item.BulkSub = bs
	}

	p.compStore.Set(comp.Name, item)

	return nil
}

func (p *pubSubManager) Close(comp v1alpha1.Component) error {
	p.lock.Lock()

	ps, ok := p.compStore.Get(comp.Name)
	if !ok {
		p.lock.Unlock()
		return nil
	}

	p.lock.Unlock()

	if !ps.closed.CompareAndSwap(false, true) {
		return nil
	}

	if err := ps.Component.Close(); err != nil {
		return err
	}

	close(ps.closeCh)

	p.compStore.Del(comp.Name)

	return nil
}
