package slack

import (
	"context"

	"github.com/slack-go/slack"
)

// Service Facade

type SlackService interface {
	Send(ctx context.Context, channel string, blocks []slack.Block) error
}

// Service Implementation

type slackService struct {
	client *slack.Client
}

func NewSlackService(token string) SlackService {
	return &slackService{
		client: slack.New(token),
	}
}

func (s *slackService) Send(ctx context.Context, channel string, blocks []slack.Block) error {
	_, _, _, err := s.client.SendMessageContext(ctx, channel, slack.MsgOptionBlocks(blocks...))
	return err
}
