package checker

import (
	"context"
	"net/url"
	"time"

	keiho "github.com/yannickalex07/dmon/pkg"
)

type HTTPChecker struct {
	url url.URL
}

func (c HTTPChecker) Check(ctx context.Context, since time.Time) ([]keiho.Notification, error) {
	println(c.url.String())
	return nil, nil
}
