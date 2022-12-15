package slack

import (
	log "github.com/sirupsen/logrus"

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
	blocks := createTimeoutBlocks(cfg, job)
	s.send(blocks)
}

func (s SlackHandler) send(blocks []slack.Block) {
	client := slack.New(s.Token)

	_, _, _, err := client.SendMessage(s.Channel, slack.MsgOptionBlocks(blocks...))
	if err != nil {
		log.Errorf("Failed to Send Message with error: %s!\n", err.Error())
	}
}
