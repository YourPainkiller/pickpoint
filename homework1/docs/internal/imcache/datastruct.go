package imcache

import (
	"container/list"
	"time"
)

func NewCached[V any](expiredAt time.Time, keyLink *list.Element, value V) *Cached[V] {
	return &Cached[V]{
		expiredAt: expiredAt,
		keyLink:   keyLink,
		value:     value,
	}
}

type Cached[V any] struct {
	expiredAt time.Time
	keyLink   *list.Element
	value     V
}

func (c *Cached[V]) Expired(now time.Time) bool {
	return c.expiredAt.Before(now)
}

func (c *Cached[V]) Value() V {
	return c.value
}

func (c Cached[V]) Link() *list.Element {
	return c.keyLink
}
