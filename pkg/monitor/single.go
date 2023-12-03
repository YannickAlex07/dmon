package monitor

import (
	"context"
	"time"

	siren "github.com/yannickalex07/dmon/pkg"
	"google.golang.org/appengine/log"
)

type SingleMonitor struct{}

func Start(ctx context.Context, checkers []siren.Checker, handlers []siren.Handler, storage siren.Storage) error {
	notifications := []siren.Notification{}

	// running the checkers
	for _, checker := range checkers {
		n, err := checker.Check(ctx, time.Now())
		if err != nil {
			log.Errorf(ctx, "error while checking: %s", err)
		}

		notifications = append(notifications, n...)
	}

	// sending to notifications to the handlers
	for _, notification := range notifications {
		// TODO: filter down notifications based on the storage
		// 1. hash the notification using -> https://pkg.go.dev/github.com/mitchellh/hashstructure
		// 2. check if the hash exists in the storage

		for _, handler := range handlers {
			if err := handler.Handle(ctx, notification); err != nil {
				log.Errorf(ctx, "error while handling: %s", err)
			}
		}
	}

	return nil
}
