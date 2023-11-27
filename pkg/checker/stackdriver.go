package checker

import (
	"context"
	"time"
)

type StackdriverChecker struct{}

func (c StackdriverChecker) Check(ctx context.Context, since time.Time) error {
	return nil
}
