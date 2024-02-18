package checker

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"time"

	keiho "github.com/yannickalex07/dmon/pkg"
	dataflow "github.com/yannickalex07/dmon/pkg/external/gcp/dataflow"
)

type notificationType string

const (
	errNotification     = "ERROR"
	timeoutNotification = "TIMEOUT"
)

// Checker

// A checker for Dataflow.
// Will check for failed jobs as well as batch jobs that run for too long.
type DataflowChecker struct {
	Service dataflow.DataflowService

	// A custom filter that can be used to filter out specific jobs to check.
	// Don't use that field to filter for failed jobs or jobs that run for too long,
	// this will already be done by the Checker itself.
	JobFilter func(dataflow.Job) bool

	// Configure when a job is marked as timed out.
	Timeout time.Duration
}

func (c DataflowChecker) Check(ctx context.Context, since time.Time) ([]keiho.Notification, error) {
	// list all jobs
	log.Println("listing jobs")
	jobs, err := c.Service.ListJobs(ctx)
	if err != nil {
		return nil, err
	}

	log.Printf("found %d jobs", len(jobs))
	notifications := []keiho.Notification{}
	for _, job := range jobs {
		// filter down jobs by the provided filter
		if c.JobFilter != nil && !c.JobFilter(job) {
			continue
		}

		// check all updated jobs
		if job.Status.UpdatedAt.After(since) {
			log.Printf("checking udpated job: %s", job.Id)
			// check if the job failed
			if job.Status.IsFailed() {
				log.Printf("job failed: %s", job.Id)

				// request error logs
				logs := []string{}

				log.Println("fetching logs")
				l, err := c.Service.GetLogs(ctx, job.Id, dataflow.LEVEL_ERROR)
				if err != nil {
					// log error event
					logs = append(logs, "Failed to fetch logs...")
					log.Printf("failed to fetch logs: %v", err)
				} else {
					log.Println("fetched logs")
					for _, m := range l {
						logs = append(logs, m.Text)
					}
				}

				// create the notification
				log.Println("creating notification")
				n := keiho.Notification{
					Key:         c.createNotificationKey(errNotification, job.Id, job.StartTime),
					Title:       "❌ Dataflow Job Failed",
					Description: fmt.Sprintf("The job `%s` with id `%s` failed at *%s*!", job.Name, job.Id, job.Status.UpdatedAt.Format(time.RFC1123)),
					Logs:        logs,
					Links:       c.links(job),
				}

				log.Printf("created notification: Title(%s) && Description(%s)", n.Title, n.Description)

				notifications = append(notifications, n)
			}
		}

		// check runtime of running batch jobs
		if !job.IsStreaming() && job.Status.IsRunning() {
			log.Printf("checking runtime of job: %s", job.Id)
			if job.Runtime() >= c.Timeout {
				log.Printf("job is running for too long: %s", job.Id)
				n := keiho.Notification{
					Key:         c.createNotificationKey(timeoutNotification, job.Id, job.StartTime),
					Title:       "⏱️ Dataflow Job Running For Too Long",
					Description: fmt.Sprintf("The job `%s` with id `%s` crossed the maximum timeout limit with a runtime of *%s*.", job.Name, job.Id, job.Runtime().Round(time.Second)),
					Logs:        []string{},
					Links:       c.links(job),
				}

				log.Printf("created notification: %v", n)
				notifications = append(notifications, n)
			}
		}
	}

	return notifications, nil
}

func (c *DataflowChecker) links(job dataflow.Job) map[string]*url.URL {
	links := map[string]*url.URL{}

	// the url to the Dataflow UI
	u, err := url.Parse(fmt.Sprintf("https://console.cloud.google.com/dataflow/jobs/%s/%s?project=%s&authuser=1&hl=en", job.Location, job.Id, job.Project))
	if err == nil {
		links["Open In Dataflow"] = u
	}

	return links
}

func (c *DataflowChecker) createNotificationKey(nType notificationType, jobId string, startTime time.Time) string {
	return fmt.Sprintf("DATAFLOW-%s-%s-%s", nType, jobId, startTime.Format(time.RFC3339))
}
