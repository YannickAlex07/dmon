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

func (s MemoryStorage) GetLatestExecutionTime() (time.Time, error) {
	return s.lastRunTime, nil
}

func (s *MemoryStorage) SetLatestExecutionTime(t time.Time) {
	s.lastRunTime = t
}

func (s MemoryStorage) WasTimeoutHandled(id string) bool {
	item := s.cache.Get(id)
	return item != nil
}

func (s *MemoryStorage) HandleTimeout(id string, t time.Time) {
	s.cache.Set(id, t.Format(time.RFC3339), ttlcache.DefaultTTL)
}
