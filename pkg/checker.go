package siren

import (
	"context"
	"time"
)

type Checker interface {
	Check(ctx context.Context, since time.Time) ([]Notification, error)
}
