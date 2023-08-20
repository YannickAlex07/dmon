package handler

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/yannickalex07/dmon/pkg/model"

	"github.com/slack-go/slack"
)

const dataflowUrl string = "https://console.cloud.google.com/dataflow/jobs/%s/%s?project=%s&authuser=1&hl=en"

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

func (s SlackHandler) HandleError(ctx context.Context, job model.Job, entries []model.LogEntry) error {
	blocks := s.createErrorBlocks(job, entries)
	return s.send(blocks)
}

func (s SlackHandler) HandleTimeout(ctx context.Context, job model.Job) error {
	blocks := s.createTimeoutBlocks(job)
	return s.send(blocks)
}

func (s SlackHandler) send(blocks []slack.Block) error {
	client := slack.New(s.Token)

	_, _, _, err := client.SendMessage(s.Channel, slack.MsgOptionBlocks(blocks...))
	if err != nil {
		return fmt.Errorf("failed to send message with error: %w", err)
	}

	return nil
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
		if len(entries) > 0 {
			// Error Text
			cleaned := strings.TrimSpace(entries[0].Text)
			msgParts := strings.Split(cleaned, "\n")
			msg := msgParts[len(msgParts)-1] // last line is a blank line - before that comes the last error message
			errorText := fmt.Sprintf("Error Message: ```%s```", msg)

			errorTextBlock := slack.NewTextBlockObject("mrkdwn", errorText, false, false)
			errorSectionBlock := slack.NewSectionBlock(errorTextBlock, nil, nil)
			blocks = append(blocks, errorSectionBlock)
		} else {
			errorText := "```Failed to fetch log entries.```"
			errorTextBlock := slack.NewTextBlockObject("mrkdwn", errorText, false, false)
			errorSectionBlock := slack.NewSectionBlock(errorTextBlock, nil, nil)
			blocks = append(blocks, errorSectionBlock)
		}
	}

	// Dataflow Button
	if s.IncludeDataflowButton {
		gcpTextBlock := slack.NewTextBlockObject("plain_text", "Open in Dataflow UI", false, false)
		gcpButtonBlock := slack.NewButtonBlockElement("dataflow_ui", "", gcpTextBlock)
		gcpButtonBlock.URL = fmt.Sprintf(dataflowUrl, s.GCPConfig.Location, job.Id, s.GCPConfig.Id)
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
		gcpButtonBlock.URL = fmt.Sprintf(dataflowUrl, s.GCPConfig.Location, job.Id, s.GCPConfig.Id)
		gcpButtonActionBlock := slack.NewActionBlock("dataflow-button", gcpButtonBlock)
		blocks = append(blocks, gcpButtonActionBlock)
	}

	return blocks
}
