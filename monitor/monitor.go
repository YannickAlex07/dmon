package monitor

import (
	"log"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/yannickalex07/dmon/interfaces"
	"github.com/yannickalex07/dmon/models"
)

func Start(cfg models.Config, api interfaces.API, handlers []interfaces.Handler) {

	notifiedTimeouts := make(map[string]time.Time)
	lastRun := time.Now()

	// monitoring function
	f := func() {
		// log that a request happens

		// clear our notified timeouts that expired
		for key, t := range notifiedTimeouts {
			tDuration := time.Since(t)
			if tDuration > cfg.ExpireTimeoutDuration() {
				delete(notifiedTimeouts, key)
			}
		}

		// get jobs
		jobs, err := api.Jobs(cfg.Project.Id, cfg.Project.Location)
		if err != nil {
			// log message
		}

		// update last runtime
		prevRunTime := lastRun
		lastRun = time.Now()

		// check jobs
		for _, job := range jobs {

			// check for status updates
			if job.Status.UpdatedAt.Before(prevRunTime) {

				// handle failure state
				if job.Status.IsFailed() {
					// fetch error messages
					messages, err := api.Messages(cfg.Project.Id, cfg.Project.Location, job.Id)
					if err != nil {
						// log error
					}

					for _, handler := range handlers {
						handler.HandleError(cfg, job, messages)
					}
				}
			}

			// check for jobs with timeout
			if job.Status.IsRunning() && !job.IsStreaming() {
				totalRunTime := time.Since(job.StartTime)
				_, ok := notifiedTimeouts[job.Id]

				// notify handlers of a timeout
				if totalRunTime > cfg.MaxTimeoutDuration() && !ok {
					for _, handler := range handlers {
						handler.HandleTimeout(cfg, job)
					}
					notifiedTimeouts[job.Id] = time.Now()
				}
			}
		}
	}

	// start scheduler
	log.Println("Start monitoring...")
	scheduler := gocron.NewScheduler(time.UTC)
	scheduler.Every(cfg.RequestInterval).Minute().Do(f)
	scheduler.StartBlocking()
}
