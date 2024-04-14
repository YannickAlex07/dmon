package local

import (
	"context"
	"errors"
	"time"

	"github.com/jellydator/ttlcache/v3"
)

type MemoryState struct {
	// The TTL for every notification that the state will store
	cache *ttlcache.Cache[string, interface{}]
}

func NewMemoryState(ttl time.Duration) MemoryState {
	cache := ttlcache.New(
		ttlcache.WithTTL[string, interface{}](ttl),
	)

	go cache.Start()

	return MemoryState{
		cache: cache,
	}
}

func (m MemoryState) Store(ctx context.Context, key string, value interface{}, shouldExpire bool) error {
	ttl := ttlcache.NoTTL
	if shouldExpire {
		ttl = ttlcache.DefaultTTL
	}

	m.cache.Set(key, value, ttl)
	return nil
}

func (m MemoryState) Get(ctx context.Context, key string) (interface{}, error) {
	hasKey := m.cache.Has(key)

	if !hasKey {
		return nil, errors.New("key not found")
	}

	item := m.cache.Get(key).Value()
	return item, nil
}

func (m MemoryState) Exists(ctx context.Context, key string) (bool, error) {
	return m.cache.Has(key), nil
}
