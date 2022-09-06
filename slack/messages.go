package slack

import (
	"fmt"

	"github.com/slack-go/slack"
)

func sendMessage(token string, channel string, blocks []slack.Block) {
	client := slack.New(token)

	_, _, _, err := client.SendMessage(channel, slack.MsgOptionBlocks(blocks...))
	if err != nil {
		fmt.Printf("Failed to Send Message with error: %s!\n", err.Error())
	}
}
