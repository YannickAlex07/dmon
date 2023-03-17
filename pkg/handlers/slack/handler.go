package slack

import (
	log "github.com/sirupsen/logrus"

	"github.com/slack-go/slack"
	"github.com/yannickalex07/dmon/pkg/models"
)

type SlackGCPConfig struct {
	Id       string
	Location string
}

type SlackHandler struct {
	Token   string
	Channel string

	IncludeErrorSection   bool
	IncludeDataflowButton bool

	GCPConfig SlackGCPConfig
}

func (s SlackHandler) HandleError(job models.Job, entries []models.LogEntry) {
	blocks := s.createErrorBlocks(job, entries)
	s.send(blocks)
}

func (s SlackHandler) HandleTimeout(job models.Job) {
	blocks := s.createTimeoutBlocks(job)
	s.send(blocks)
}

func (s SlackHandler) send(blocks []slack.Block) {
	client := slack.New(s.Token)

	_, _, _, err := client.SendMessage(s.Channel, slack.MsgOptionBlocks(blocks...))
	if err != nil {
		log.Errorf("Failed to Send Message with error: %s!\n", err.Error())
	}
}
