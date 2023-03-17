package monitor

import (
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/yannickalex07/dmon/pkg/interfaces"
	"github.com/yannickalex07/dmon/pkg/models"
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
	lastRunTime := stateStore.GetLatestRuntime()

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
				log.Infof("Requesting error log entries for job %s", job.Id)

				entries, err := api.ErrorLogs(cfg.Project.Id, cfg.Project.Location, job.Id)
				if err != nil {
					log.Errorf("Failed to query error entries for job %s with error %s", job.Id, err.Error())
					return
				}

				log.Debugf("Found %d error entries for job %s", len(entries), job.Id)

				// notifying handlers
				log.Infof("Notifying handlers for failed job %s", job.Id)

				for _, handler := range handlers {
					handler.HandleError(cfg, job, entries)
				}

				log.Debugf("Notified handlers for job %s", job.Id)
			}
		}

		if job.Status.IsRunning() && !job.IsStreaming() {

			log.Debugf("Found running batch job %s", job.Id)
			totalRunTime := job.Runtime()

			// check if time runs longer than allowed
			log.Debugf("Checking if job %s has timeouted", job.Id)

			if totalRunTime > cfg.MaxTimeoutDuration() {

				log.Infof("Job %s crossed max allowed timeout duration with a total runtime of %s", job.Id, totalRunTime.Round(time.Second))

				// check if notification for job was already send
				wasNotified := stateStore.TimeoutAlreadyHandled(job.Id)
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

	stateStore.SetLatestRuntime(time.Now().UTC())
	log.Info("Run finished.")
}
