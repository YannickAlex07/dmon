package dataflow

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/yannickalex07/dmon/pkg/util"
	dataflow "google.golang.org/api/dataflow/v1b3"
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

// Check if the job is a streaming job.
func (j DataflowJob) IsBatch() bool {
	return j.Type == "JOB_TYPE_BATCH"
}

// Check the current runtime of the job. This is calculating by taking the time
// since the start time provided by the Dataflow backend.
func (j DataflowJob) Runtime() time.Duration {
	return time.Since(j.StartTime)
}

// Logs

// A Google log entry
type LogEntry struct {
	Text string
	Time time.Time
}

// Service Facade

type DataflowService interface {
	ListJobs(ctx context.Context) ([]DataflowJob, error)
	GetLogs(ctx context.Context, jobId string) ([]LogEntry, error)
}

// Servcie Implementation

type dataflowService struct {
	project  string
	location string

	service *dataflow.Service
}

func NewDataflowService(ctx context.Context, project string, location string) DataflowService {
	// create dataflow service
	service, err := dataflow.NewService(ctx)
	if err != nil {
		log.Fatalf("unable to create translate service, shutting down: %v", err)
	}

	return &dataflowService{
		project:  project,
		location: location,
		service:  service,
	}
}

func (s *dataflowService) ListJobs(ctx context.Context) ([]DataflowJob, error) {
	// create list request
	jobService := dataflow.NewProjectsLocationsJobsService(s.service)
	req := jobService.List(s.project, s.location)

	// loop through pages
	jobs := []DataflowJob{}
	err := req.Pages(ctx, func(res *dataflow.ListJobsResponse) error {
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

func (s *dataflowService) GetLogs(ctx context.Context, jobId string) ([]LogEntry, error) {
	jobService := dataflow.NewProjectsLocationsJobsMessagesService(s.service)
	req := jobService.List(s.project, s.location, jobId)

	entries := []LogEntry{}
	err := req.Pages(ctx, func(res *dataflow.ListJobMessagesResponse) error {
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
			e := LogEntry{
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
