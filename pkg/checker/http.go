package checker

import (
	"context"
	"net/url"
	"time"

	siren "github.com/yannickalex07/dmon/pkg"
)

type HTTPChecker struct {
	url url.URL
}

func (c HTTPChecker) Check(ctx context.Context, since time.Time) ([]siren.Notification, error) {
	println(c.url.String())
	return nil, nil
}
