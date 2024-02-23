package keiho

import (
	"context"
	"time"
)

type Handler interface {
	Handle(ctx context.Context, notification Notification) error
}

type Storage interface {
	Store(ctx context.Context, key string, value interface{}, shouldExpire bool) error
	Get(ctx context.Context, key string) (interface{}, error)
	Exists(ctx context.Context, key string) (bool, error)
}

type Checker interface {
	Check(ctx context.Context, since time.Time) ([]Notification, error)
}
