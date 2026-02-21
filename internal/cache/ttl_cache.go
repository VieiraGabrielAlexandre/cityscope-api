package cache

import (
	"sync"
	"time"

	"golang.org/x/sync/singleflight"
)

type TTLCache[V any] struct {
	mu    sync.RWMutex
	items map[string]entry[V]
	sf    singleflight.Group
	clock func() time.Time
}

type entry[V any] struct {
	value     V
	expiresAt time.Time
}

func NewTTLCache[V any]() *TTLCache[V] {
	return &TTLCache[V]{
		items: make(map[string]entry[V]),
		clock: time.Now,
	}
}

func (c *TTLCache[V]) Get(key string) (V, bool) {
	c.mu.RLock()
	e, ok := c.items[key]
	c.mu.RUnlock()

	var zero V
	if !ok {
		return zero, false
	}
	if c.clock().After(e.expiresAt) {
		c.mu.Lock()
		e2, ok2 := c.items[key]
		if ok2 && c.clock().After(e2.expiresAt) {
			delete(c.items, key)
		}
		c.mu.Unlock()
		return zero, false
	}
	return e.value, true
}

func (c *TTLCache[V]) Set(key string, value V, ttl time.Duration) {
	if ttl <= 0 {
		return
	}
	c.mu.Lock()
	c.items[key] = entry[V]{value: value, expiresAt: c.clock().Add(ttl)}
	c.mu.Unlock()
}

func (c *TTLCache[V]) GetOrSet(key string, ttl time.Duration, loader func() (V, error)) (V, error) {
	if v, ok := c.Get(key); ok {
		return v, nil
	}

	vAny, err, _ := c.sf.Do(key, func() (any, error) {
		if v, ok := c.Get(key); ok {
			return v, nil
		}
		v, err := loader()
		if err != nil {
			var zero V
			return zero, err
		}
		c.Set(key, v, ttl)
		return v, nil
	})

	if err != nil {
		var zero V
		return zero, err
	}
	return vAny.(V), nil
}
