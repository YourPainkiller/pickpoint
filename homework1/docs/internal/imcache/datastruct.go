package imcache

import "time"

func NewCached[V any](expiredAt time.Time, value V) *Cached[V] {
	return &Cached[V]{
		expiredAt: expiredAt,
		value:     value,
	}
}

type Cached[V any] struct {
	expiredAt time.Time
	value     V
}

func (c *Cached[V]) Expired(now time.Time) bool {
	return c.expiredAt.Before(now)
}

func (c *Cached[V]) Value() V {
	return c.value
}
