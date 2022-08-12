package dataflow

import (
	"context"
)

type JobsService interface {
	List(projectId string, location string) JobsRequest
}

type JobsRequest interface {
	Pages(context.Context, func(JobsResponse) error) error
}

type JobsResponse interface{}

func ListJobs(service JobsService, projectId string, location string) {
	ctx := context.Background()
	req := service.List(projectId, location)

	err := req.Pages(ctx, func(res JobsResponse) error {
		return nil
	})

	if err != nil {
		println("Failed!")
	}
}
