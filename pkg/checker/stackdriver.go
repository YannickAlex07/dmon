package checker

import (
	"context"
	"time"

	keiho "github.com/yannickalex07/dmon/pkg"
)

type StackdriverChecker struct{}

func (c StackdriverChecker) Check(ctx context.Context, since time.Time) ([]keiho.Notification, error) {
	return nil, nil
}
