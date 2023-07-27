package handler

import (
	"fmt"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/yannickalex07/dmon/pkg/model"

	"github.com/slack-go/slack"
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

func (s SlackHandler) HandleError(job model.Job, entries []model.LogEntry) {
	blocks := s.createErrorBlocks(job, entries)
	s.send(blocks)
}

func (s SlackHandler) HandleTimeout(job model.Job) {
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

func (s SlackHandler) createErrorBlocks(job model.Job, entries []model.LogEntry) []slack.Block {
	blocks := make([]slack.Block, 0)

	// Title
	titleBlock := slack.NewTextBlockObject("plain_text", "❌ Job Failed", true, false)
	titleHeaderBlock := slack.NewHeaderBlock(titleBlock)
	blocks = append(blocks, titleHeaderBlock)

	// Info Section
	infoText := fmt.Sprintf("The job `%s` with id `%s` failed at *%s*!", job.Name, job.Id, job.Status.UpdatedAt.Format(time.RFC1123))
	infoTextBlock := slack.NewTextBlockObject("mrkdwn", infoText, false, false)
	infoSectionBlock := slack.NewSectionBlock(infoTextBlock, nil, nil)
	blocks = append(blocks, infoSectionBlock)

	// Error Section
	if s.IncludeErrorSection {
		// Error Text
		msgParts := strings.Split(entries[0].Text, "\n")
		msg := msgParts[len(msgParts)-2] // last line is a blank line - before that comes the last error message
		errorText := fmt.Sprintf("Error Message: ```%s```", msg)

		errorTextBlock := slack.NewTextBlockObject("mrkdwn", errorText, false, false)
		errorSectionBlock := slack.NewSectionBlock(errorTextBlock, nil, nil)
		blocks = append(blocks, errorSectionBlock)
	}

	// Dataflow Button
	if s.IncludeDataflowButton {
		gcpTextBlock := slack.NewTextBlockObject("plain_text", "Open in Dataflow UI", false, false)
		gcpButtonBlock := slack.NewButtonBlockElement("dataflow_ui", "", gcpTextBlock)
		gcpButtonBlock.URL = fmt.Sprintf("https://console.cloud.google.com/dataflow/jobs/%s/%s?project=%s&authuser=1&hl=en", s.GCPConfig.Location, job.Id, s.GCPConfig.Id)
		gcpButtonActionBlock := slack.NewActionBlock("dataflow-button", gcpButtonBlock)
		blocks = append(blocks, gcpButtonActionBlock)
	}

	return blocks
}

func (s SlackHandler) createTimeoutBlocks(job model.Job) []slack.Block {
	blocks := make([]slack.Block, 0)

	// Title
	titleBlock := slack.NewTextBlockObject("plain_text", "⚠️ Job Timeout", true, false)
	titleHeaderBlock := slack.NewHeaderBlock(titleBlock)
	blocks = append(blocks, titleHeaderBlock)

	// Info Section
	infoText := fmt.Sprintf("The job `%s` with id `%s` crossed the maximum timeout limit with a runtime of *%s*.", job.Name, job.Id, job.Runtime().Round(time.Second))
	infoTextBlock := slack.NewTextBlockObject("mrkdwn", infoText, false, false)
	infoSectionBlock := slack.NewSectionBlock(infoTextBlock, nil, nil)
	blocks = append(blocks, infoSectionBlock)

	if s.IncludeDataflowButton {
		gcpTextBlock := slack.NewTextBlockObject("plain_text", "Open in Dataflow UI", false, false)
		gcpButtonBlock := slack.NewButtonBlockElement("dataflow_ui", "", gcpTextBlock)
		gcpButtonBlock.URL = fmt.Sprintf("https://console.cloud.google.com/dataflow/jobs/%s/%s?project=%s&authuser=1&hl=en", s.GCPConfig.Location, job.Id, s.GCPConfig.Id)
		gcpButtonActionBlock := slack.NewActionBlock("dataflow-button", gcpButtonBlock)
		blocks = append(blocks, gcpButtonActionBlock)
	}

	return blocks
}
