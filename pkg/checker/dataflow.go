package checker

import (
	"context"
	"net/url"
	"time"

	dataflow "github.com/yannickalex07/dmon/internal/gcp/dataflow"
	siren "github.com/yannickalex07/dmon/pkg"
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

func (c DataflowChecker) Check(ctx context.Context, since time.Time) ([]siren.Notification, error) {
	// list all jobs
	jobs, err := c.Service.ListJobs(ctx)
	if err != nil {
		return nil, err
	}

	notifications := []siren.Notification{}
	for _, job := range jobs {
		// filter down jobs by the provided filter
		if !c.JobFilter(job) {
			continue
		}

		// check all updated jobs
		if job.Status.UpdatedAt.After(since) {
			// check if the job failed
			if job.Status.IsFailed() {
				// request error logs
				logs := []string{}

				l, err := c.Service.GetLogs(ctx, job.Id, dataflow.LEVEL_ERROR)
				if err != nil {
					// log error event
					logs = append(logs, "Failed to fetch logs...")
				} else {
					for _, m := range l {
						logs = append(logs, m.Text)
					}
				}

				// create the notification
				n := siren.Notification{
					Title:       "❌ Dataflow Job Failed",
					Description: "",
					Logs:        logs,
					Links:       c.links(job),
				}

				notifications = append(notifications, n)
			}
		}

		// check runtime of running batch jobs
		if !job.IsStreaming() && job.Status.IsRunning() {
			if job.Runtime() >= c.Timeout {
				n := siren.Notification{
					Title:       "⏱️ Dataflow Job Running For Too Long",
					Description: "",
					Logs:        []string{},
					Links:       c.links(job),
				}

				notifications = append(notifications, n)
			}
		}
	}

	return notifications, nil
}

func (c DataflowChecker) links(job dataflow.Job) map[string]*url.URL {
	links := map[string]*url.URL{}

	// the url to the Dataflow UI
	u, err := url.Parse("")
	if err == nil {
		links["Open In Dataflow"] = u
	}

	return links
}
