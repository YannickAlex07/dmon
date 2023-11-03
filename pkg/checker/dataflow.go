package checker

import (
	"context"
	"fmt"
	"time"

	"github.com/yannickalex07/dmon/pkg/util"
	dataflow "google.golang.org/api/dataflow/v1b3"
	"google.golang.org/api/option"
)

// Job Status

type jobStatus struct {
	Status    string
	UpdatedAt time.Time
}

func (js jobStatus) IsFailed() bool {
	return js.Status == "JOB_STATE_FAILED"
}

func (js jobStatus) IsRunning() bool {
	return js.Status == "JOB_STATE_RUNNING"
}

// Job

type job struct {
	Id   string
	Name string
	Type string

	StartTime time.Time

	Status jobStatus
}

func (j job) IsStreaming() bool {
	return j.Type == "JOB_TYPE_STREAMING"
}

func (j job) Runtime() time.Duration {
	return time.Since(j.StartTime)
}

// Checker

type DataflowChecker struct {
	Project  string
	Location string

	// Additional options for the Dataflow service
	// Can be used to override the endpoint of the API
	ServiceOptions []option.ClientOption
}

func (c DataflowChecker) Check(ctx context.Context, since time.Time) error {
	// list all jobs
	_, err := c.listJobs(ctx)
	if err != nil {
		return err
	}

	// filter down jobs by prefix

	// filter down updated & failed jobs

	// filter down and check streaming jobs

	return nil
}

func (c DataflowChecker) listJobs(ctx context.Context) ([]job, error) {
	// create dataflow service
	service, err := dataflow.NewService(ctx, c.ServiceOptions...)
	if err != nil {
		return nil, err
	}

	// create list request
	jobService := dataflow.NewProjectsLocationsJobsService(service)
	req := jobService.List(c.Project, c.Location)

	// loop through pages
	jobs := []job{}
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
			job := job{
				Id:        j.Id,
				Name:      j.Name,
				Type:      j.Type,
				StartTime: startTime,
				Status: jobStatus{
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
