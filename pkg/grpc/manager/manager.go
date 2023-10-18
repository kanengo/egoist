package manager

import (
	"context"
	"sync"
	"time"

	"github.com/kanengo/goutil/pkg/log"
	"github.com/kanengo/goutil/pkg/utils"
	"go.uber.org/zap"
	"golang.org/x/sync/singleflight"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
)

type Manager struct {
	actives map[string]*Channel
	idles   map[string]idleChannel

	appId string

	rwMutex sync.RWMutex
	idleMu  sync.RWMutex

	dialTimeout        time.Duration
	maxWriteBufferSize int
	maxReadBufferSize  int

	sg singleflight.Group

	ctx    context.Context
	cancel context.CancelFunc

	wg sync.WaitGroup
}

type idleChannel struct {
	*Channel
	IdleTime time.Time
}

func NewManager(ctx context.Context) *Manager {
	var cancel context.CancelFunc
	ctx, cancel = context.WithCancel(ctx)
	m := &Manager{
		ctx:    ctx,
		cancel: cancel,

		actives: map[string]*Channel{},
		idles:   map[string]idleChannel{},
	}

	if m.dialTimeout == 0 {
		m.dialTimeout = time.Second * 5
	}

	return m
}

func (m *Manager) Init(ctx context.Context) error {
	m.checkIdleChannels()
	return nil
}

func (m *Manager) Close() error {
	m.cancel()

	m.wg.Wait()

	return nil
}

func (m *Manager) GetChannel(ctx context.Context, target string) (c *Channel, err error) {
	return m.getActiveChannel(ctx, target)
}

func (m *Manager) getActiveChannel(ctx context.Context, target string) (c *Channel, err error) {
	var ok bool
	m.rwMutex.RLock()
	c, ok = m.actives[target]
	m.rwMutex.RUnlock()

	if ok {
		return c, nil
	}

	c, err = m.newChannel(ctx, target)

	if err != nil {
		return nil, err
	}

	return c, nil
}

func (m *Manager) hasActiveTarget(target string) bool {
	m.rwMutex.RLock()
	_, ok := m.actives[target]
	m.rwMutex.RUnlock()

	return ok
}
func (m *Manager) hasIdleTarget(target string) bool {
	m.idleMu.RLock()
	_, ok := m.idles[target]
	m.idleMu.RUnlock()

	return ok
}

func (m *Manager) addIdleChannel(ic idleChannel) {
	m.idleMu.Lock()
	m.idles[ic.Channel.target] = ic
	m.idleMu.Unlock()
}

func (m *Manager) delIdleChannel(target string) {
	if !m.hasIdleTarget(target) {
		return
	}

	m.idleMu.Lock()
	delete(m.idles, target)
	m.idleMu.Unlock()
}

func (m *Manager) delActiveChannel(target string) {
	if !m.hasActiveTarget(target) {
		return
	}

	m.rwMutex.Lock()
	delete(m.actives, target)
	m.rwMutex.Unlock()
}

func (m *Manager) addActiveChannel(c *Channel) {
	if m.hasActiveTarget(c.target) {
		return
	}
	m.rwMutex.Lock()
	m.actives[c.target] = c
	m.rwMutex.Unlock()
}

func (m *Manager) checkIdleChannels() {
	m.wg.Add(1)
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer func() {
			ticker.Stop()
			m.wg.Done()
		}()
		for {
			select {
			case <-ticker.C:
			case <-m.ctx.Done():
				break
			}

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
				deletedChannels := make(map[string]*Channel, len(idleTimeoutTargets))
				m.rwMutex.Lock()
				for _, target := range idleTimeoutTargets {
					c, ok := m.actives[target]
					if !ok {
						continue
					}
					deletedChannels[c.target] = c
					delete(m.actives, target)
				}
				m.rwMutex.Unlock()

				m.idleMu.Lock()
				for _, target := range idleTimeoutTargets {
					ic, ok := m.idles[target]
					if !ok {
						continue
					}

					if _, ok := deletedChannels[target]; !ok {
						deletedChannels[target] = ic.Channel
					}
					delete(m.idles, target)
				}
				m.idleMu.Unlock()

				for _, c := range deletedChannels {
					_ = c.Close()
				}
			}
		}
	}()
}

func (m *Manager) checkChannelState(c *Channel) {
	m.wg.Add(1)
	go func() {
		utils.CheckGoPanic(m.ctx, nil)
		defer func() {
			_ = c.Close()
			m.wg.Done()
		}()
		for {
			b := c.cli.WaitForStateChange(m.ctx, c.cli.GetState())
			if !b {
				break
			}
			state := c.cli.GetState()
			log.Debug("[gRPC]channel state change", zap.Any("state", state))
			if state == connectivity.Idle {
				ic := idleChannel{
					Channel:  c,
					IdleTime: time.Now(),
				}
				m.addIdleChannel(ic)
			} else if state == connectivity.Ready {
				m.addActiveChannel(c)

				m.delIdleChannel(c.target)
			} else {
				//
				m.delActiveChannel(c.target)

				m.delIdleChannel(c.target)

				break
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

			m.checkChannelState(channel)
		}

		return channel, nil
	})

	if err != nil {
		return nil, err
	}

	return v.(*Channel), nil
}
