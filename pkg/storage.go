package gmon

import (
	"context"
)

type Storage interface {
	Store(ctx context.Context, key string, value interface{}) error
	Get(ctx context.Context, key string) (interface{}, error)
}
