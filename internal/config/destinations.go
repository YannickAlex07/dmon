package config

import (
	"context"
	"errors"

	"github.com/yannickalex07/inframon/internal/config/destinations"
	inframon "github.com/yannickalex07/inframon/pkg"
)

type DestinationsConfig struct {
	Slack []destinations.SlackDestinationConfig `yaml:"slack"`
	Log   []destinations.LogDestinationConfig   `yaml:"log"`
}

func (c *DestinationsConfig) Get(ctx context.Context) ([]inframon.Destination, error) {
	var d []inframon.Destination

	// slack
	for _, s := range c.Slack {
		d = append(d, s.ToDestination())
	}

	// log
	for _, l := range c.Log {
		d = append(d, l.ToDestination())
	}

	if len(d) == 0 {
		return nil, errors.New("no destinations configured")
	}

	return d, nil
}
