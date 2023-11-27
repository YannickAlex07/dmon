package checker

import (
	"context"
	"fmt"
	"net/url"
	"time"

	siren "github.com/yannickalex07/dmon/pkg"
	"github.com/yannickalex07/dmon/pkg/util"
	dataflow "google.golang.org/api/dataflow/v1b3"
	"google.golang.org/api/option"
)

// Job Status

// The Status of a Dataflow Job
type DataflowJobStatus struct {
	// The raw status string of the job. You can check this manually
	// or use the provided `IsXXX()`-methods on this struct.
	Status string

	// The time the status was last updated by the Dataflow backend.
	UpdatedAt time.Time
}

// Check if the job has failed according to its status.
func (js DataflowJobStatus) IsFailed() bool {
	return js.Status == "JOB_STATE_FAILED"
}

// Check if the job is running according to its status.
func (js DataflowJobStatus) IsRunning() bool {
	return js.Status == "JOB_STATE_RUNNING"
}

// Job

// A Dataflow Job
type DataflowJob struct {
	Id   string
	Name string

	// The raw type of the job. You can check this field manually or use
	// the provided `IsXXX()`-methods on this struct to check it.
	Type string

	// The time that the job started according to the Dataflow backend.
	StartTime time.Time

	// The current status of the Dataflow job, containing the state it is in
	// as well as when it was last updated.
	Status DataflowJobStatus
}

// Check if the job is a streaming job.
func (j DataflowJob) IsStreaming() bool {
	return j.Type == "JOB_TYPE_STREAMING"
}

// Check the current runtime of the job. This is calculating by taking the time
// since the start time provided by the Dataflow backend.
func (j DataflowJob) Runtime() time.Duration {
	return time.Since(j.StartTime)
}

// Logs

// A Google log entry
type logEntry struct {
	Text string
	Time time.Time
}

// Checker

// A checker for Dataflow.
// Will check for failed jobs as well as batch jobs that run for too long.
type DataflowChecker struct {
	Project  string
	Location string

	// A custom filter that can be used to filter out specific jobs to check.
	// Don't use that field to filter for failed jobs or jobs that run for too long,
	// this will already be done by the Checker itself.
	JobFilter func(DataflowJob) bool

	// Configure when a job is marked as timed out.
	Timeout time.Duration

	// Additional options for the Dataflow service
	// Can be used to override the endpoint of the API
	ServiceOptions []option.ClientOption
}

func (c DataflowChecker) Check(ctx context.Context, since time.Time) ([]siren.Notification, error) {
	// list all jobs
	jobs, err := c.listJobs(ctx)
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

				l, err := c.listErrorLogs(ctx, job.Id)
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
					Title:    "❌ Dataflow Job Failed",
					Overview: "",
					Logs:     logs,
					Links:    c.links(job),
				}

				notifications = append(notifications, n)
			}
		}

		// check runtime of running batch jobs
		if !job.IsStreaming() && job.Status.IsRunning() {
			if job.Runtime() >= c.Timeout {
				n := siren.Notification{
					Title:    "⏱️ Dataflow Job Running For Too Long",
					Overview: "",
					Logs:     []string{},
					Links:    c.links(job),
				}

				notifications = append(notifications, n)
			}
		}
	}

	return notifications, nil
}

func (c DataflowChecker) links(job DataflowJob) map[string]*url.URL {
	links := map[string]*url.URL{}

	// the url to the Dataflow UI
	u, err := url.Parse("")
	if err == nil {
		links["Open In Dataflow"] = u
	}

	return links
}

func (c DataflowChecker) listJobs(ctx context.Context) ([]DataflowJob, error) {
	// create dataflow service
	service, err := dataflow.NewService(ctx, c.ServiceOptions...)
	if err != nil {
		return nil, err
	}

	// create list request
	jobService := dataflow.NewProjectsLocationsJobsService(service)
	req := jobService.List(c.Project, c.Location)

	// loop through pages
	jobs := []DataflowJob{}
	err = req.Pages(ctx, func(res *dataflow.ListJobsResponse) error {
		for _, j := range res.Jobs {
			// parse start time
			startTime, err := util.ParseTimestamp(j.StartTime)
			if err != nil {
				return fmt.Errorf("failed to parse start time with: %w", err)
			}

			// parse updated time
			statusTime, err := util.ParseTimestamp(j.CurrentStateTime)
			if err != nil {
				return fmt.Errorf("failed to parse start time with: %w", err)
			}

			// create dataflow job
			job := DataflowJob{
				Id:        j.Id,
				Name:      j.Name,
				Type:      j.Type,
				StartTime: startTime,
				Status: DataflowJobStatus{
					Status:    j.CurrentState,
					UpdatedAt: statusTime,
				},
			}

			jobs = append(jobs, job)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return jobs, nil
}

func (c DataflowChecker) listErrorLogs(ctx context.Context, jobId string) ([]logEntry, error) {
	// create dataflow service
	service, err := dataflow.NewService(ctx, c.ServiceOptions...)
	if err != nil {
		return nil, err
	}

	jobService := dataflow.NewProjectsLocationsJobsMessagesService(service)
	req := jobService.List(c.Project, c.Location, jobId)

	entries := []logEntry{}
	err = req.Pages(ctx, func(res *dataflow.ListJobMessagesResponse) error {
		for _, message := range res.JobMessages {
			// skip any entry that is not an error
			if message.MessageImportance != "JOB_MESSAGE_ERROR" {
				continue
			}

			// parse timestamps
			t, err := util.ParseTimestamp(message.Time)
			if err != nil {
				return fmt.Errorf("failed to parse entry time with: %w", err)
			}

			// add entry
			e := logEntry{
				Text: message.MessageText,
				Time: t,
			}

			entries = append(entries, e)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return entries, nil
}
