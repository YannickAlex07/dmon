package dataflow

import (
	"context"
	"fmt"

	"github.com/yannickalex07/dmon/pkg/model"
	"github.com/yannickalex07/dmon/pkg/util"
	dataflow "google.golang.org/api/dataflow/v1b3"
)

func (client DataflowClient) ErrorLogs(ctx context.Context, jobId string) ([]model.LogEntry, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	// create service and request
	service, err := dataflow.NewService(ctx)
	if err != nil {
		return nil, err
	}

	jobService := dataflow.NewProjectsLocationsJobsMessagesService(service)
	req := jobService.List(client.Project, client.Location, jobId)

	// request list of jobs
	var entries []model.LogEntry
	err = req.Pages(ctx, func(res *dataflow.ListJobMessagesResponse) error {
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
			e := model.LogEntry{
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
