package dataflow

import (
	"context"
	"fmt"
	"net/url"
	"time"

	log "github.com/sirupsen/logrus"
	inframon "github.com/yannickalex07/inframon/pkg"
)

type notificationType string

const (
	errNotification     = "ERROR"
	timeoutNotification = "TIMEOUT"
)

// Source

// A source to get notifications from Dataflow.
// Will check for failed jobs as well as batch jobs that run for too long.
type DataflowSource struct {
	Service DataflowService

	// A custom filter that can be used to filter out specific jobs to check.
	// Don't use that field to filter for failed jobs or jobs that run for too long,
	// this will already be done by the Checker itself.
	JobFilter func(Job) bool

	// Configure when a job is marked as timed out.
	Timeout time.Duration
}

func (c DataflowSource) Check(ctx context.Context, since time.Time) ([]inframon.Notification, error) {
	// list all jobs
	jobs, err := c.Service.ListJobs(ctx)
	if err != nil {
		return nil, err
	}

	log.Debugf("Found %d jobs", len(jobs))
	notifications := []inframon.Notification{}
	for _, job := range jobs {
		contextLogger := log.WithFields(log.Fields{
			"job_id":     job.Id,
			"job_status": job.Status.Status,
		})

		// filter down jobs by the provided filter
		if c.JobFilter != nil && !c.JobFilter(job) {
			contextLogger.Debug("Job skipped because it doesn't match the filter")
			continue
		}

		// check all updated jobs
		if job.Status.UpdatedAt.After(since) {
			contextLogger.Debug("Checking udpated job")

			// check if the job failed
			if job.Status.IsFailed() {
				contextLogger.Debug("Job has failed since last check")

				// request error logs
				logs := []string{}

				contextLogger.Debug("Fetching logs")
				l, err := c.Service.GetLogs(ctx, job.Id, LEVEL_ERROR)
				if err != nil {
					// log error event
					contextLogger.WithField("err", err).Error("Failed to fetch logs")
					logs = append(logs, "Failed to fetch logs...")
				} else {
					contextLogger.Debug("Fetched logs")
					for _, m := range l {
						logs = append(logs, m.Text)
					}
				}

				// create the notification
				contextLogger.Debug("Creating notification")
				n := inframon.Notification{
					Key:         c.createNotificationKey(errNotification, job.Id, job.StartTime),
					Title:       "❌ Dataflow Job Failed",
					Description: fmt.Sprintf("The job `%s` with id `%s` failed at *%s*!", job.Name, job.Id, job.Status.UpdatedAt.Format(time.RFC1123)),
					Logs:        logs,
					Links:       c.links(job),
				}

				contextLogger.WithFields(log.Fields{
					"title":       n.Title,
					"description": n.Description,
					"key":         n.Key,
				}).Debug("Created notification")

				notifications = append(notifications, n)
			}
		}

		// check runtime of running batch jobs
		if !job.IsStreaming() && job.Status.IsRunning() {
			runtime := job.Runtime()
			contextLogger.WithField("runtime", runtime).Debug("Checking runtime")

			if job.Runtime() >= c.Timeout {
				contextLogger.Debug("Job crossed timeout limit")
				n := inframon.Notification{
					Key:         c.createNotificationKey(timeoutNotification, job.Id, job.StartTime),
					Title:       "⏱️ Dataflow Job Running For Too Long",
					Description: fmt.Sprintf("The job `%s` with id `%s` crossed the maximum timeout limit with a runtime of *%s*.", job.Name, job.Id, job.Runtime().Round(time.Second)),
					Logs:        []string{},
					Links:       c.links(job),
				}

				contextLogger.WithFields(log.Fields{
					"title":       n.Title,
					"description": n.Description,
					"key":         n.Key,
				}).Debug("Created notification")

				notifications = append(notifications, n)
			}
		}
	}

	return notifications, nil
}

func (c *DataflowSource) links(job Job) map[string]*url.URL {
	links := map[string]*url.URL{}

	// the url to the Dataflow UI
	u, err := url.Parse(fmt.Sprintf("https://console.cloud.google.com/dataflow/jobs/%s/%s?project=%s&authuser=1&hl=en", job.Location, job.Id, job.Project))
	if err == nil {
		links["Open In Dataflow"] = u
	}

	return links
}

func (c *DataflowSource) createNotificationKey(nType notificationType, jobId string, startTime time.Time) string {
	return fmt.Sprintf("DATAFLOW-%s-%s-%s", nType, jobId, startTime.Format(time.RFC3339))
}
