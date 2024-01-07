package dataflow

import (
	"context"
	"fmt"
	"log"

	"github.com/yannickalex07/dmon/pkg/util"
	dataflow "google.golang.org/api/dataflow/v1b3"
	"google.golang.org/api/option"
)

// Service Facade

type DataflowService interface {
	ListJobs(ctx context.Context) ([]Job, error)
	GetLogs(ctx context.Context, jobId string, minLevel MessageLevel) ([]LogMessage, error)
}

// Servcie Implementation

type dataflowService struct {
	project  string
	location string

	service *dataflow.Service
}

func NewDataflowService(ctx context.Context, project string, location string, options []option.ClientOption) DataflowService {
	// create dataflow service
	service, err := dataflow.NewService(ctx, options...)
	if err != nil {
		log.Fatalf("unable to create dataflow service: %v", err)
	}

	return &dataflowService{
		project:  project,
		location: location,
		service:  service,
	}
}

func (s *dataflowService) ListJobs(ctx context.Context) ([]Job, error) {
	// create list request
	jobService := dataflow.NewProjectsLocationsJobsService(s.service)
	req := jobService.List(s.project, s.location)

	// loop through pages
	jobs := []Job{}
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
				return fmt.Errorf("failed to parse status time with: %w", err)
			}

			// create dataflow job
			job := Job{
				Id:        j.Id,
				Name:      j.Name,
				Project:   s.project,
				Location:  s.location,
				Type:      j.Type,
				StartTime: startTime,
				Status: JobStatus{
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

func (s *dataflowService) GetLogs(ctx context.Context, jobId string, minLevel MessageLevel) ([]LogMessage, error) {
	jobService := dataflow.NewProjectsLocationsJobsMessagesService(s.service)
	req := jobService.List(s.project, s.location, jobId)
	req.MinimumImportance(string(minLevel))

	entries := []LogMessage{}
	err := req.Pages(ctx, func(res *dataflow.ListJobMessagesResponse) error {
		for _, message := range res.JobMessages {
			// parse timestamps
			t, err := util.ParseTimestamp(message.Time)
			if err != nil {
				return fmt.Errorf("failed to parse message time with: %w", err)
			}

			// add entry
			e := LogMessage{
				Text:  message.MessageText,
				Level: MessageLevelFromString(string(message.MessageImportance)),
				Time:  t,
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
