package slack

import (
	"fmt"
	"strings"

	"github.com/slack-go/slack"
	"github.com/yannickalex07/dataflow-monitor/dataflow"
)

func errorBlocks(job dataflow.Job, messages []dataflow.Message, projectId string, location string) []slack.Block {
	blocks := make([]slack.Block, 0)

	// Title
	titleBlock := slack.NewTextBlockObject("plain_text", "🚨 Pipeline Crashed 🚨", true, false)
	titleHeaderBlock := slack.NewHeaderBlock(titleBlock)
	blocks = append(blocks, titleHeaderBlock)

	// Info Section
	infoText := fmt.Sprintf("The job `%s` with the job id `%s` failed!", job.Name, job.Id)
	infoTextBlock := slack.NewTextBlockObject("mrkdwn", infoText, false, false)
	infoSectionBlock := slack.NewSectionBlock(infoTextBlock, nil, nil)
	blocks = append(blocks, infoSectionBlock)

	// GCP Button
	gcpTextBlock := slack.NewTextBlockObject("plain_text", "Open in Dataflow", false, false)
	gcpButtonBlock := slack.NewButtonBlockElement("1", "1", gcpTextBlock)
	gcpButtonBlock.URL = fmt.Sprintf("https://console.cloud.google.com/dataflow/jobs/%s/%s?project=%s&authuser=1&hl=en", location, job.Id, projectId)

	gcpButtonActionBlock := slack.NewActionBlock("1", gcpButtonBlock)
	blocks = append(blocks, gcpButtonActionBlock)

	// Error header
	errorTitleBlock := slack.NewTextBlockObject("plain_text", "Error Message", true, false)
	errorTitleHeaderBlock := slack.NewHeaderBlock(errorTitleBlock)
	blocks = append(blocks, errorTitleHeaderBlock)

	// Error Text
	msgParts := strings.Split(messages[0].Text, "\n")
	msg := msgParts[len(msgParts)-2] // last line is a blank line - before that comes the last error message
	errorText := fmt.Sprintf("The last error message was: ```%s```", msg)

	errorTextBlock := slack.NewTextBlockObject("mrkdwn", errorText, false, false)
	errorSectionBlock := slack.NewSectionBlock(errorTextBlock, nil, nil)
	blocks = append(blocks, errorSectionBlock)

	return blocks
}

func SendErrorMessage(token string, channel string, job dataflow.Job, messages []dataflow.Message, projectId string, location string) {
	blocks := errorBlocks(job, messages, projectId, location)
	sendMessage(token, channel, blocks)
}
