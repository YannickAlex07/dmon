package storage

import (
	"context"
	"time"
)

type Storage interface {
	GetLatestExecutionTime(ctx context.Context) (time.Time, error)
	SetLatestExecutionTime(ctx context.Context, t time.Time) error

	IsTimeoutStored(ctx context.Context, id string) (bool, error)
	StoreTimeout(ctx context.Context, id string, t time.Time) error
}
