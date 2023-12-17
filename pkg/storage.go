package siren

import (
	"context"
)

type Storage interface {
	Store(ctx context.Context, key string, value interface{}, shouldExpire bool) error
	Get(ctx context.Context, key string) (interface{}, error)
	Exists(ctx context.Context, key string) (bool, error)
}
