package monitor

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/mitchellh/hashstructure"
	keiho "github.com/yannickalex07/dmon/pkg"
)

type SingleMonitor struct{}

func (SingleMonitor) Start(ctx context.Context, checkers []keiho.Checker, handlers []keiho.Handler, storage keiho.Storage) error {
	notifications := []keiho.Notification{}

	log.Println("running monitor")

	now := time.Now().UTC()

	// fetch last runtime from storage
	log.Println("fetching last runtime")
	lastRuntime, err := storage.Get(ctx, "KEIHO_LAST_RUNTIME")
	if err != nil {
		return err
	}
	log.Printf("last runtime: %v", lastRuntime)

	lastRuntimeTime := now
	if lastRuntime != nil {
		log.Println("parsing last runtime")

		// TODO: move to seperate method
		lastRuntimeStr, ok := lastRuntime.(string)
		if !ok {
			return errors.New("failed to parse last runtime to string")
		}

		lastRuntimeTime, err = time.Parse(time.RFC3339, lastRuntimeStr)
		if err != nil {
			log.Printf("failed to parse last runtime: %v", err)
			return err
		}

		log.Println("parsed last runtime")
	}

	// running the checkers
	log.Println("going thorugh checkers")
	for _, checker := range checkers {
		log.Printf("running checker: %v", checker)
		n, err := checker.Check(ctx, lastRuntimeTime)
		if err != nil {
			// TODO: log error
			log.Printf("failed to check: %v", err)
		}

		log.Printf("checker returned: %v", n)
		notifications = append(notifications, n...)
	}

	log.Println("going through notifications")
	// sending to notifications to the handlers
	for _, notification := range notifications {
		// TODO: filter down notifications based on the storage
		// 1. hash the notification using -> https://pkg.go.dev/github.com/mitchellh/hashstructure
		log.Println("hashing notification")
		hash, err := hashstructure.Hash(notification, nil)
		if err != nil {
			// TODO: log error
			log.Printf("failed to hash notification: %v", err)
		}

		log.Printf("hashed notification: %d", hash)

		// 2. check if the hash exists in the storage
		log.Println("checking if hash exists in storage")
		exists, err := storage.Exists(ctx, fmt.Sprintf("%d", hash))
		if err != nil {
			// TODO: log error
			log.Printf("failed to check if hash exists in storage: %v", err)
			// assuming exists to be false due to error
			exists = false
		}

		log.Printf("hash exists in storage: %v", exists)

		if !exists {
			log.Printf("hash does not exist in storage, sending notification to handlers")
			for _, handler := range handlers {
				if err := handler.Handle(ctx, notification); err != nil {
					// TODO: log error
					log.Printf("failed to handle notification: %v", err)
				}
			}

			// store notification in storage
			log.Println("storing notification in storage")
			err = storage.Store(ctx, fmt.Sprintf("%d", hash), notification, true)
			if err != nil {
				// TODO: log error
				log.Printf("failed to store notification in storage: %v", err)
			}
		}
	}

	// store the execution time in storage
	log.Println("storing execution time in storage")
	nowStr := now.Format(time.RFC3339)
	err = storage.Store(ctx, "KEIHO_LAST_RUNTIME", nowStr, false)
	if err != nil {
		log.Printf("failed to store execution time in storage: %v", err)
		return err
	}

	log.Println("finished running monitor")
	return nil
}
