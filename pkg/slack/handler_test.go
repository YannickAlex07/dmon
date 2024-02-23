package slack_test

import (
	"context"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	keiho "github.com/yannickalex07/dmon/pkg"
	"github.com/yannickalex07/dmon/pkg/slack"
)

func TestSlackHandlerWithoutLogs(t *testing.T) {
	// Arrange
	ctx := context.Background()
	replacer := strings.NewReplacer(
		"\t", "",
		"\n", "",
		" ", "",
	)

	service := SlackServiceMock{}
	handler := slack.SlackHandler{Service: &service, Channel: "channel"}

	notification := keiho.Notification{
		Key:         "Test",
		Title:       "Test",
		Description: "Test",
		Logs:        []string{},
		Links: map[string]*url.URL{
			"Test": {Path: "test.com"},
		},
	}

	expectedMessage := Message{
		channel: "channel",
		blocks: replacer.Replace(`
		[
			{
				"type": "header",
				"text": {
					"type": "plain_text",
					"text": "Test",
					"emoji": true
				}
			},
			{
				"type": "section",
				"text": {
					"type": "mrkdwn",
					"text": "Test"
				}
			},
			{
				"type":"actions",
				"block_id":"test",
				"elements":[
					{
						"type": "button",
						"text": {
							"type": "plain_text",
							"text": "Test"
						},
						"action_id": "test",
						"url": "test.com",
						"value": "test"
					}
				]
			}
		]
		`),
	}

	// Act
	err := handler.Handle(ctx, notification)
	if err != nil {
		t.Fatal(err)
	}

	// Assert
	assert.Equal(t, expectedMessage, service.Messages[0])
}

func TestSlackHandlerWithSingleLogMessage(t *testing.T) {
	// Arrange
	ctx := context.Background()
	replacer := strings.NewReplacer(
		"\t", "",
		"\n", "",
		" ", "",
		"LOGMESSAGE", "Last Log Messages: ```Test```",
	)

	service := SlackServiceMock{}
	handler := slack.SlackHandler{Service: &service, Channel: "channel"}

	notification := keiho.Notification{
		Key:         "Test",
		Title:       "Test",
		Description: "Test",
		Logs: []string{
			"Test",
		},
		Links: map[string]*url.URL{
			"Test": {Path: "test.com"},
		},
	}

	expectedMessage := Message{
		channel: "channel",
		blocks: replacer.Replace(`
		[
			{
				"type": "header",
				"text": {
					"type": "plain_text",
					"text": "Test",
					"emoji": true
				}
			},
			{
				"type": "section",
				"text": {
					"type": "mrkdwn",
					"text": "Test"
				}
			},
			{
				"type": "section",
				"text": {
					"type": "mrkdwn",
					"text": "LOGMESSAGE"
				}
			},
			{
				"type":"actions",
				"block_id":"test",
				"elements":[
					{
						"type": "button",
						"text": {
							"type": "plain_text",
							"text": "Test"
						},
						"action_id": "test",
						"url": "test.com",
						"value": "test"
					}
				]
			}
		]
		`),
	}

	// Act
	err := handler.Handle(ctx, notification)
	if err != nil {
		t.Fatal(err)
	}

	// Assert
	assert.Equal(t, expectedMessage, service.Messages[0])
}

func TestSlackHandlerWithMultipleLogMessages(t *testing.T) {
	// Arrange
	ctx := context.Background()
	replacer := strings.NewReplacer(
		"\t", "",
		"\n", "",
		" ", "",
		"LOGMESSAGE", "Last Log Messages: ```2\\n3\\n4\\n5\\n6```",
	)

	service := SlackServiceMock{}
	handler := slack.SlackHandler{Service: &service, Channel: "channel"}

	notification := keiho.Notification{
		Key:         "Test",
		Title:       "Test",
		Description: "Test",
		Logs: []string{
			"1",
			"2",
			"3",
			"4",
			"5",
			"6",
		},
		Links: map[string]*url.URL{
			"Test": {Path: "test.com"},
		},
	}

	expectedMessage := Message{
		channel: "channel",
		blocks: replacer.Replace(`
		[
			{
				"type": "header",
				"text": {
					"type": "plain_text",
					"text": "Test",
					"emoji": true
				}
			},
			{
				"type": "section",
				"text": {
					"type": "mrkdwn",
					"text": "Test"
				}
			},
			{
				"type": "section",
				"text": {
					"type": "mrkdwn",
					"text": "LOGMESSAGE"
				}
			},
			{
				"type":"actions",
				"block_id":"test",
				"elements":[
					{
						"type": "button",
						"text": {
							"type": "plain_text",
							"text": "Test"
						},
						"action_id": "test",
						"url": "test.com",
						"value": "test"
					}
				]
			}
		]
		`),
	}

	// Act
	err := handler.Handle(ctx, notification)
	if err != nil {
		t.Fatal(err)
	}

	// Assert
	assert.Equal(t, expectedMessage, service.Messages[0])
}
