package dataflow

import (
	"context"
	"fmt"
	"strings"

	"github.com/yannickalex07/dmon/pkg/model"
	"github.com/yannickalex07/dmon/pkg/util"
	dataflow "google.golang.org/api/dataflow/v1b3"
)

func (client DataflowClient) Jobs(ctx context.Context) ([]model.Job, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	// create service and request
	service, err := dataflow.NewService(ctx)
	if err != nil {
		return nil, err
	}

	jobService := dataflow.NewProjectsLocationsJobsService(service)
	req := jobService.List(client.Project, client.Location)

	// request list of jobs
	var jobs []model.Job
	err = req.Pages(ctx, func(res *dataflow.ListJobsResponse) error {
		for _, job := range res.Jobs {

			// check if the name matches a prefix as long as one is required
			if client.Prefix != "" {
				matches := strings.HasPrefix(job.Name, client.Prefix)
				if !matches {
					continue
				}
			}

			// parse timestamps
			startTime, err := util.ParseTimestamp(job.StartTime)
			if err != nil {
				return fmt.Errorf("failed to parse state time with: %w", err)
			}

			statusTime, err := util.ParseTimestamp(job.CurrentStateTime)
			if err != nil {
				return fmt.Errorf("failed to parse current status time with %w", err)
			}

			// add job
			j := model.Job{
				Id:   job.Id,
				Name: job.Name,
				Type: job.Type,
				Status: model.Status{
					Status:    job.CurrentState,
					UpdatedAt: statusTime,
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

	return jobs, nil
}
