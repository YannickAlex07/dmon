package gmon

import (
	"context"
	"time"
)

type State interface {
	Store(ctx context.Context, key string, value interface{}) error
	StoreWithTTL(ctx context.Context, key string, value interface{}, duration time.Duration) error

	Get(ctx context.Context, key string) (interface{}, error)
}
