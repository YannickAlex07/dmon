package monitor

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/mitchellh/hashstructure"
	keiho "github.com/yannickalex07/dmon/pkg"
)

type SingleMonitor struct{}

func (SingleMonitor) Start(ctx context.Context, checkers []keiho.Checker, handlers []keiho.Handler, storage keiho.Storage) error {
	notifications := []keiho.Notification{}

	now := time.Now().UTC()

	// fetch last runtime from storage
	lastRuntime, err := storage.Get(ctx, "KEIHO_LAST_RUNTIME")
	if err != nil {
		return err
	}

	// TODO: move to seperate method
	lastRuntimeStr, ok := lastRuntime.(string)
	if !ok {
		return errors.New("failed to parse last runtime to string")
	}

	lastRuntimeTime, err := time.Parse(time.RFC3339, lastRuntimeStr)
	if err != nil {
		return err
	}

	// running the checkers
	for _, checker := range checkers {
		n, err := checker.Check(ctx, lastRuntimeTime)
		if err != nil {
			// TODO: log error
		}

		notifications = append(notifications, n...)
	}

	// sending to notifications to the handlers
	for _, notification := range notifications {
		// TODO: filter down notifications based on the storage
		// 1. hash the notification using -> https://pkg.go.dev/github.com/mitchellh/hashstructure
		hash, err := hashstructure.Hash(notification, nil)
		if err != nil {
			// TODO: log error
		}

		// 2. check if the hash exists in the storage
		exists, err := storage.Exists(ctx, fmt.Sprintf("%d", hash))
		if err != nil {
			// TODO: log error
			// assuming exists to be false due to error
			exists = false
		}

		if !exists {
			for _, handler := range handlers {
				if err := handler.Handle(ctx, notification); err != nil {
					// TODO: log error
				}
			}

			// store notification in storage
			err = storage.Store(ctx, fmt.Sprintf("%d", hash), notification, true)
			if err != nil {
				// TODO: log error
			}
		}

	}

	// store the execution time in storage
	nowStr := now.Format(time.RFC3339)
	err = storage.Store(ctx, "KEIHO_LAST_RUNTIME", nowStr, false)
	if err != nil {
		return err
	}

	return nil
}
