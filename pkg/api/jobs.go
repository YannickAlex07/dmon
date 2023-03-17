package api

import (
	"context"
	"errors"

	"github.com/yannickalex07/dmon/pkg/models"
	"github.com/yannickalex07/dmon/pkg/util"
	dataflow "google.golang.org/api/dataflow/v1b3"
)

func (api API) Jobs(project string, location string) ([]models.Job, error) {
	ctx := context.Background()

	// create service and request
	service, err := dataflow.NewService(ctx)
	if err != nil {
		return nil, err
	}

	jobService := dataflow.NewProjectsLocationsJobsService(service)
	req := jobService.List(project, location)

	// request list of jobs
	var jobs []models.Job
	err = req.Pages(ctx, func(res *dataflow.ListJobsResponse) error {
		for _, job := range res.Jobs {

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
