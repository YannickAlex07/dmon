package monitor

import (
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/yannickalex07/dmon/interfaces"
	"github.com/yannickalex07/dmon/models"
)

func Monitor(cfg models.Config, api interfaces.API, handlers []interfaces.Handler, stateStore interfaces.Storage) {
	log.Info("Starting new run.")

	// Dataflow API request
	log.Info("Requesting jobs from Dataflow API")

	jobs, err := api.Jobs(cfg.Project.Id, cfg.Project.Location)
	if err != nil {
		log.Errorf("Failed to list jobs with error %s", err.Error())
		return
	}

	log.Debugf("Found %d jobs", len(jobs))

	// checking job status
	lastRunTime := stateStore.GetLatestRunTime()

	for _, job := range jobs {
		// job was updated after last run
		if job.Status.UpdatedAt.After(lastRunTime) {
			log.WithFields(log.Fields{
				"id":        job.Id,
				"name":      job.Name,
				"status":    job.Status.Status,
				"updatedAt": job.Status.UpdatedAt,
			}).Info("Found Job with newer status")

			// handeling failed job
			if job.Status.IsFailed() {
				log.Info("Job %s has new failed status", job.Id)

				// requesting error messages from Dataflow
				log.Infof("Requesting error messages for job %s", job.Id)

				messages, err := api.Messages(cfg.Project.Id, cfg.Project.Location, job.Id)
				if err != nil {
					log.Errorf("Failed to query error messages for job %s with error %s", job.Id, err.Error())
					return
				}

				log.Debugf("Found %d error messages for job %s", len(messages), job.Id)

				// notifying handlers
				log.Infof("Notifying handlers for failed job %s", job.Id)

				for _, handler := range handlers {
					handler.HandleError(cfg, job, messages)
				}

				log.Debugf("Notified handlers for job %s", job.Id)
			}

			// handeling job timeout
			log.Info("Checking for running batch jobs.")

			if job.Status.IsRunning() && !job.IsStreaming() {
				log.Debugf("Batch job %s is currently running", job.Id)

				totalRunTime := time.Since(job.StartTime)

				// check if time runs longer than allowed
				log.Debugf("Checking if job %s has timeouted", job.Id)

				if totalRunTime > cfg.MaxTimeoutDuration() {
					log.Infof("Job %s crossed max allowed timeout duration", job.Id)

					// check if notification for job was already send
					wasNotified := stateStore.WasTimeoutHandled(job.Id)
					if !wasNotified {
						log.Infof("Timeout for job %s was not yet handled - handeling it now", job.Id)

						for _, handler := range handlers {
							handler.HandleTimeout(cfg, job)
						}

						stateStore.TimeoutHandled(job.Id)
						log.Infof("Timeout of job %s was handled", job.Id)
					}
				}
			}
		}
	}

	stateStore.SetLatestRunTime(time.Now().UTC())
	log.Info("Run finished.")
}
