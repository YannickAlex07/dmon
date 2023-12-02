package checker

import (
	"context"
	"time"

	siren "github.com/yannickalex07/dmon/pkg"
)

type StackdriverChecker struct{}

func (c StackdriverChecker) Check(ctx context.Context, since time.Time) ([]siren.Notification, error) {
	return nil, nil
}
