package api

import (
	"context"
	"errors"

	"github.com/yannickalex07/dmon/pkg/models"
	"github.com/yannickalex07/dmon/pkg/util"
	dataflow "google.golang.org/api/dataflow/v1b3"
)

func (api API) ErrorLogs(project string, location string, jobId string) ([]models.LogEntry, error) {
	ctx := context.Background()

	// create service and request
	service, err := dataflow.NewService(ctx)
	if err != nil {
		return nil, err
	}

	jobService := dataflow.NewProjectsLocationsJobsMessagesService(service)
	req := jobService.List(project, location, jobId)

	// request list of jobs
	var entries []models.LogEntry
	err = req.Pages(ctx, func(res *dataflow.ListJobMessagesResponse) error {
		for _, message := range res.JobMessages {
			// skip any entry that is not an error
			if message.MessageImportance != "JOB_MESSAGE_ERROR" {
				continue
			}

			// parse timestamps
			t, err := util.ParseTimestamp(message.Time)
			if err != nil {
				return errors.New("failed to parse entry time")
			}

			// add entry
			e := models.LogEntry{
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
