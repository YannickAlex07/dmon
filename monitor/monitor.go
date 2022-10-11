package monitor

import (
	"time"

	"github.com/go-co-op/gocron"
	log "github.com/sirupsen/logrus"
	"github.com/yannickalex07/dmon/interfaces"
	"github.com/yannickalex07/dmon/models"
)

func Start(cfg models.Config, api interfaces.API, handlers []interfaces.Handler) {
	notifiedTimeouts := make(map[string]time.Time)
	lastRun := time.Now()

	// monitoring function
	f := func() {
		log.Info("Starting new request...")

		// clear our notified timeouts that expired
		log.Debug("Clearing out timouts...")

		cleaned := 0
		for key, t := range notifiedTimeouts {
			tDuration := time.Since(t)
			if tDuration > cfg.ExpireTimeoutDuration() {
				delete(notifiedTimeouts, key)
				cleaned += 1
			}
		}

		log.Debugf("Cleared out %d timeouts.", cleaned)

		// get jobs
		log.Debug("Requesting jobs from Dataflow API...")

		jobs, err := api.Jobs(cfg.Project.Id, cfg.Project.Location)
		if err != nil {
			log.Errorf("Failed to query jobs with err %s.", err.Error())
			return
		}

		log.Debugf("Found %d jobs.", len(jobs))

		// update last runtime
		prevRunTime := lastRun
		lastRun = time.Now()

		// check jobs
		log.Debug("Iterating through jobs...")

		for _, job := range jobs {

			// check for status updates
			if job.Status.UpdatedAt.Before(prevRunTime) {

				log.WithFields(log.Fields{
					"id":     job.Id,
					"name":   job.Name,
					"status": job.Status.Status,
				}).Info("Found newer job status")

				// handle failure state
				if job.Status.IsFailed() {
					// fetch error messages
					log.Debugf("Requesting error messages for job %s...", job.Id)

					messages, err := api.Messages(cfg.Project.Id, cfg.Project.Location, job.Id)
					if err != nil {
						log.Errorf("Failed to query error messages with err %s", err.Error())
						return
					}

					log.Debugf("Found %d error messages.", len(messages))

					// handle errors
					log.Infof("Handling error for job %s...", job.Id)

					for _, handler := range handlers {
						handler.HandleError(cfg, job, messages)
					}

					log.Debug("Notified each handler.")
				}
			}

			// check for jobs with timeout
			log.Debug("Checking for timeouts...")

			if job.Status.IsRunning() && !job.IsStreaming() {
				log.Debugf("Found currently running job %s...", job.Id)

				totalRunTime := time.Since(job.StartTime)
				_, ok := notifiedTimeouts[job.Id]

				if ok {
					log.Debugf("Job %s already present in notified timeouts.", job.Id)
				}

				// notify handlers of a timeout
				if totalRunTime > cfg.MaxTimeoutDuration() && !ok {
					log.Infof("Job %s crossed max allowed timeout duration. Sending to handlers...", job.Id)

					for _, handler := range handlers {
						handler.HandleTimeout(cfg, job)
					}

					notifiedTimeouts[job.Id] = time.Now()

					log.Debug("Job send to handlers and stored in notified timeouts.")
				}
			}
		}
	}

	log.Info("Starting monitor...")

	// start scheduler
	scheduler := gocron.NewScheduler(time.UTC)
	scheduler.Every(cfg.RequestInterval).Minute().Do(f)
	scheduler.StartBlocking()
}
