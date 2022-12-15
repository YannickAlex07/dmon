package storage

import (
	"time"

	"github.com/jellydator/ttlcache/v3"
)

type MemoryStorage struct {
	cache       *ttlcache.Cache[string, string]
	lastRunTime time.Time
}

func NewMemoryStore(ttl time.Duration) *MemoryStorage {
	cache := ttlcache.New(
		ttlcache.WithTTL[string, string](ttl),
	)

	go cache.Start()

	return &MemoryStorage{
		cache:       cache,
		lastRunTime: time.Now().UTC(),
	}
}

func (s MemoryStorage) GetLatestRuntime() time.Time {
	return s.lastRunTime
}

func (s *MemoryStorage) SetLatestRuntime(newRunTime time.Time) {
	s.lastRunTime = newRunTime
}

func (s MemoryStorage) TimeoutAlreadyHandled(id string) bool {
	item := s.cache.Get(id)
	return item != nil
}

func (s *MemoryStorage) TimeoutHandled(id string) {
	s.cache.Set(id, time.Now().UTC().String(), ttlcache.DefaultTTL)
}
