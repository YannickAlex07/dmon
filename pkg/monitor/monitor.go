package monitor

import (
	"context"
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

func Monitor(ctx context.Context, cfg MonitorConfig, client dataflow.Dataflow, handlers []handler.Handler, stateStore storage.Storage) error {
	log.Info("Starting new run.")

	// checking job status
	lastExecutionTime, err := stateStore.GetLatestExecutionTime(ctx)
	if err != nil {
		log.Warningf("Failed to fetch latest execution time from state store! Using now().")
		lastExecutionTime = time.Now().UTC()
	}

	// Dataflow API request
	log.Info("Requesting jobs from Dataflow API")

	jobs, err := client.Jobs(ctx)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to list jobs with error %w", err)
		log.Errorf(wrappedErr.Error())

		return wrappedErr
	}

	log.Debugf("Found %d jobs", len(jobs))

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

				entries, err := client.ErrorLogs(ctx, job.Id)
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
					err := handler.HandleError(ctx, job, entries)
					if err != nil {
						log.Errorf("handler failed to handle job error: %s", err.Error())
					}
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
				isStored, err := stateStore.IsTimeoutStored(ctx, job.Id)
				if err != nil {
					log.Errorf("failed to fetch if timeout is stored: %s", err.Error())
					isStored = false
				}

				if !isStored {
					log.Infof("Timeout for job %s was not yet handled - handeling it now", job.Id)

					for _, handler := range handlers {
						err := handler.HandleTimeout(ctx, job)
						if err != nil {
							log.Errorf("handler failed to handle job timeout: %s", err.Error())
						}
					}

					err := stateStore.StoreTimeout(ctx, job.Id, time.Now().UTC())
					if err != nil {
						log.Errorf("failed to store timeout with err: %s", err.Error())
					}

					log.Infof("Timeout of job %s was handled", job.Id)
				}
			}
		}
	}

	err = stateStore.SetLatestExecutionTime(ctx, time.Now().UTC())
	if err != nil {
		log.Errorf("failed to set latest execution time: %s", err.Error())
	}

	log.Info("Run finished.")

	return nil
}
