package slack

import (
	"context"
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
	inframon "github.com/yannickalex07/inframon/pkg"
)

type SlackDestination struct {
	Service SlackService

	Channel string
}

func (s *SlackDestination) Handle(ctx context.Context, notification inframon.Notification) error {
	return s.Service.Send(ctx, s.Channel, s.convertToBlocks(notification))
}

func (s *SlackDestination) convertToBlocks(notification inframon.Notification) []slack.Block {
	blocks := []slack.Block{}

	// title
	titleBlock := slack.NewTextBlockObject("plain_text", notification.Title, true, false)
	titleHeaderBlock := slack.NewHeaderBlock(titleBlock)
	blocks = append(blocks, titleHeaderBlock)

	// description
	infoTextBlock := slack.NewTextBlockObject("mrkdwn", notification.Description, false, false)
	infoSectionBlock := slack.NewSectionBlock(infoTextBlock, nil, nil)
	blocks = append(blocks, infoSectionBlock)

	// logs
	if len(notification.Logs) > 0 {
		var logParts []string

		if len(notification.Logs) > 5 {
			logParts = notification.Logs[len(notification.Logs)-5:]
		} else {
			logParts = notification.Logs
		}

		logStr := strings.Join(logParts, "\n")
		errorText := fmt.Sprintf("Last Log Messages: ```%s```", logStr)

		log.Println(errorText)

		errorTextBlock := slack.NewTextBlockObject("mrkdwn", errorText, false, false)
		errorSectionBlock := slack.NewSectionBlock(errorTextBlock, nil, nil)
		blocks = append(blocks, errorSectionBlock)
	}

	// buttons
	for title, link := range notification.Links {
		idTitle := strings.Replace(title, " ", "_", -1)
		idTitle = strings.ToLower(idTitle)

		textBlock := slack.NewTextBlockObject("plain_text", title, false, false)

		buttonBlock := slack.NewButtonBlockElement(idTitle, idTitle, textBlock)
		buttonBlock.URL = link.String()

		buttonActionBlock := slack.NewActionBlock(idTitle, buttonBlock)
		blocks = append(blocks, buttonActionBlock)
	}

	return blocks
}
