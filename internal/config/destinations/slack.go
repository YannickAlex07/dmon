package destinations

import (
	inframon "github.com/yannickalex07/inframon/pkg"
	"github.com/yannickalex07/inframon/pkg/slack"
)

type SlackDestinationConfig struct {
	Token   string `validate:"empty=false" yaml:"token"`
	Channel string `validate:"empty=false" yaml:"channel"`
}

func (c *SlackDestinationConfig) ToDestination() inframon.Destination {
	d := slack.SlackDestination{
		Service: slack.NewSlackService(c.Token),
		Channel: c.Channel,
	}

	return &d
}
