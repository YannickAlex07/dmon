package storage

import (
	"context"
	"time"

	"github.com/jellydator/ttlcache/v3"
)

type MemoryStorage struct {
	// The TTL for every notification that the state will store
	cache *ttlcache.Cache[string, interface{}]
}

func NewMemoryStorage(ttl time.Duration) MemoryStorage {
	cache := ttlcache.New(
		ttlcache.WithTTL[string, interface{}](ttl),
	)

	go cache.Start()

	return MemoryStorage{
		cache: cache,
	}
}

func (m MemoryStorage) Store(ctx context.Context, key string, value interface{}, shouldExpire bool) error {
	ttl := ttlcache.NoTTL
	if shouldExpire {
		ttl = ttlcache.DefaultTTL
	}

	m.cache.Set(key, value, ttl)
	return nil
}

func (m MemoryStorage) Get(ctx context.Context, key string) (interface{}, error) {
	hasKey := m.cache.Has(key)

	if !hasKey {
		return nil, nil
	}

	item := m.cache.Get(key).Value()
	return item, nil
}

func (m MemoryStorage) Exists(ctx context.Context, key string) (bool, error) {
	return m.cache.Has(key), nil
}
