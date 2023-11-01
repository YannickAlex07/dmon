package gcp

import (
	"context"

	gmon "github.com/yannickalex07/dmon/pkg"
	dataflow "google.golang.org/api/dataflow/v1b3"
	"google.golang.org/api/option"
)

// Models

type dataflowJob struct {
	Id   string
	Name string
	Type string
}

// Checker

type DataflowChecker struct {
	Project  string
	Location string

	// Additional options for the Dataflow service
	// Can be used to override the endpoint of the API
	ServiceOptions []option.ClientOption
}

func (c DataflowChecker) Check(ctx context.Context, handlers []gmon.Handler) error {
	_, err := c.listJobs(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (c DataflowChecker) listJobs(ctx context.Context) ([]dataflowJob, error) {
	// create dataflow service
	service, err := dataflow.NewService(ctx, c.ServiceOptions...)
	if err != nil {
		return nil, err
	}

	// create list request
	jobService := dataflow.NewProjectsLocationsJobsService(service)
	req := jobService.List(c.Project, c.Location)

	// loop through pages
	jobs := []dataflowJob{}
	err = req.Pages(ctx, func(res *dataflow.ListJobsResponse) error {
		for _, j := range res.Jobs {

			// create dataflow job
			job := dataflowJob{
				Id:   j.Id,
				Name: j.Name,
				Type: j.Type,
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
