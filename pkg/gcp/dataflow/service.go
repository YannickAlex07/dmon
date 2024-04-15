package dataflow

import (
	"context"
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/yannickalex07/inframon/pkg/util"
	dataflow "google.golang.org/api/dataflow/v1b3"
	"google.golang.org/api/option"
)

// Service Facade

// A service interface that can be used to interact with the Dataflow API
type DataflowService interface {
	// List all jobs from the Dataflow API
	ListJobs(ctx context.Context) ([]Job, error)

	// Get the logs associated with a specific job id
	GetLogs(ctx context.Context, jobId string, minLevel MessageLevel) ([]LogMessage, error)
}

// Servcie Implementation

// An implementation of the service to interact with Dataflow
type dataflowService struct {
	project  string
	location string

	service *dataflow.Service
}

// Creates a new DataflowService to interact with dataflow in the specified project and location.
// Use the ClientOptions to override the Dataflow API endpoint for testing.
func NewDataflowService(ctx context.Context, project string, location string, options []option.ClientOption) (DataflowService, error) {
	// create dataflow service
	service, err := dataflow.NewService(ctx, options...)
	if err != nil {
		log.Errorf("Failed to create the Dataflow service: %v", err)
		return nil, err
	}

	return &dataflowService{
		project:  project,
		location: location,
		service:  service,
	}, nil
}

// Get all jobs from the Dataflow API that can be found in the project and location
func (s *dataflowService) ListJobs(ctx context.Context) ([]Job, error) {
	log.Debug("Listing jobs...")

	// create list request
	jobService := dataflow.NewProjectsLocationsJobsService(s.service)
	req := jobService.List(s.project, s.location)

	// loop through pages
	jobs := []Job{}
	err := req.Pages(ctx, func(res *dataflow.ListJobsResponse) error {
		for _, j := range res.Jobs {
			contextLogger := log.WithField("job_id", j.Id)
			contextLogger.Debug("Parsing job...")

			// parse start time
			startTime, err := util.ParseTimestamp(j.StartTime)
			if err != nil {
				contextLogger.Errorf("Failed to parse start time with: %v", err)
				return fmt.Errorf("failed to parse start time with: %w", err)
			}

			// parse updated time
			statusTime, err := util.ParseTimestamp(j.CurrentStateTime)
			if err != nil {
				contextLogger.Errorf("Failed to parse status update time with: %v", err)
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

	log.Debugf("Found and parsed %d jobs", len(jobs))

	return jobs, nil
}

// Get the logs associated with a specific job id and a certain minimum level
func (s *dataflowService) GetLogs(ctx context.Context, jobId string, minLevel MessageLevel) ([]LogMessage, error) {
	contextLogger := log.WithFields(log.Fields{
		"job_id":    jobId,
		"min_level": minLevel,
	})
	contextLogger.Debug("Getting logs...")

	jobService := dataflow.NewProjectsLocationsJobsMessagesService(s.service)
	req := jobService.List(s.project, s.location, jobId)
	req.MinimumImportance(string(minLevel))

	entries := []LogMessage{}
	err := req.Pages(ctx, func(res *dataflow.ListJobMessagesResponse) error {
		for _, message := range res.JobMessages {
			contextLogger := contextLogger.WithField("message_id", message.Id)
			contextLogger.Debug("Parsing message...")

			// parse timestamps
			t, err := util.ParseTimestamp(message.Time)
			if err != nil {
				log.Errorf("Failed to parse message time with: %v", err)
				return fmt.Errorf("failed to parse message time with: %w", err)
			}

			for _, m := range strings.Split(message.MessageText, "\n") {
				// add entry
				e := LogMessage{
					Text:  m,
					Level: MessageLevelFromString(string(message.MessageImportance)),
					Time:  t,
				}

				entries = append(entries, e)
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	contextLogger.Debugf("Found and parsed %d log entries", len(entries))

	return entries, nil
}
