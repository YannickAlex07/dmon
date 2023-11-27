package storage

import "time"

type MemoryStorage struct {
	// The TTL for every notification that the state will store
	TTL time.Duration
}
