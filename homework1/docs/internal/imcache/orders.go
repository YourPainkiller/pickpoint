package imcache

import (
	"context"
	"homework1/internal/dto"
	"time"

	"github.com/opentracing/opentracing-go"
)

func NewOrdersCache(ttl time.Duration) *OrdersCache {
	return &OrdersCache{
		cli: NewTTLClient[int, *dto.OrderDto](ttl),
	}
}

type OrdersCache struct {
	cli *TTLClient[int, *dto.OrderDto]
}

func (p *OrdersCache) Get(ctx context.Context, key int) (*dto.OrderDto, bool) {
	span, _ := opentracing.StartSpanFromContext(ctx, "service.cache.get")
	defer span.Finish()

	return p.cli.Get(key)
}

func (p *OrdersCache) Set(ctx context.Context, key int, order *dto.OrderDto, now time.Time) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "service.cache.set")
	defer span.Finish()

	p.cli.Set(key, order, now)
	return nil
}
