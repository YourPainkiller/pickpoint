package imcache

import (
	"container/list"
	"sync"
	"time"
)

func NewTTLClient[K comparable, V any](ttl time.Duration, size int) *TTLClient[K, V] {
	return &TTLClient[K, V]{
		ttl:  ttl,
		data: make(map[K]*Cached[V]),
		ll:   list.New(),
		size: size,
	}
}

type TTLClient[K comparable, V any] struct {
	ttl  time.Duration
	lock sync.RWMutex
	data map[K]*Cached[V]
	ll   *list.List
	size int
}

func (c *TTLClient[K, V]) Get(key K) (V, bool) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	v, ok := c.data[key]

	if ok && !v.Expired(time.Now()) {
		c.ll.MoveToFront(v.keyLink)
		return v.Value(), true
	}
	return (&Cached[V]{}).Value(), false
}

func (c *TTLClient[K, V]) Set(key K, value V, now time.Time) {
	c.lock.Lock()
	defer c.lock.Unlock()

	v, ok := c.data[key]
	if ok {
		c.ll.MoveToFront(v.keyLink)
		c.data[key].expiredAt = now.Add(c.ttl)
	} else {
		if c.ll.Len() == c.size {
			backElem := c.ll.Back()
			c.ll.Remove(backElem)
			k := backElem.Value.(K)
			delete(c.data, k)
		}
		link := c.ll.PushFront(list.Element{Value: key})
		wrapped := NewCached[V](now.Add(c.ttl), link, value)
		c.data[key] = wrapped
	}

}
