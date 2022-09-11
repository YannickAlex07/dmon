package api

import (
	"context"
	"errors"

	"github.com/yannickalex07/dmon/models"
	"github.com/yannickalex07/dmon/util"
	dataflow "google.golang.org/api/dataflow/v1b3"
)

func (api API) Messages(project string, location string, jobId string) ([]models.Message, error) {
	ctx := context.Background()

	// create service and request
	service, err := dataflow.NewService(ctx)
	if err != nil {
		return nil, err
	}

	jobService := dataflow.NewProjectsLocationsJobsMessagesService(service)
	req := jobService.List(project, location, jobId)

	// request list of jobs
	var messages []models.Message
	err = req.Pages(ctx, func(res *dataflow.ListJobMessagesResponse) error {
		for _, message := range res.JobMessages {

			// parse timestamps
			t, err := util.ParseTimestamp(message.Time)
			if err != nil {
				return errors.New("failed to parse message time")
			}

			// add message
			m := models.Message{
				Text: message.MessageText,
				Time: t,
			}

			messages = append(messages, m)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return messages, nil
}
