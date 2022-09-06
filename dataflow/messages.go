package dataflow

import (
	"context"
	"errors"
	"time"

	dataflow "google.golang.org/api/dataflow/v1b3"
)

type Message struct {
	Text  string
	Time  time.Time
	level string
}

func ListMessages(projectId string, location string, jobId string, onlyErrors bool) ([]Message, error) {
	ctx := context.Background()

	// Configure Service
	dataflowService, err := dataflow.NewService(ctx)
	if err != nil {
		return nil, err
	}

	// Create request
	messagesService := dataflow.NewProjectsLocationsJobsMessagesService(dataflowService)
	listRequest := messagesService.List(projectId, location, jobId)

	if onlyErrors {
		listRequest = listRequest.MinimumImportance("JOB_MESSAGE_ERROR")
	}

	// Cycle through pages
	var messages []Message
	err = listRequest.Pages(ctx, func(res *dataflow.ListJobMessagesResponse) error {
		for _, message := range res.JobMessages {

			t, err := time.Parse(time.RFC3339, message.Time)
			if err != nil {
				return errors.New("couldn't parse time")
			}

			messages = append(messages, Message{
				Text:  message.MessageText,
				Time:  t,
				level: message.MessageImportance,
			})
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return messages, nil
}
