package client

import (
	"context"
	"errors"
	"strings"

	"github.com/yannickalex07/dmon/pkg/models"
	"github.com/yannickalex07/dmon/pkg/util"
	dataflow "google.golang.org/api/dataflow/v1b3"
)

func (client DataflowClient) Jobs() ([]models.Job, error) {
	ctx := context.Background()

	// create service and request
	service, err := dataflow.NewService(ctx)
	if err != nil {
		return nil, err
	}

	jobService := dataflow.NewProjectsLocationsJobsService(service)
	req := jobService.List(client.Project, client.Location)

	// request list of jobs
	var jobs []models.Job
	err = req.Pages(ctx, func(res *dataflow.ListJobsResponse) error {
		for _, job := range res.Jobs {

			// check if the name matches a prefix as long as one is required
			if client.Prefix != nil && *client.Prefix != "" {
				matches := strings.HasPrefix(job.Name, *client.Prefix)
				if !matches {
					continue
				}
			}

			// parse timestamps
			startTime, err := util.ParseTimestamp(job.StartTime)
			if err != nil {
				return errors.New("failed to parse start time")
			}

			statusTime, err := util.ParseTimestamp(job.CurrentStateTime)
			if err != nil {
				return errors.New("failed to parse current status time")
			}

			// add job
			j := models.Job{
				Id:   job.Id,
				Name: job.Name,
				Type: job.Type,
				Status: models.Status{
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
