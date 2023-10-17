package manager

import (
	"context"
	"sync"
	"time"

	"github.com/kanengo/goutil/pkg/log"
	"go.uber.org/zap"
	"golang.org/x/sync/singleflight"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
)

type Manager struct {
	actives map[string]*Channel
	idles   map[string]IdleChannel

	appId string

	rwMutex sync.RWMutex
	idleMu  sync.RWMutex

	dialTimeout        time.Duration
	maxWriteBufferSize int
	maxReadBufferSize  int

	sg singleflight.Group
}

type IdleChannel struct {
	*Channel
	IdleTime time.Time
}

func NewManager() *Manager {
	m := &Manager{}

	if m.dialTimeout == 0 {
		m.dialTimeout = time.Second * 5
	}

	return m
}

func (m *Manager) Init(ctx context.Context) error {
	m.checkIdleChannels()
	return nil
}

func (m *Manager) checkIdleChannels() {
	go func() {
		for {
			time.Sleep(5 * time.Minute)
			var idleTimeoutTargets []string
			now := time.Now()
			m.idleMu.RLock()
			for target, idleChannel := range m.idles {
				if now.Sub(idleChannel.IdleTime) >= 5*time.Minute {
					idleTimeoutTargets = append(idleTimeoutTargets, target)
				}
			}
			m.idleMu.RUnlock()

			if len(idleTimeoutTargets) > 0 {
				deletedChannels := make([]*Channel, 0, len(idleTimeoutTargets))
				m.rwMutex.Lock()
				for _, target := range idleTimeoutTargets {
					c, ok := m.actives[target]
					if !ok {
						continue
					}
					deletedChannels = append(deletedChannels, c)
					delete(m.actives, target)
				}
				m.rwMutex.Unlock()
			}
		}

	}()
}

func (m *Manager) newChannel(ctx context.Context, target string) (*Channel, error) {
	v, err, _ := m.sg.Do(target, func() (interface{}, error) {
		channel := &Channel{
			appId:  m.appId,
			target: target,
		}

		opts := []grpc.DialOption{
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithBlock(),
			grpc.WithIdleTimeout(time.Minute * 5),
		}

		if m.maxWriteBufferSize != 0 {
			opts = append(opts, grpc.WithWriteBufferSize(m.maxWriteBufferSize<<20))
		}

		if m.maxWriteBufferSize != 0 {
			opts = append(opts, grpc.WithReadBufferSize(m.maxWriteBufferSize<<20))
		}

		dialCtx, cancel := context.WithTimeout(ctx, m.dialTimeout)
		cli, err := grpc.DialContext(dialCtx, target, opts...)
		cancel()
		if err != nil {
			return nil, err
		}

		channel.cli = cli
		if channel.cli.GetState() == connectivity.Ready {
			m.rwMutex.Lock()
			m.actives[channel.target] = channel
			m.rwMutex.Unlock()

			go func() {
				for {
					b := channel.cli.WaitForStateChange(ctx, channel.cli.GetState())
					if !b {
						return
					}
					state := channel.cli.GetState()
					log.Debug("[gRPC]channel state change", zap.Any("state", state))
					if state == connectivity.Idle {
						idleChannel := IdleChannel{
							Channel:  channel,
							IdleTime: time.Now(),
						}
						m.idleMu.Lock()
						m.idles[channel.target] = idleChannel
						m.idleMu.Unlock()
					} else if state == connectivity.Ready {
						m.rwMutex.Lock()
						m.actives[channel.target] = channel
						m.rwMutex.Unlock()

						m.idleMu.RLock()
						_, ok := m.idles[channel.target]
						m.idleMu.RUnlock()

						if ok {
							m.idleMu.Lock()
							delete(m.idles, channel.target)
							m.idleMu.Unlock()
						}
					} else {
						m.idleMu.RLock()
						_, ok := m.idles[channel.target]
						m.idleMu.RUnlock()

						if ok {
							m.idleMu.Lock()
							delete(m.idles, channel.target)
							m.idleMu.Unlock()
						}

						m.rwMutex.RLock()
						_, ok = m.actives[channel.target]
						m.rwMutex.RUnlock()

						if ok {
							m.rwMutex.Lock()
							delete(m.actives, channel.target)
							m.rwMutex.Unlock()
						}
						break
					}
				}
			}()
		}

		return channel, nil
	})

	if err != nil {
		return nil, err
	}

	return v.(*Channel), nil
}
