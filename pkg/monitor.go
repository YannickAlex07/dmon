package keiho

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"

	"github.com/go-co-op/gocron"
)

const LAST_RUNTIME_KEY = "KEIHO_LAST_RUNTIME"

type Monitor struct {
	Checkers []Checker
	Handlers []Handler
	Storage  Storage
}

func (m *Monitor) StartWithSchedule(ctx context.Context, schedule string) error {
	// monitor func
	scheduler := gocron.NewScheduler(time.UTC)

	scheduler.Cron(schedule).Do(func() {
		m.Start(ctx)
	})

	scheduler.StartBlocking()

	return nil
}

func (m *Monitor) Start(ctx context.Context) error {
	now := time.Now().UTC()

	// fetch last runtime from storage
	lastRuntimeTime, err := m.fetchLastRuntime(ctx)
	if err != nil {
		// TODO: log error here and inform user we assume earliest time
		lastRuntimeTime = now
	}

	// running the checkers
	notifications, err := m.runCheckers(ctx, lastRuntimeTime)
	if err != nil {
		// TODO: log error
		log.Printf("failed to run checkers: %v", err)
	}

	// running the handlers
	err = m.runHandlers(ctx, notifications)
	if err != nil {
		log.Printf("failed to run handlers: %v", err)
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
	var wg sync.WaitGroup
	resultsChan := make(chan Notification)

	for _, checker := range m.Checkers {
		wg.Add(1)

		go func(c Checker) {
			defer wg.Done()

			// run the checker
			notifications, err := c.Check(ctx, since)
			if err != nil {
				log.Println(err)
				return
			}

			// append the notifications
			for _, n := range notifications {
				resultsChan <- n
			}
		}(checker)
	}

	// close the channel when all checkers are done
	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	// collect the results
	notifications := []Notification{}
	for n := range resultsChan {
		notifications = append(notifications, n)
	}

	return notifications, nil
}

func (m *Monitor) runHandlers(ctx context.Context, notifications []Notification) error {
	for _, notification := range notifications {
		// check if the hash exists in the storage
		log.Printf("checking if notification exists in storage - key: %s", notification.Key)
		exists, err := m.Storage.Exists(ctx, notification.Key)
		if err != nil {
			// TODO: log error
			log.Printf("failed to check if hash exists in storage: %v", err)
			// assuming exists to be false due to error
			exists = false
		}

		log.Printf("notification exists in storage: %v", exists)

		if !exists {
			var wg sync.WaitGroup

			log.Printf("notification does not exist in storage, sending notification to handlers")
			for _, handler := range m.Handlers {
				wg.Add(1)

				go func(h Handler) {
					defer wg.Done()
					if err := h.Handle(ctx, notification); err != nil {
						log.Printf("failed to handle notification: %v", err)
					}
				}(handler)
			}

			wg.Wait()

			// store notification in storage
			log.Println("storing notification in storage")
			err = m.Storage.Store(ctx, notification.Key, notification, true)
			if err != nil {
				// TODO: log error
				log.Printf("failed to store notification in storage: %v", err)
			}
		}
	}

	return nil
}
