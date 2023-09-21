package components

import "sync"

type CompStore[T any] struct {
	comps map[string]T
	rw    sync.RWMutex
}

func NewCompStore[T any]() *CompStore[T] {
	cs := &CompStore[T]{
		comps: make(map[string]T),
	}

	return cs
}

func (c *CompStore[T]) Get(name string) (T, bool) {
	c.rw.RLock()
	defer c.rw.RUnlock()

	val, ok := c.comps[name]

	return val, ok
}

func (c *CompStore[T]) Set(name string, val T) {
	c.rw.Lock()
	defer c.rw.Unlock()
	c.comps[name] = val
}

func (c *CompStore[T]) Del(name string) {
	c.rw.Lock()
	defer c.rw.Unlock()
	delete(c.comps, name)
}
