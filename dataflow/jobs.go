package dataflow

import (
	"context"
	"errors"
	"time"

	dataflow "google.golang.org/api/dataflow/v1b3"
)

// Types

type Job struct {
	Id        string
	Name      string
	Status    Status
	StartTime time.Time
}

type Status struct {
	status     string
	updateTime time.Time
}

func (s Status) IsNewer(t time.Time) bool {
	return t.Before(s.updateTime)
}

func (s Status) IsFailed() bool {
	return s.status == "JOB_STATE_FAILED"
}

func (s Status) IsCanceled() bool {
	return s.status == "JOB_STATE_CANCELLED"
}

func (s Status) IsDone() bool {
	return s.status == "JOB_STATE_DONE"
}

func (s Status) IsRunning() bool {
	return s.status == "JOB_STATE_RUNNING"
}

// Functions

func ListJobs(projectId string, location string) ([]Job, error) {
	ctx := context.Background()

	// Configure Service
	dataflowService, err := dataflow.NewService(ctx)
	if err != nil {
		return nil, err
	}

	jobsService := dataflow.NewProjectsLocationsJobsService(dataflowService)
	listRequest := jobsService.List(projectId, location)

	// Request Pages
	var jobs []Job
	err = listRequest.Pages(ctx, func(res *dataflow.ListJobsResponse) error {
		for _, job := range res.Jobs {
			// Parse Timestamps
			updateTime, err := time.Parse(time.RFC3339, job.CurrentStateTime)
			if err != nil {
				return errors.New("couldn't parse time")
			}

			startTime, err := time.Parse(time.RFC3339, job.StartTime)
			if err != nil {
				return errors.New("couldn't parse time")
			}

			// Parse Job
			j := Job{
				Id:   job.Id,
				Name: job.Name,
				Status: Status{
					status:     job.CurrentState,
					updateTime: updateTime,
				},
				StartTime: startTime,
			}

			jobs = append(jobs, j)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Return Jobs
	return jobs, nil
}
