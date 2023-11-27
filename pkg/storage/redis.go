package storage

import "time"

type RedisStorage struct {
	// The TTL for every notification that the state will store
	TTL time.Duration
}
