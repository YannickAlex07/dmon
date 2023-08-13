package monitor

import (
	"errors"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/yannickalex07/dmon/pkg/dataflow"
	"github.com/yannickalex07/dmon/pkg/handler"
	"github.com/yannickalex07/dmon/pkg/model"
	"github.com/yannickalex07/dmon/pkg/storage"
)

type MonitorConfig struct {
	MaxJobTimeout time.Duration
}

func Monitor(cfg MonitorConfig, client dataflow.Dataflow, handlers []handler.Handler, stateStore storage.Storage) error {
	log.Info("Starting new run.")

	// Dataflow API request
	log.Info("Requesting jobs from Dataflow API")

	jobs, err := client.Jobs()
	if err != nil {
		errMsg := fmt.Sprintf("Failed to list jobs with error %s", err.Error())
		log.Errorf(errMsg)

		return errors.New(errMsg)
	}

	log.Debugf("Found %d jobs", len(jobs))

	// checking job status
	lastExecutionTime, err := stateStore.GetLatestExecutionTime()
	if err != nil {
		log.Warningf("Failed to fetch latest execution time from state store! Using now().")
		lastExecutionTime = time.Now().UTC()
	}

	for _, job := range jobs {
		// job was updated after last run
		if job.Status.UpdatedAt.After(lastExecutionTime) {
			log.WithFields(log.Fields{
				"id":        job.Id,
				"name":      job.Name,
				"status":    job.Status.Status,
				"updatedAt": job.Status.UpdatedAt,
			}).Info("Found Job with newer status")

			// handeling failed job
			if job.Status.IsFailed() {
				log.Infof("Job %s has new failed status", job.Id)

				// requesting error messages from Dataflow
				log.Infof("Requesting error log entries for job %s", job.Id)

				entries, err := client.ErrorLogs(job.Id)
				if err != nil {
					errMsg := fmt.Sprintf("Failed to query error entries for job %s with error %s", job.Id, err.Error())
					log.Errorf(errMsg)

					// we don't interrupt the application here and just pass 0 entries.
					entries = make([]model.LogEntry, 0)
				}

				log.Debugf("Found %d error entries for job %s", len(entries), job.Id)

				// notifying handlers
				log.Infof("Notifying handlers for failed job %s", job.Id)

				for _, handler := range handlers {
					handler.HandleError(job, entries)
				}

				log.Debugf("Notified handlers for job %s", job.Id)
			}
		}

		if job.Status.IsRunning() && !job.IsStreaming() {

			log.Debugf("Found running batch job %s", job.Id)
			totalRunTime := job.Runtime()

			// check if time runs longer than allowed
			log.Debugf("Checking if job %s has timeouted", job.Id)

			if totalRunTime > cfg.MaxJobTimeout {

				log.Infof("Job %s crossed max allowed timeout duration with a total runtime of %s", job.Id, totalRunTime.Round(time.Second))

				// check if notification for job was already send
				wasNotified := stateStore.WasTimeoutHandled(job.Id)
				if !wasNotified {

					log.Infof("Timeout for job %s was not yet handled - handeling it now", job.Id)

					for _, handler := range handlers {
						handler.HandleTimeout(job)
					}

					stateStore.HandleTimeout(job.Id, time.Now().UTC())
					log.Infof("Timeout of job %s was handled", job.Id)
				}
			}
		}
	}

	stateStore.SetLatestExecutionTime(time.Now().UTC())
	log.Info("Run finished.")

	return nil
}
