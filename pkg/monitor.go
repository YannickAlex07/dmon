package inframon

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/go-co-op/gocron/v2"
	log "github.com/sirupsen/logrus"
)

// The last runtime will be stored under this key
const LAST_RUNTIME_KEY = "INFRAMON_LAST_RUNTIME"

// The monitor can be used to observe a number of source and in case of any issues
// notify a set of destinations.
//
// Monitor can be started once with `Start()` or on a cron schedule with `StartWithSchedule()`.
type Monitor struct {
	Sources      []Source
	Destinations []Destination
	State        State
}

// This will start the monitor on a specified cron schedule. This will essentially run forever as
// long as the context is not cancelled.
//
// If you want to end the monitor, just cancel the context.
func (m *Monitor) StartWithSchedule(ctx context.Context, schedule string) error {
	// create our cron scheduler
	scheduler, err := gocron.NewScheduler(gocron.WithLocation(time.UTC))
	if err != nil {
		log.Errorf("Failed to create cron scheduler: %v", err)
		return err
	}

	// we create a new job that is automatically registered
	// with the scheduler.
	_, err = scheduler.NewJob(
		gocron.CronJob(schedule, false),
		gocron.NewTask(m.Start, ctx),
	)
	if err != nil {
		log.Errorf("Failed to create cron job: %v", err)
		return err
	}

	// we will start the scheduler and let it run until the context is cancelled
	log.Infof("Starting scheduler with schedule %s", schedule)
	scheduler.Start()
	<-ctx.Done()

	// once the context is cancelled, we will shutdown the scheduler
	log.Info("Context finished, shutting down scheduler")
	err = scheduler.Shutdown()
	if err != nil {
		log.Errorf("Failed to shutdown scheduler: %v", err)
		return err
	}

	log.Infof("Scheduler shutted down!")
	return nil
}

// This will start the monitor once. This will check each source for any notifications
// and will then forward any notifications to all destinations.
func (m *Monitor) Start(ctx context.Context) error {
	// get the current time
	now := time.Now().UTC()

	// fetch last runtime from state
	log.Info("Fetching last runtime from state")
	lastRuntimeTime, err := m.fetchLastRuntime(ctx)
	if err != nil {
		log.Warnf("Assuming now() as the last runtime due to error: %v", err)
		lastRuntimeTime = now
	}

	contextLogger := log.WithField("last_runtime", lastRuntimeTime)

	// running the checkers
	contextLogger.Info("Checking sources...")
	notifications, err := m.checkSources(ctx, lastRuntimeTime)
	if err != nil {
		contextLogger.Errorf("Source failed: %v", err)
	}

	// running the handlers
	contextLogger.Infof("Notifying destinations about %d notifications...", len(notifications))
	err = m.notifyDestinations(ctx, notifications)
	if err != nil {
		contextLogger.Errorf("Failed to notify destination: %v", err)
	}

	// store the execution time in storage
	nowStr := now.Format(time.RFC3339)
	contextLogger = contextLogger.WithField("new_runtime", nowStr)
	contextLogger.Info("Storing execution time in storage...")

	err = m.State.Store(ctx, LAST_RUNTIME_KEY, nowStr, false)
	if err != nil {
		contextLogger.Errorf("Failed to store execution time in storage: %v", err)
		return err
	}

	contextLogger.Info("Finished current monitor run")
	return nil
}

// This function will fetch the last runtime from the state. If the last runtime is not found
// or couldn't be parsed, an error will be returned. This usually will lead the monitor to assume
// now as the last runtime.
//
// The runtime needs to be in RFC3339 format.
func (m *Monitor) fetchLastRuntime(ctx context.Context) (time.Time, error) {
	// fetch the last runtime as string from storage
	lastRuntime, err := m.State.Get(ctx, LAST_RUNTIME_KEY)
	if err != nil {
		log.Errorf("Failed to get last runtime from storage: %v", err)
		return time.Time{}, err
	}

	// if we have not found a runtime string, we return an error
	if lastRuntime == nil {
		log.Warn("Last runtime is nil")
		return time.Time{}, errors.New("last runtime is nil")
	}

	// if we have found a runtime string, we attempt to parse it
	// first we try to cast the interface{} to a string
	lastRuntimeStr, ok := lastRuntime.(string)
	if !ok {
		log.Errorf("Failed to parse last runtime to string: %+v", lastRuntime)
		return time.Time{}, errors.New("failed to parse last runtime to string")
	}

	// then we parse the string to a time.Time
	t, err := time.Parse(time.RFC3339, lastRuntimeStr)
	if err != nil {
		log.Errorf("Failed to parse last runtime to time.Time: %v", err)
		return time.Time{}, err
	}

	return t, nil
}

// Check all sources for any notifications since the specified time.
func (m *Monitor) checkSources(ctx context.Context, since time.Time) ([]Notification, error) {
	var wg sync.WaitGroup
	resultsChan := make(chan Notification)

	for _, checker := range m.Sources {
		wg.Add(1)

		go func(s Source) {
			defer wg.Done()

			log.Debugf("Checking source: %#v", s)

			// run the checker
			notifications, err := s.Check(ctx, since)
			if err != nil {
				log.Errorf("Source failed: %v", err)
				return
			}

			// append the notifications
			log.Debugf("Found %d notifications for source %#v", len(notifications), s)
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

	log.Debugf("Found a total of %d notifications", len(notifications))

	return notifications, nil
}

// Notify all destinations about new notifications found by the sources.
func (m *Monitor) notifyDestinations(ctx context.Context, notifications []Notification) error {
	for _, notification := range notifications {
		// check if the hash exists in the storage
		contextLogger := log.WithField("notification_key", notification.Key)
		contextLogger.Debug("Checking if notification exists in storage...")

		exists, err := m.State.Exists(ctx, notification.Key)
		if err != nil {
			contextLogger.Errorf("Failed to storage for key: %v", err)

			// assuming exists to be false due to error
			contextLogger.Warn("Assuming notification does not exist in storage due to error")
			exists = false
		}

		if !exists {
			var wg sync.WaitGroup

			contextLogger.Debug("Notification does not exist in storage")
			for _, destination := range m.Destinations {
				wg.Add(1)

				go func(d Destination) {
					defer wg.Done()

					log.Debugf("Notifying destination: %#v", d)

					if err := d.Handle(ctx, notification); err != nil {
						contextLogger.Errorf("Failed to handle notification: %v", err)
					}
				}(destination)
			}

			wg.Wait()

			// store notification in storage
			log.Debugf("Storing notification in storage...")
			err = m.State.Store(ctx, notification.Key, notification, true)
			if err != nil {
				log.Errorf("Failed to store notification in storage: %v", err)
			}
		}
	}

	return nil
}
