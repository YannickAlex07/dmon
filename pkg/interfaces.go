package inframon

import (
	"context"
	"time"
)

type Destination interface {
	Handle(ctx context.Context, notification Notification) error
}

type State interface {
	Store(ctx context.Context, key string, value interface{}, shouldExpire bool) error
	Get(ctx context.Context, key string) (interface{}, error)
	Exists(ctx context.Context, key string) (bool, error)
}

type Source interface {
	Check(ctx context.Context, since time.Time) ([]Notification, error)
}
