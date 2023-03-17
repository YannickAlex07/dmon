package slack

import (
	"fmt"
	"strings"
	"time"

	"github.com/slack-go/slack"
	"github.com/yannickalex07/dmon/pkg/models"
)

func createErrorBlocks(cfg models.Config, job models.Job, entries []models.LogEntry) []slack.Block {
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
	if cfg.Slack.IncludeErrorSection {
		// Error Text
		msgParts := strings.Split(entries[0].Text, "\n")
		msg := msgParts[len(msgParts)-2] // last line is a blank line - before that comes the last error message
		errorText := fmt.Sprintf("Error Message: ```%s```", msg)

		errorTextBlock := slack.NewTextBlockObject("mrkdwn", errorText, false, false)
		errorSectionBlock := slack.NewSectionBlock(errorTextBlock, nil, nil)
		blocks = append(blocks, errorSectionBlock)
	}

	// Dataflow Button
	if cfg.Slack.IncludeDataflowButton {
		gcpTextBlock := slack.NewTextBlockObject("plain_text", "Open in Dataflow UI", false, false)
		gcpButtonBlock := slack.NewButtonBlockElement("dataflow_ui", "", gcpTextBlock)
		gcpButtonBlock.URL = fmt.Sprintf("https://console.cloud.google.com/dataflow/jobs/%s/%s?project=%s&authuser=1&hl=en", cfg.Project.Location, job.Id, cfg.Project.Id)
		gcpButtonActionBlock := slack.NewActionBlock("dataflow-button", gcpButtonBlock)
		blocks = append(blocks, gcpButtonActionBlock)
	}

	return blocks
}

func createTimeoutBlocks(cfg models.Config, job models.Job) []slack.Block {
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

	if cfg.Slack.IncludeDataflowButton {
		gcpTextBlock := slack.NewTextBlockObject("plain_text", "Open in Dataflow UI", false, false)
		gcpButtonBlock := slack.NewButtonBlockElement("dataflow_ui", "", gcpTextBlock)
		gcpButtonBlock.URL = fmt.Sprintf("https://console.cloud.google.com/dataflow/jobs/%s/%s?project=%s&authuser=1&hl=en", cfg.Project.Location, job.Id, cfg.Project.Id)
		gcpButtonActionBlock := slack.NewActionBlock("dataflow-button", gcpButtonBlock)
		blocks = append(blocks, gcpButtonActionBlock)
	}

	return blocks
}
