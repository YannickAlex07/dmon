package slack_test

import (
	"context"
	"encoding/json"

	"github.com/slack-go/slack"
)

type Message struct {
	channel string
	blocks  string
}

// DataflowServiceMock

type SlackServiceMock struct {
	Messages []Message
}

func (s *SlackServiceMock) Send(ctx context.Context, channel string, blocks []slack.Block) error {
	serializedBlocks, err := json.Marshal(blocks)
	if err != nil {
		return err
	}

	s.Messages = append(s.Messages, Message{channel, string(serializedBlocks)})
	return nil
}
