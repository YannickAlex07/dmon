package slack

import (
	"fmt"

	"github.com/slack-go/slack"
	"github.com/yannickalex07/dmon/models"
)

type SlackHandler struct {
	Token   string
	Channel string
}

func (s SlackHandler) HandleError(cfg models.Config, job models.Job, entries []models.LogEntry) {
	blocks := createErrorBlocks(cfg, job, entries)
	s.send(blocks)
}

func (s SlackHandler) HandleTimeout(cfg models.Config, job models.Job) {
	// TODO: Implement Timeout Handler
}

func (s SlackHandler) send(blocks []slack.Block) {
	client := slack.New(s.Token)

	_, _, _, err := client.SendMessage(s.Channel, slack.MsgOptionBlocks(blocks...))
	if err != nil {
		fmt.Printf("Failed to Send Message with error: %s!\n", err.Error())
	}
}
