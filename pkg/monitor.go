package keiho

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/mitchellh/hashstructure"
)

const LAST_RUNTIME_KEY = "KEIHO_LAST_RUNTIME"

type Monitor struct {
	Checkers []Checker
	Handlers []Handler
	Storage  Storage
}

func (m *Monitor) StartWithSchedule(ctx context.Context, schedule string) error { return nil }

func (m *Monitor) Start(ctx context.Context) error {
	now := time.Now().UTC()

	// fetch last runtime from storage
	lastRuntimeTime, err := m.fetchLastRuntime(ctx)
	if err != nil {
		// TODO: log error here and inform user we assume now
		lastRuntimeTime = now
	}

	// running the checkers
	notifications, err := m.runCheckers(ctx, lastRuntimeTime)
	if err != nil {
		// TODO: log error
		log.Printf("failed to run checkers: %v", err)
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
		exists, err := m.Storage.Exists(ctx, fmt.Sprintf("%d", hash))
		if err != nil {
			// TODO: log error
			log.Printf("failed to check if hash exists in storage: %v", err)
			// assuming exists to be false due to error
			exists = false
		}

		log.Printf("hash exists in storage: %v", exists)

		if !exists {
			log.Printf("hash does not exist in storage, sending notification to handlers")
			for _, handler := range m.Handlers {
				if err := handler.Handle(ctx, notification); err != nil {
					// TODO: log error
					log.Printf("failed to handle notification: %v", err)
				}
			}

			// store notification in storage
			log.Println("storing notification in storage")
			err = m.Storage.Store(ctx, fmt.Sprintf("%d", hash), notification, true)
			if err != nil {
				// TODO: log error
				log.Printf("failed to store notification in storage: %v", err)
			}
		}
	}

	// store the execution time in storage
	log.Println("storing execution time in storage")
	nowStr := now.Format(time.RFC3339)
	err = m.Storage.Store(ctx, "KEIHO_LAST_RUNTIME", nowStr, false)
	if err != nil {
		log.Printf("failed to store execution time in storage: %v", err)
		return err
	}

	return nil
}

func (m *Monitor) fetchLastRuntime(ctx context.Context) (time.Time, error) {
	// fetch the last runtime as string from storage
	lastRuntime, err := m.Storage.Get(ctx, LAST_RUNTIME_KEY)
	if err != nil {
		return time.Time{}, err
	}

	// if we have not found a runtime string, we return an error
	if lastRuntime == nil {
		return time.Time{}, errors.New("last runtime not found")
	}

	// if we have found a runtime string, we attempt to parse it
	// first we try to cast the interface{} to a string
	lastRuntimeStr, ok := lastRuntime.(string)
	if !ok {
		return time.Time{}, errors.New("failed to parse last runtime to string")
	}

	// then we parse the string to a time.Time
	t, err := time.Parse(time.RFC3339, lastRuntimeStr)
	if err != nil {
		return time.Time{}, err
	}

	return t, nil
}

func (m *Monitor) runCheckers(ctx context.Context, since time.Time) ([]Notification, error) {
	// TODO: run checkers in parallel
	notifications := []Notification{}

	log.Println("going thorugh checkers")
	for _, checker := range m.Checkers {
		log.Printf("running checker: %v", checker)
		n, err := checker.Check(ctx, since)
		if err != nil {
			// TODO: log error
			log.Printf("failed to check: %v", err)
		}

		log.Printf("checker returned: %v", n)
		notifications = append(notifications, n...)
	}

	return notifications, nil
}
