package gcp

import (
	"context"

	gmon "github.com/yannickalex07/dmon/pkg"
)

type StackdriverChecker struct{}

func (c StackdriverChecker) Check(ctx context.Context, handlers []gmon.Handler) error {
	return nil
}
