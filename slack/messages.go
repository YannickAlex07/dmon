package slack

import (
	"fmt"

	"github.com/slack-go/slack"
	"github.com/yannickalex07/dataflow-monitor/dataflow"
)

func ErrorBlocks(job dataflow.Job) []slack.Block {
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
	gcpButtonBlock.URL = "https://www.google.com"

	gcpButtonActionBlock := slack.NewActionBlock("1", gcpButtonBlock)
	blocks = append(blocks, gcpButtonActionBlock)

	// Error header
	errorTitleBlock := slack.NewTextBlockObject("plain_text", "Error Message", true, false)
	errorTitleHeaderBlock := slack.NewHeaderBlock(errorTitleBlock)
	blocks = append(blocks, errorTitleHeaderBlock)

	// Error Text
	errorText := fmt.Sprintf("The last error message was: ```%s```", "**NOT IMPLEMENTED YET**")
	errorTextBlock := slack.NewTextBlockObject("mrkdwn", errorText, false, false)
	errorSectionBlock := slack.NewSectionBlock(errorTextBlock, nil, nil)
	blocks = append(blocks, errorSectionBlock)

	return blocks
}

func SendMessage(token string, channel string, blocks []slack.Block) {
	client := slack.New(token)

	_, _, _, err := client.SendMessage(channel, slack.MsgOptionBlocks(blocks...))
	if err != nil {
		fmt.Printf("Failed to Send Message with error: %s!\n", err.Error())
	}
}
