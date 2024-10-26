package imcache

import (
	"sync"
	"time"
)

func NewTTLClient[K comparable, V any](ttl time.Duration) *TTLClient[K, V] {
	return &TTLClient[K, V]{
		ttl:  ttl,
		data: make(map[K]*Cached[V]),
	}
}

type TTLClient[K comparable, V any] struct {
	ttl  time.Duration
	lock sync.RWMutex
	data map[K]*Cached[V]
}

func (c *TTLClient[K, V]) Get(key K) (V, bool) {
	c.lock.RLock()
	v, ok := c.data[key]
	c.lock.RUnlock()

	if ok && !v.Expired(time.Now()) {
		return v.Value(), true
	}

	return (&Cached[V]{}).Value(), false
}

func (c *TTLClient[K, V]) Set(key K, value V, now time.Time) {
	wrapped := NewCached[V](now.Add(c.ttl), value)
	c.lock.Lock()
	c.data[key] = wrapped
	c.lock.Unlock()
}
